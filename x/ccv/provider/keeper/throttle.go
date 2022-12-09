package keeper

import (
	"encoding/binary"
	"fmt"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	providertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	ccvtypes "github.com/cosmos/interchain-security/x/ccv/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// This file contains functionality relevant to the throttling of slash and vsc matured packets, aka circuit breaker logic.

// HandlePendingSlashPackets handles all or some portion of pending slash packets received by consumer chains,
// depending on throttling logic. The slash meter is decremented appropriately in this method.
func (k Keeper) HandlePendingSlashPackets(ctx sdktypes.Context) {

	meter := k.GetSlashMeter(ctx)

	// Don't start iterating if meter is 0 or negative in value
	if !meter.IsPositive() {
		return
	}

	handledEntries := []providertypes.SlashPacketEntry{}

	// Iterate through ordered (by received time) slash packet entries from any consumer chain
	k.IteratePendingSlashPacketEntries(ctx, func(entry providertypes.SlashPacketEntry) (stop bool) {

		// Obtain validator from consensus address.
		// The slash packet validator address may be known only on the consumer chain
		// in this case, it must be mapped back to the consensus address on the provider chain
		consumerConsAddr := sdktypes.ConsAddress(entry.ValAddr)
		providerConsAddr := k.GetProviderAddrFromConsumerAddr(ctx, entry.ConsumerChainID, consumerConsAddr)

		// Note: if validator is not found or unbonded, this will be handled appropriately in HandleSlashPacket
		val, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, providerConsAddr)

		// Obtain the validator power relevant to the slash packet that's about to be handled
		// (this power will be removed via jailing or tombstoning)
		var valPower int64
		if !found {
			// If validator is jailed or not found, it's power is 0. This path is explicitly defined since the
			// staking keeper's LastValidatorPower values are not updated till the staking keeper's endblocker.
			valPower = 0
			k.Logger(ctx).Debug("handling slash entry, validator not found", "chainID", entry.ConsumerChainID, "provider cons addr", providerConsAddr)
		} else if val.IsJailed() {
			// If validator is jailed or not found, it's power is 0. This path is explicitly defined since the
			// staking keeper's LastValidatorPower values are not updated till the staking keeper's endblocker.
			valPower = 0
			k.Logger(ctx).Debug("handling slash entry, validator already jailed", "chainID", entry.ConsumerChainID, "provider cons addr", providerConsAddr)
		} else {
			valPower = k.stakingKeeper.GetLastValidatorPower(ctx, val.GetOperator())
			k.Logger(ctx).Debug("handling slash entry", "chainID", entry.ConsumerChainID, "provider cons addr", providerConsAddr, "val power at beginning of block", valPower)
		}

		// Subtract this power from the slash meter
		meter = meter.Sub(sdktypes.NewInt(valPower))

		k.Logger(ctx).Debug("meter after subtracting val power", "meter", meter, "val power", valPower)

		// Handle slash packet by passing in chainID and appropriate callbacks, relevant packet data is deleted in this method
		k.HandlePacketDataForChain(ctx, entry.ConsumerChainID, k.HandleSlashPacket, k.HandleVSCMaturedPacket)

		// Store handled entry to be deleted after iteration is completed
		handledEntries = append(handledEntries, entry)

		// Do not handle anymore slash packets if the meter is 0 or negative in value
		stop = !meter.IsPositive()

		k.Logger(ctx).Debug("handled slash packet entry", "chainID", entry.ConsumerChainID, "meter", meter, "stopping", stop)

		return stop
	})

	// Handled entries are deleted after iteration is completed
	k.DeletePendingSlashPacketEntries(ctx, handledEntries...)

	// Persist current value for slash meter
	k.SetSlashMeter(ctx, meter)
}

