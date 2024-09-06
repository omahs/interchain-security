package ibc_testing

import (
	"encoding/json"
	"fmt"
	"testing"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"

	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	testkeeper "github.com/cosmos/interchain-security/v6/testutil/keeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cometbft/cometbft/abci/types"
	tmencoding "github.com/cometbft/cometbft/crypto/encoding"
	tmtypes "github.com/cometbft/cometbft/types"

	testutil "github.com/cosmos/interchain-security/v6/testutil/integration"
	consumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"
)

type (
	AppIniter       func() (ibctesting.TestingApp, map[string]json.RawMessage)
	ValSetAppIniter func([]types.ValidatorUpdate) AppIniter
)

// Contains generic setup code for running integration tests against a provider, consumer,
// and/or democracy consumer app.go implementation. You should not need to modify or replicate this file
// to run integration tests against your app.go implementations!

const (
	// Default number of consumer chains
	NumConsumers = 5
)

var (
	firstConsumerChainID string
	FirstConsumerID      string
	provChainID          string
	democConsumerChainID string
	consumerTopNParams   [NumConsumers]uint32
)

func init() {
	// Disable revision format
	ibctesting.ChainIDSuffix = ""
	firstConsumerChainID = ibctesting.GetChainID(2)
	FirstConsumerID = ""
	provChainID = ibctesting.GetChainID(1)
	democConsumerChainID = ibctesting.GetChainID(5000)
	// TopN parameter values per consumer chain initiated
	// sorted in ascending order i.e. testchain2, testchain3, ..., testchain6
	consumerTopNParams = [NumConsumers]uint32{100, 100, 100, 100, 100}
}

// ConsumerBundle serves as a way to store useful in-mem consumer app chain state
// and relevant IBC paths in the context of CCV integration testing.
type ConsumerBundle struct {
	ConsumerId   string
	Chain        *ibctesting.TestChain
	App          testutil.ConsumerApp
	Path         *ibctesting.Path
	TransferPath *ibctesting.Path
	TopN         uint32
}

// GetCtx returns the context for the ConsumerBundle
func (cb ConsumerBundle) GetCtx() sdk.Context {
	return cb.Chain.GetContext()
}

// GetKeeper returns the keeper for the ConsumerBundle
func (cb ConsumerBundle) GetKeeper() consumerkeeper.Keeper {
	return cb.App.GetConsumerKeeper()
}

// AddProvider adds a new provider chain to the coordinator and returns the test chain and app type
func AddProvider[T testutil.ProviderApp](t *testing.T, coordinator *ibctesting.Coordinator, appIniter AppIniter) (
	*ibctesting.TestChain, T,
) {
	t.Helper()
	ibctesting.DefaultTestingAppInit = appIniter
	provider := ibctesting.NewTestChain(t, coordinator, provChainID)
	coordinator.Chains[provChainID] = provider

	providerToReturn, ok := provider.App.(T)
	if !ok {
		panic(fmt.Sprintf("provider app type returned from app initer does not match app type passed in as type param: %T, %T",
			provider.App, *new(T)))
	}
	return provider, providerToReturn
}

// AddDemocracyConsumer adds a new democ consumer chain to the coordinator and returns the test chain and app type
func AddDemocracyConsumer[T testutil.DemocConsumerApp](
	coordinator *ibctesting.Coordinator,
	s *suite.Suite,
	appIniter ValSetAppIniter,
) (*ibctesting.TestChain, T) {
	s.T().Helper()

	// generate validators private/public key
	valSet, valUpdates, signers, err := testutil.CreateValidators(4, "")
	require.NoError(s.T(), err)

	ibctesting.DefaultTestingAppInit = appIniter(valUpdates)
	democConsumer := ibctesting.NewTestChainWithValSet(s.T(), coordinator, democConsumerChainID, valSet, signers)
	coordinator.Chains[democConsumerChainID] = democConsumer

	democConsumerToReturn, ok := democConsumer.App.(T)
	if !ok {
		panic(fmt.Sprintf("democ consumer app type returned from app initer does not match app type passed in as type param: %T, %T",
			democConsumer.App, *new(T)))
	}
	return democConsumer, democConsumerToReturn
}

