package keeper

import (
	"fmt"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	consumertypes "github.com/cosmos/interchain-security/v3/x/ccv/consumer/types"
)

// TODO: will need good integration tests making sure this state is properly init, cleared, etc.

//
// Throttling with retries follows a FSM design:
//
// 1. "No slash": If no slash record exists, the consumer is permitted to send packets from the pending packets queue.
// The consumer starts in this state from genesis.
//
// 2. On the event that a slash packet is obtained from the head of the pending packets queue and sent,
// a consumer transitions from "No Slash" to "Standby". A slash record is created upon entry to this state,
// and the consumer is now restricted from sending anymore packets.
//
// The slash packet remains at the head of the pending packets queue within the "Standby" state.
//
// - If the consumer receives a V1Result ack from the provider, the consumer checks for a slash record,
// and if found, the consumer transitions from "Standby" to "No Slash". The slash record is cleared upon this transition,
// and the slash packet is popped from the pending packets queue.
// TODO: Make note of above design for v1 throttling providers in the ADR, and explain that consumers must upgrade first in prod (where double queueing may exist for some time).
//
// - Else if the consumer receives a reply from the provider that the slash packet was successfully handled,
// the consumer transitions from "Standby" to "No Slash". The slash record is cleared upon this transition,
// and the slash packet is popped from the pending packets queue.
//
// - Else if the consumer receives a reply from the provider that the slash packet was bounced (not handled),
// then SlashRecord.WaitingOnReply is set false, and the consumer retries sending the slash packet after a delay period.
//
// Once a retry is sent, the consumer enters a new cycle of the "Standby" state and the process repeats.
//
// Once the slash packet is successfully handled, the consumer transitions from "Standby" to "No Slash",
// the slash record is cleared upon this transition, and the slash packet is popped from the pending packets queue.
//
// This design is implemented below, and in relay.go under SendPackets() and OnAcknowledgementPacket().
//

// PacketSendingPermitted returns whether the consumer is allowed to send packets
// from the pending packets queue.
func (k Keeper) PacketSendingPermitted(ctx sdktypes.Context) bool {
	record, found := k.GetSlashRecord(ctx)
	if !found {
		// no slash record exists, send is permitted
		return true
	}
	if record.WaitingOnReply {
		// We are waiting on a reply from provider, block sending
		return false
	}
	// TODO: implement retry delay period as param
	// retryDelayPeriod := k.GetParams(ctx).RetryDelayPeriod
	retryDelayPeriod := time.Hour
	timeSinceSend := ctx.BlockTime().Sub(record.SendTime)
	// If retry delay period has elapsed, we can send again
	retryPeriodElapsed := timeSinceSend >= retryDelayPeriod
	return retryPeriodElapsed
}

func (k Keeper) UpdateSlashRecordOnSend(ctx sdktypes.Context) {
	record := consumertypes.NewSlashRecord(
		ctx.BlockTime(), // sendTime
		true,            // waitingOnReply
	)
	// We don't mind overwriting here, since this is either a retry or the first time we send a slash
	k.SetSlashRecord(ctx, record)
}

func (k Keeper) UpdateSlashRecordOnBounce(ctx sdktypes.Context) {
	record, found := k.GetSlashRecord(ctx)
	if !found {
		// This should never happen
		panic("could not find slash record, but reply was received from provider")
	}
	record.WaitingOnReply = false
	k.SetSlashRecord(ctx, record)
}

func (k Keeper) GetSlashRecord(ctx sdktypes.Context) (record consumertypes.SlashRecord, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(consumertypes.SlashRecordKey())
	if bz == nil {
		return record, false
	}
	err := record.Unmarshal(bz)
	if err != nil {
		// This should never happen
		panic(fmt.Sprintf("could not unmarshal slash record: %v", err))
	}
	return record, true
}

func (k Keeper) SetSlashRecord(ctx sdktypes.Context, record consumertypes.SlashRecord) {
	store := ctx.KVStore(k.storeKey)
	bz, err := record.Marshal()
	if err != nil {
		// This should never happen
		panic(fmt.Sprintf("could not marshal slash record: %v", err))
	}
	store.Set(consumertypes.SlashRecordKey(), bz)
}

func (k Keeper) ClearSlashRecord(ctx sdktypes.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(consumertypes.SlashRecordKey())
}