// HandlePacketDataForChain handles only the first queued slash packet relevant to the passed consumer chainID,
// and then handles any trailing vsc matured packets in that (consumer chain specific) queue.
//
// Note: Any packet data which is handled in this method is also deleted from the (consumer chain specific) queue.
func (k Keeper) HandlePacketDataForChain(ctx sdktypes.Context, consumerChainID string,
	slashPacketHandler func(sdktypes.Context, string, ccvtypes.SlashPacketData) (bool, error),
	vscMaturedPacketHandler func(sdktypes.Context, string, ccvtypes.VSCMaturedPacketData),
) {

	// Store ibc sequence numbers to delete data after iteration is completed
	seqNums := []uint64{}

	idx := 0
	k.IteratePendingPacketData(ctx, consumerChainID, func(ibcSeqNum uint64, data interface{}) (stop bool) {

		switch data := data.(type) {

		case ccvtypes.SlashPacketData:
			if idx == 0 {
				_, err := slashPacketHandler(ctx, consumerChainID, data)
				if err != nil {
					panic(fmt.Sprintf("failed to handle slash packet: %s", err))
				}
			} else {
				// Break iteration, since we've already handled one slash packet
				stop = true
				return stop
			}
		case ccvtypes.VSCMaturedPacketData:
			if idx == 0 {
				panic("data is corrupt, first data struct in queue should be slash packet data")
			} else {
				vscMaturedPacketHandler(ctx, consumerChainID, data)
			}
		default:
			panic(fmt.Sprintf("unexpected pending packet data type: %T", data))
		}
		seqNums = append(seqNums, ibcSeqNum)

		// Increment idx
		idx++

		// Continue iterating through the queue until we reach the end or a 2nd slash packet
		stop = false
		return stop
	})

	if 0 < len(seqNums) {
		k.Logger(ctx).Debug("handled consumer packets", "chainID", consumerChainID, "largest seq num handled", seqNums[len(seqNums)-1])
	}

	// Delete handled data after iteration is completed
	k.DeletePendingPacketData(ctx, consumerChainID, seqNums...)
}

// InitializeSlashMeter initializes the slash meter to it's max value (also its allowance
// per replenish period), and sets the last replenish time to the current block time.
func (k Keeper) InitializeSlashMeter(ctx sdktypes.Context) {
	k.SetSlashMeter(ctx, k.GetSlashMeterAllowance(ctx))
	k.SetLastSlashMeterReplenishTime(ctx, ctx.BlockTime())
}

// CheckForSlashMeterReplenishment checks if the slash gas meter should be replenished, and if so, replenishes it.
// Note: initial slash meter replenish time is set in InitGenesis
func (k Keeper) CheckForSlashMeterReplenishment(ctx sdktypes.Context) {
	lastReplenishTime := k.GetLastSlashMeterReplenishTime(ctx)
	replenishPeriod := k.GetSlashMeterReplenishPeriod(ctx)
	if ctx.BlockTime().UTC().After(lastReplenishTime.Add(replenishPeriod)) {
		k.ReplenishSlashMeter(ctx)
	}
}

func (k Keeper) ReplenishSlashMeter(ctx sdktypes.Context) {

	meter := k.GetSlashMeter(ctx)
	allowance := k.GetSlashMeterAllowance(ctx)

	// Replenish gas up to gas allowance per period. That is, if meter was negative
	// before being replenished, it'll gain some additional gas. However, if the meter
	// was 0 or positive in value, it'll be replenished only up to it's allowance for the period.
	meter = meter.Add(allowance)
	if meter.GT(allowance) {
		meter = allowance
	}
	k.SetSlashMeter(ctx, meter)
	k.SetLastSlashMeterReplenishTime(ctx, ctx.BlockTime())

	k.Logger(ctx).Debug("replenished slash meter", "meter", meter)
}

// GetSlashMeterAllowance returns the allowance of voting power units (int) that the slash meter
// is given per replenish period, this also serves as the max value for the meter.
func (k Keeper) GetSlashMeterAllowance(ctx sdktypes.Context) sdktypes.Int {

	strFrac := k.GetSlashMeterReplenishFraction(ctx)
	decFrac, err := sdktypes.NewDecFromStr(strFrac)
	if err != nil {
		panic(fmt.Sprintf("failed to parse slash meter replenish fraction: %s", err))
	}

	// Compute slash gas allowance in units of tendermint voting power (integer),
	// noting that total power changes over time
	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	roundedInt := sdktypes.NewInt(decFrac.MulInt(totalPower).RoundInt64())
	if roundedInt.IsZero() {
		// TODO: Log warning that replenish fraction is too small to
		// add any allowance to the meter, considering bankers rounding.

		// Return non-zero allowance to guarantee some slash packets are eventually handled
		return sdktypes.NewInt(1)
	}
	return roundedInt
}

//
// CRUD section
//

// QueuePendingSlashPacketEntry queues an entry to the "parent" slash packet queue, used for throttling val power changes
// related to jailing/tombstoning over time. This "parent" queue is used to coordinate the order of slash packet handling
// between chains, whereas the chain specific queue is used to coordinate the order of slash and vsc matured packets
// relevant to each chain.
func (k Keeper) QueuePendingSlashPacketEntry(ctx sdktypes.Context, entry providertypes.SlashPacketEntry) {
	store := ctx.KVStore(k.storeKey)
	key := providertypes.PendingSlashPacketEntryKey(entry)
	// Note: Val address is stored as value to assist in debugging. This could be removed for efficiency.
	store.Set(key, entry.ValAddr)
}