// AddConsumer adds a new consumer chain with "testchain<index+2>" as chainID to the coordinator
// and returns the test chain and app type. A new client is created on the provider to the new
// consumer chain (see CreateConsumerClient). The new consumer is initialized with the
// InitialValSet from the genesis state generated by the provider (see MakeConsumerGenesis).
//
// This method must be called after AddProvider.
func AddConsumer[Tp testutil.ProviderApp, Tc testutil.ConsumerApp](
	coordinator *ibctesting.Coordinator,
	s *suite.Suite,
	index int,
	appIniter ValSetAppIniter,
) *ConsumerBundle {
	// check index isn't bigger that the number of consumers
	s.Require().LessOrEqual(index, NumConsumers)

	// consumer chain ID
	chainID := ibctesting.GetChainID(index + 2)

	// create client to the consumer on the provider chain
	providerChain := coordinator.Chains[provChainID]
	providerApp := providerChain.App.(Tp)
	providerKeeper := providerApp.GetProviderKeeper()

	consumerMetadata := testkeeper.GetTestConsumerMetadata()

	initializationParameters := testkeeper.GetTestInitializationParameters()
	// NOTE: we cannot use the time.Now() because the coordinator chooses a hardcoded start time
	// using time.Now() could set the spawn time to be too far in the past or too far in the future
	initializationParameters.SpawnTime = coordinator.CurrentTime
	// NOTE: the initial height passed to CreateConsumerClient
	// must be the height on the consumer when InitGenesis is called
	initializationParameters.InitialHeight = clienttypes.Height{RevisionNumber: 0, RevisionHeight: 2}

	powerShapingParameters := testkeeper.GetTestPowerShapingParameters()
	powerShapingParameters.Top_N = consumerTopNParams[index] // isn't used in CreateConsumerClient

	consumerId := providerKeeper.FetchAndIncrementConsumerId(providerChain.GetContext())
	if chainID == firstConsumerChainID {
		FirstConsumerID = consumerId
	}
	providerKeeper.SetConsumerChainId(providerChain.GetContext(), consumerId, chainID)
	err := providerKeeper.SetConsumerMetadata(providerChain.GetContext(), consumerId, consumerMetadata)
	s.Require().NoError(err)
	err = providerKeeper.SetConsumerInitializationParameters(providerChain.GetContext(), consumerId, initializationParameters)
	s.Require().NoError(err)
	err = providerKeeper.SetConsumerPowerShapingParameters(providerChain.GetContext(), consumerId, powerShapingParameters)
	s.Require().NoError(err)
	providerKeeper.SetConsumerPhase(providerChain.GetContext(), consumerId, providertypes.CONSUMER_PHASE_INITIALIZED)
	err = providerKeeper.AppendConsumerToBeLaunched(providerChain.GetContext(), consumerId, coordinator.CurrentTime)
	s.Require().NoError(err)

	// opt-in all validators
	lastVals, err := providerApp.GetProviderKeeper().GetLastBondedValidators(providerChain.GetContext())
	s.Require().NoError(err)

	for _, v := range lastVals {
		consAddr, _ := v.GetConsAddr()
		providerKeeper.SetOptedIn(providerChain.GetContext(), consumerId, providertypes.NewProviderConsAddress(consAddr))
		err = providerKeeper.AppendOptedInConsumerId(providerChain.GetContext(), providertypes.NewProviderConsAddress(consAddr), consumerId)
		s.Require().NoError(err)
	}

	// commit the state on the provider chain
	// and create the client and genesis of consumer
	coordinator.CommitBlock(providerChain)

	// get genesis state created by the provider
	consumerGenesisState, found := providerKeeper.GetConsumerGenesis(
		providerChain.GetContext(),
		consumerId,
	)
	s.Require().True(found, "consumer genesis not found in AddConsumer")

	_, found = providerKeeper.GetConsumerClientId(
		providerChain.GetContext(),
		consumerId,
	)
	s.Require().True(found, "clientID not found in AddConsumer")

	// use InitialValSet as the valset on the consumer
	var valz []*tmtypes.Validator
	for _, update := range consumerGenesisState.Provider.InitialValSet {
		// tmPubKey update.PubKey
		tmPubKey, err := tmencoding.PubKeyFromProto(update.PubKey)
		s.Require().NoError(err, "failed to convert tendermint pubkey")
		valz = append(valz, &tmtypes.Validator{
			PubKey:           tmPubKey,
			VotingPower:      update.Power,
			Address:          tmPubKey.Address(),
			ProposerPriority: 0,
		})
	}

	// create and instantiate consumer chain
	ibctesting.DefaultTestingAppInit = appIniter(consumerGenesisState.Provider.InitialValSet)
	testChain := ibctesting.NewTestChainWithValSet(s.T(), coordinator, chainID,
		tmtypes.NewValidatorSet(valz), providerChain.Signers)
	coordinator.Chains[chainID] = testChain

	consumerToReturn, ok := testChain.App.(Tc)
	if !ok {
		panic(fmt.Sprintf("consumer app type returned from app initer does not match app type passed in as type param: %T, %T",
			testChain.App, *new(Tc)))
	}

	return &ConsumerBundle{
		ConsumerId: consumerId,
		Chain:      testChain,
		App:        consumerToReturn,
		TopN:       powerShapingParameters.Top_N,
	}
}
