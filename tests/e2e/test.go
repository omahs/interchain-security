package e2e

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *CCVTestSuite) TestSlash() {
	providerStakingKeeper := suite.providerApp.GetE2eStakingKeeper()
	providerSlashingKeeper := suite.providerApp.GetE2eSlashingKeeper()

	valIndex := 0
	bondAmt := sdk.NewInt(10000000)

	// get the power before delegate
	validator, _ := suite.getValByIdx(valIndex)
	powerBeforeDelegate := validator.GetConsensusPower(sdk.DefaultPowerReduction)
	heightDelegate := suite.providerCtx().BlockHeight()
	fmt.Printf("delegate at height %d\n", heightDelegate)

	// delegate some tokens
	delAddr := suite.providerChain.SenderAccount.GetAddress()
	initBalance, shares, valAddr := delegateByIdx(suite, delAddr, bondAmt, valIndex)

	// commit state on provider
	suite.coordinator.CommitBlock(suite.providerChain)

	incrementTime(suite, time.Hour)

	// create validator signing info to test slashing
	valConsAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	providerSlashingKeeper.SetValidatorSigningInfo(
		suite.providerChain.GetContext(),
		valConsAddr,
		slashingtypes.ValidatorSigningInfo{Address: valConsAddr.String()},
	)

	fmt.Printf("slash at height %d\n", heightDelegate-5)
	providerStakingKeeper.Slash(
		suite.providerCtx(),
		valConsAddr,
		heightDelegate-5,
		powerBeforeDelegate,
		sdk.NewDec(1),
		stakingtypes.Downtime,
	)

	// commit state on provider
	suite.coordinator.CommitBlock(suite.providerChain)

	// undelegate all the delegated shares
	undelegate(suite, delAddr, valAddr, shares)

	incrementTimeByUnbondingPeriod(suite, Provider)

	_, found := providerStakingKeeper.GetUnbondingDelegation(suite.providerCtx(), delAddr, valAddr)
	suite.Require().False(found)

	suite.Require().Equal(
		initBalance,
		getBalance(suite, suite.providerCtx(), delAddr),
		"delegator shouldn't have been slashed",
	)
}