// GetAllPendingSlashPacketEntries returns all pending slash packet entries in the queue
// This method is used for testing purposes only
func (k Keeper) GetAllPendingSlashPacketEntries(ctx sdktypes.Context) (entries []providertypes.SlashPacketEntry) {
	k.IteratePendingSlashPacketEntries(ctx, func(entry providertypes.SlashPacketEntry) (stop bool) {
		entries = append(entries, entry)
		// Continue iteration
		stop = false
		return stop
	})
	return entries
}

// IteratePendingSlashPackets iterates over the pending slash packet entry queue and calls the provided callback
func (k Keeper) IteratePendingSlashPacketEntries(ctx sdktypes.Context,
	cb func(providertypes.SlashPacketEntry) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdktypes.KVStorePrefixIterator(store, []byte{providertypes.PendingSlashPacketEntryBytePrefix})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		recvTime, chainID := providertypes.ParsePendingSlashPacketEntryKey(iterator.Key())
		valAddr := iterator.Value()
		entry := providertypes.NewSlashPacketEntry(recvTime, chainID, valAddr)
		stop := cb(entry)
		if stop {
			break
		}
	}
}

// DeletePendingSlashPackets deletes the given entries from the pending slash packet entry queue
func (k Keeper) DeletePendingSlashPacketEntries(ctx sdktypes.Context, entries ...providertypes.SlashPacketEntry) {
	store := ctx.KVStore(k.storeKey)
	for _, entry := range entries {
		store.Delete(providertypes.PendingSlashPacketEntryKey(entry))
	}
}

// Pending packet data type enum, used to encode the type of packet data stored at each entry in the mutual queue.
const (
	slashPacketData byte = iota
	vscMaturedPacketData
)

// PanicIfTooMuchPendingPacketData is a sanity check to ensure that the pending packet data queue
// does not grow too large for a single consumer chain.
func (k Keeper) PanicIfTooMuchPendingPacketData(ctx sdktypes.Context, consumerChainID string) {
	size := k.GetPendingPacketDataSize(ctx, consumerChainID)
	if size > uint64(k.GetMaxPendingSlashingPackets(ctx)) {
		panic(fmt.Sprintf("pending packet data queue for chain %s is too large: %d", consumerChainID, size))
	}
}

// GetPendingPacketDataSize returns the size of the pending packet data queue for the given consumer chain
func (k Keeper) GetPendingPacketDataSize(ctx sdktypes.Context, consumerChainID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := providertypes.PendingPacketDataSizeKey(consumerChainID)
	var size uint64
	bz := store.Get(key)
	if bz == nil {
		size = 0
	} else {
		size = binary.LittleEndian.Uint64(bz)
	}
	return size
}

// SetPendingPacketDataSize sets the size of the pending packet data queue for the given consumer chain
func (k Keeper) SetPendingPacketDataSize(ctx sdktypes.Context, consumerChainID string, size uint64) {
	store := ctx.KVStore(k.storeKey)
	key := providertypes.PendingPacketDataSizeKey(consumerChainID)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, size)
	store.Set(key, bz)
}

// IncrementPendingPacketDataSize increments the size of the pending packet data queue for the given consumer chain
func (k Keeper) IncrementPendingPacketDataSize(ctx sdktypes.Context, consumerChainID string) {
	size := k.GetPendingPacketDataSize(ctx, consumerChainID)
	k.SetPendingPacketDataSize(ctx, consumerChainID, size+1)
}

// QueuePendingSlashPacketData queues the given slash packet data for the given consumer chain's queue
// Note: This queue is shared between pending slash packet data and pending vsc matured packet data
func (k Keeper) QueuePendingSlashPacketData(
	ctx sdktypes.Context, consumerChainID string, ibcSeqNum uint64, data ccvtypes.SlashPacketData) {

	k.PanicIfTooMuchPendingPacketData(ctx, consumerChainID)

	store := ctx.KVStore(k.storeKey)
	bz, err := data.Marshal()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal slash packet data: %v", err))
	}
	bz = append([]byte{slashPacketData}, bz...)
	store.Set(providertypes.PendingPacketDataKey(consumerChainID, ibcSeqNum), bz)

	k.IncrementPendingPacketDataSize(ctx, consumerChainID)
}

// QueuePendingVSCMaturedPacketData queues the given vsc matured packet data for the given consumer chain's queue
// Note: This queue is shared between pending slash packet data and pending vsc matured packet data
func (k Keeper) QueuePendingVSCMaturedPacketData(
	ctx sdktypes.Context, consumerChainID string, ibcSeqNum uint64, data ccvtypes.VSCMaturedPacketData) {

	k.PanicIfTooMuchPendingPacketData(ctx, consumerChainID)

	store := ctx.KVStore(k.storeKey)
	bz, err := data.Marshal()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal vsc matured packet data: %v", err))
	}
	bz = append([]byte{vscMaturedPacketData}, bz...)
	store.Set(providertypes.PendingPacketDataKey(consumerChainID, ibcSeqNum), bz)

	k.IncrementPendingPacketDataSize(ctx, consumerChainID)
}

// IteratePendingPacketData iterates over the pending packet data queue for a specific consumer chain
// (ordered by ibc seq number) and calls the provided callback
func (k Keeper) IteratePendingPacketData(ctx sdktypes.Context, consumerChainID string, cb func(uint64, interface{}) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iteratorPrefix := providertypes.ChainIdWithLenKey(providertypes.PendingPacketDataBytePrefix, consumerChainID)
	iterator := sdktypes.KVStorePrefixIterator(store, iteratorPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var packetData interface{}
		var err error
		bz := iterator.Value()
		switch bz[0] {
		case slashPacketData:
			spd := ccvtypes.SlashPacketData{}
			err = spd.Unmarshal(bz[1:])
			packetData = spd
		case vscMaturedPacketData:
			vpd := ccvtypes.VSCMaturedPacketData{}
			err = vpd.Unmarshal(bz[1:])
			packetData = vpd
		default:
			panic("invalid packet data type")
		}
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal pending packet data: %v", err))
		}
		_, ibcSeqNum, err := providertypes.ParsePendingPacketDataKey(iterator.Key())
		if err != nil {
			panic(fmt.Sprintf("failed to parse pending packet data key: %v", err))
		}
		stop := cb(ibcSeqNum, packetData)
		if stop {
			break
		}
	}
}

// DeletePendingPacketData deletes the given entries (specified by their ibc seq number) from the pending packet data queue
func (k Keeper) DeletePendingPacketData(ctx sdktypes.Context, consumerChainID string, ibcSeqNumbers ...uint64) {
	store := ctx.KVStore(k.storeKey)
	for _, ibcSeqNum := range ibcSeqNumbers {
		store.Delete(providertypes.PendingPacketDataKey(consumerChainID, ibcSeqNum))
	}
	// Decrement the size of the pending packet data queue
	sizeBeforeDeletion := k.GetPendingPacketDataSize(ctx, consumerChainID)
	k.SetPendingPacketDataSize(ctx, consumerChainID, sizeBeforeDeletion-uint64(len(ibcSeqNumbers)))
}

// GetSlashMeter returns a meter (persisted as a signed int) which stores an amount of voting power, corresponding
// to an allowance of validators that can be jailed/tombstoned over time.
//
// Note: the value of this int should always be in the range of tendermint's [-MaxVotingPower, MaxVotingPower]
func (k Keeper) GetSlashMeter(ctx sdktypes.Context) sdktypes.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(providertypes.SlashMeterKey())
	if bz == nil {
		panic("slash meter not set")
	}
	value := sdktypes.ZeroInt()
	err := value.Unmarshal(bz)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal slash meter: %v", err))
	}
	return value
}

// SetSlashMeter sets the slash meter to the given signed int value
//
// Note: the value of this int should always be in the range of tendermint's [-MaxVotingPower, MaxVotingPower]
func (k Keeper) SetSlashMeter(ctx sdktypes.Context, value sdktypes.Int) {
	if value.GT(sdktypes.NewInt(tmtypes.MaxTotalVotingPower)) {
		panic("slash meter value cannot be greater than tendermint's MaxTotalVotingPower")
	}
	if value.LT(sdktypes.NewInt(-tmtypes.MaxTotalVotingPower)) {
		panic("slash meter value cannot be less than negative tendermint's MaxTotalVotingPower")
	}
	store := ctx.KVStore(k.storeKey)
	bz, err := value.Marshal()
	if err != nil {
		panic(fmt.Sprintf("failed to marshal slash meter: %v", err))
	}
	store.Set(providertypes.SlashMeterKey(), bz)
}

// GetLastSlashMeterReplenishTime returns the last UTC time the slash meter was replenished
func (k Keeper) GetLastSlashMeterReplenishTime(ctx sdktypes.Context) time.Time {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(providertypes.LastSlashMeterReplenishTimeKey())
	if bz == nil {
		panic("last slash replenish time not set")
	}
	time, err := sdktypes.ParseTimeBytes(bz)
	if err != nil {
		panic(fmt.Sprintf("failed to parse last slash meter replenish time: %s", err))
	}
	return time.UTC()
}

// SetLastSlashMeterReplenishTime sets the last time the slash meter was replenished
func (k Keeper) SetLastSlashMeterReplenishTime(ctx sdktypes.Context, time time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Set(providertypes.LastSlashMeterReplenishTimeKey(), sdktypes.FormatTimeBytes(time.UTC()))
}
