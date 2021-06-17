package commands_test

import (
	"errors"
	"strconv"
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	oraclespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/oracles/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckProposalSubmission(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilProposalSubmissionFails)
	t.Run("Submitting a proposal without terms fails", testProposalSubmissionWithoutTermsFails)
	t.Run("Submitting a proposal with non-positive closing timestamp fails", testProposalSubmissionWithNonPositiveClosingTimestampFails)
	t.Run("Submitting a proposal with positive closing timestamp succeeds", testProposalSubmissionWithPositiveClosingTimestampSucceeds)
	t.Run("Submitting a proposal with non-positive enactment timestamp fails", testProposalSubmissionWithNonPositiveEnactmentTimestampFails)
	t.Run("Submitting a proposal with positive enactment timestamp succeeds", testProposalSubmissionWithPositiveEnactmentTimestampSucceeds)
	t.Run("Submitting a proposal with negative validation timestamp fails", testProposalSubmissionWithNegativeValidationTimestampFails)
	t.Run("Submitting a proposal with positive validation timestamp succeeds", testProposalSubmissionWithPositiveValidationTimestampSucceeds)
	t.Run("Submitting a proposal with closing timestamp after enactment timestamp fails", testProposalSubmissionWithClosingTimestampAfterEnactmentTimestampFails)
	t.Run("Submitting a proposal with closing timestamp before enactment timestamp succeeds", testProposalSubmissionWithClosingTimestampBeforeEnactmentTimestampSucceeds)
	t.Run("Submitting a proposal with closing timestamp at enactment timestamp succeeds", testProposalSubmissionWithClosingTimestampAtEnactmentTimestampSucceeds)
	t.Run("Submitting a proposal with validation timestamp after closing timestamp fails", testProposalSubmissionWithValidationTimestampAfterClosingTimestampFails)
	t.Run("Submitting a proposal with validation timestamp at closing timestamp succeeds", testProposalSubmissionWithValidationTimestampAtClosingTimestampFails)
	t.Run("Submitting a proposal with validation timestamp before closing timestamp fails", testProposalSubmissionWithValidationTimestampBeforeClosingTimestampSucceeds)
	t.Run("Submitting an asset change without change fails", testProposalSubmissionWithoutChangeFails)
	t.Run("Submitting an asset change without new asset fails", testAssetChangeSubmissionWithoutNewsAssetFails)
	t.Run("Submitting an asset change without changes fails", testAssetChangeSubmissionWithoutChangesFails)
	t.Run("Submitting an asset change without source fails", testAssetChangeSubmissionWithoutSourceFails)
	t.Run("Submitting an built-in asset change without built-in asset fails", testBuiltInAssetChangeSubmissionWithoutBuiltInAssetFails)
	t.Run("Submitting an built-in asset change without name fails", testBuiltInAssetChangeSubmissionWithoutNameFails)
	t.Run("Submitting an built-in asset change with name succeeds", testBuiltInAssetChangeSubmissionWithNameSucceeds)
	t.Run("Submitting an built-in asset change without symbol fails", testBuiltInAssetChangeSubmissionWithoutSymbolFails)
	t.Run("Submitting an built-in asset change with symbol succeeds", testBuiltInAssetChangeSubmissionWithSymbolSucceeds)
	t.Run("Submitting an built-in asset change without decimal fails", testBuiltInAssetChangeSubmissionWithoutDecimalsFails)
	t.Run("Submitting an built-in asset change with decimal succeeds", testBuiltInAssetChangeSubmissionWithDecimalsSucceeds)
	t.Run("Submitting an built-in asset change without total supply fails", testBuiltInAssetChangeSubmissionWithoutTotalSupplyFails)
	t.Run("Submitting an built-in asset change with total supply succeeds", testBuiltInAssetChangeSubmissionWithTotalSupplySucceeds)
	t.Run("Submitting an built-in asset change with not-a-number total supply fails", testBuiltInAssetChangeSubmissionWithNaNTotalSupplyFails)
	t.Run("Submitting an built-in asset change without max faucet amount fails", testBuiltInAssetChangeSubmissionWithoutMaxFaucetAmountMintFails)
	t.Run("Submitting an built-in asset change with max faucet amount succeeds", testBuiltInAssetChangeSubmissionWithMaxFaucetAmountMintSucceeds)
	t.Run("Submitting an built-in asset change with not-a-number max faucet amount fails", testBuiltInAssetChangeSubmissionWithNaNMaxFaucetAmountMintFails)
	t.Run("Submitting an ERC20 asset change without ERC20 asset fails", testERC20AssetChangeSubmissionWithoutErc20AssetFails)
	t.Run("Submitting an ERC20 asset change without contract address fails", testErc20AssetChangeSubmissionWithoutContractAddressFails)
	t.Run("Submitting an ERC20 asset change with contract address succeeds", testErc20AssetChangeSubmissionWithContractAddressSucceeds)
	t.Run("Submitting a network parameter changes without network parameter fails", testNetworkParameterChangeSubmissionWithoutNetworkParameterFails)
	t.Run("Submitting a network parameter changes without changes fails", testNetworkParameterChangeSubmissionWithoutChangesFails)
	t.Run("Submitting a network parameter change without key fails", testNetworkParameterChangeSubmissionWithoutKeyFails)
	t.Run("Submitting a network parameter change with key succeeds", testNetworkParameterChangeSubmissionWithKeySucceeds)
	t.Run("Submitting a network parameter change without value fails", testNetworkParameterChangeSubmissionWithoutValueFails)
	t.Run("Submitting a network parameter change with value succeeds", testNetworkParameterChangeSubmissionWithValueSucceeds)
	t.Run("Submitting a market change without new market fails", testNewMarketChangeSubmissionWithoutNewMarketFails)
	t.Run("Submitting a market change without changes fails", testNewMarketChangeSubmissionWithoutChangesFails)
	t.Run("Submitting a market change without decimal places fails", testNewMarketChangeSubmissionWithoutDecimalPlacesFails)
	t.Run("Submitting a market change with decimal places above or equal to 150 fails", testNewMarketChangeSubmissionWithDecimalPlacesAboveOrEqualTo150Fails)
	t.Run("Submitting a market change with decimal places below 150 succeeds", testNewMarketChangeSubmissionWithDecimalPlacesBelow150Succeeds)
	t.Run("Submitting a new market without price monitoring succeeds", testNewMarketChangeSubmissionWithoutPriceMonitoringSucceeds)
	t.Run("Submitting a new market with price monitoring succeeds", testNewMarketChangeSubmissionWithPriceMonitoringSucceeds)
	t.Run("Submitting a price monitoring change without triggers succeeds", testPriceMonitoringChangeSubmissionWithoutTriggersSucceeds)
	t.Run("Submitting a price monitoring change with triggers succeeds", testPriceMonitoringChangeSubmissionWithTriggersSucceeds)
	t.Run("Submitting a price monitoring change without trigger horizon fails", testPriceMonitoringChangeSubmissionWithoutTriggerHorizonFails)
	t.Run("Submitting a price monitoring change with trigger horizon succeeds", testPriceMonitoringChangeSubmissionWithTriggerHorizonSucceeds)
	t.Run("Submitting a price monitoring change with wrong trigger probability fails", testPriceMonitoringChangeSubmissionWithWrongTriggerProbabilityFails)
	t.Run("Submitting a price monitoring change with right trigger probability succeeds", testPriceMonitoringChangeSubmissionWithRightTriggerProbabilitySucceeds)
	t.Run("Submitting a price monitoring change without trigger auction extension fails", testPriceMonitoringChangeSubmissionWithoutTriggerAuctionExtensionFails)
	t.Run("Submitting a price monitoring change with trigger auction extension succeeds", testPriceMonitoringChangeSubmissionWithTriggerAuctionExtensionSucceeds)
	t.Run("Submitting a new market without liquidity monitoring succeeds", testNewMarketChangeSubmissionWithoutLiquidityMonitoringSucceeds)
	t.Run("Submitting a new market with liquidity monitoring succeeds", testNewMarketChangeSubmissionWithLiquidityMonitoringSucceeds)
	t.Run("Submitting a liquidity monitoring change with wrong triggering ratio fails", testLiquidityMonitoringChangeSubmissionWithWrongTriggeringRatioFails)
	t.Run("Submitting a liquidity monitoring change with right triggering ratio succeeds", testLiquidityMonitoringChangeSubmissionWithRightTriggeringRatioSucceeds)
	t.Run("Submitting a liquidity monitoring change without target stake parameters fails", testLiquidityMonitoringChangeSubmissionWithoutTargetStakeParametersFails)
	t.Run("Submitting a liquidity monitoring change with target stake parameters succeeds", testLiquidityMonitoringChangeSubmissionWithTargetStakeParametersSucceeds)
	t.Run("Submitting a liquidity monitoring change with non-positive time window fails", testLiquidityMonitoringChangeSubmissionWithNonPositiveTimeWindowFails)
	t.Run("Submitting a liquidity monitoring change with positive time window succeeds", testLiquidityMonitoringChangeSubmissionWithPositiveTimeWindowSucceeds)
	t.Run("Submitting a liquidity monitoring change with non-positive scaling factor fails", testLiquidityMonitoringChangeSubmissionWithNonPositiveScalingFactorFails)
	t.Run("Submitting a liquidity monitoring change with positive scaling factor succeeds", testLiquidityMonitoringChangeSubmissionWithPositiveScalingFactorSucceeds)
	t.Run("Submitting a market change without instrument name fails", testNewMarketChangeSubmissionWithoutInstrumentNameFails)
	t.Run("Submitting a market change with instrument name succeeds", testNewMarketChangeSubmissionWithInstrumentNameSucceeds)
	t.Run("Submitting a market change without instrument code fails", testNewMarketChangeSubmissionWithoutInstrumentCodeFails)
	t.Run("Submitting a market change with instrument code succeeds", testNewMarketChangeSubmissionWithInstrumentCodeSucceeds)
	t.Run("Submitting a market change without product fails", testNewMarketChangeSubmissionWithoutProductFails)
	t.Run("Submitting a market change with product succeeds", testNewMarketChangeSubmissionWithProductSucceeds)
	t.Run("Submitting a future market change without future fails", testNewFutureMarketChangeSubmissionWithoutFutureFails)
	t.Run("Submitting a future market change with future succeeds", testNewFutureMarketChangeSubmissionWithFutureSucceeds)
	t.Run("Submitting a future market change without settlement asset fails", testNewFutureMarketChangeSubmissionWithoutSettlementAssetFails)
	t.Run("Submitting a future market change with settlement asset succeeds", testNewFutureMarketChangeSubmissionWithSettlementAssetSucceeds)
	t.Run("Submitting a future market change without quote name fails", testNewFutureMarketChangeSubmissionWithoutQuoteNameFails)
	t.Run("Submitting a future market change with quote name succeeds", testNewFutureMarketChangeSubmissionWithQuoteNameSucceeds)
	t.Run("Submitting a future market change without maturity fails", testNewFutureMarketChangeSubmissionWithoutMaturityFails)
	t.Run("Submitting a future market change with maturity succeeds", testNewFutureMarketChangeSubmissionWithMaturitySucceeds)
	t.Run("Submitting a future market change with wrong maturity date format fails", testNewFutureMarketChangeSubmissionWithWrongMaturityDateFormatFails)
	t.Run("Submitting a future market change with right maturity date format succeeds", testNewFutureMarketChangeSubmissionWithRightMaturityDateFormatSucceeds)
	t.Run("Submitting a future market change without oracle spec fails", testNewFutureMarketChangeSubmissionWithoutOracleSpecFails)
	t.Run("Submitting a future market change with oracle spec succeeds", testNewFutureMarketChangeSubmissionWithOracleSpecSucceeds)
	t.Run("Submitting a future market change without pub-keys fails", testNewFutureMarketChangeSubmissionWithoutPubKeysFails)
	t.Run("Submitting a future market change with wrong pub-keys fails", testNewFutureMarketChangeSubmissionWithWrongPubKeysFails)
	t.Run("Submitting a future market change with pub-keys succeeds", testNewFutureMarketChangeSubmissionWithPubKeysSucceeds)
	t.Run("Submitting a future market change without filters fails", testNewFutureMarketChangeSubmissionWithoutFiltersFails)
	t.Run("Submitting a future market change with filters succeeds", testNewFutureMarketChangeSubmissionWithFiltersSucceeds)
	t.Run("Submitting a future market change with filter without key fails", testNewFutureMarketChangeSubmissionWithFilterWithoutKeyFails)
	t.Run("Submitting a future market change with filter with key succeeds", testNewFutureMarketChangeSubmissionWithFilterWithKeySucceeds)
	t.Run("Submitting a future market change with filter without key name fails", testNewFutureMarketChangeSubmissionWithFilterWithoutKeyNameFails)
	t.Run("Submitting a future market change with filter with key name succeeds", testNewFutureMarketChangeSubmissionWithFilterWithKeyNameSucceeds)
	t.Run("Submitting a future market change with filter without key type fails", testNewFutureMarketChangeSubmissionWithFilterWithoutKeyTypeFails)
	t.Run("Submitting a future market change with filter with key type succeeds", testNewFutureMarketChangeSubmissionWithFilterWithKeyTypeSucceeds)
	t.Run("Submitting a future market change with filter without condition succeeds", testNewFutureMarketChangeSubmissionWithFilterWithoutConditionsSucceeds)
	t.Run("Submitting a future market change with filter without condition operator fails", testNewFutureMarketChangeSubmissionWithFilterWithoutConditionOperatorFails)
	t.Run("Submitting a future market change with filter with condition operator succeeds", testNewFutureMarketChangeSubmissionWithFilterWithConditionOperatorSucceeds)
	t.Run("Submitting a future market change with filter without condition value fails", testNewFutureMarketChangeSubmissionWithFilterWithoutConditionValueFails)
	t.Run("Submitting a future market change with filter with condition value succeeds", testNewFutureMarketChangeSubmissionWithFilterWithConditionValueSucceeds)
	t.Run("Submitting a future market change without oracle spec fails", testNewFutureMarketChangeSubmissionWithoutOracleSpecBindingFails)
	t.Run("Submitting a future market change with oracle spec binding succeeds", testNewFutureMarketChangeSubmissionWithOracleSpecBindingSucceeds)
	t.Run("Submitting a future market change without settlement price property fails", testNewFutureMarketChangeSubmissionWithoutSettlementPricePropertyFails)
	t.Run("Submitting a future market change with settlement price property succeeds", testNewFutureMarketChangeSubmissionWithSettlementPricePropertySucceeds)
	t.Run("Submitting a continuous trading market change without continuous trading mode fails", testNewContinuousTradingMarketChangeSubmissionWithoutContinuousTradingModeFails)
	t.Run("Submitting a continuous trading market change with continuous trading mode succeeds", testNewContinuousTradingMarketChangeSubmissionWithContinuousTradingModeSucceeds)
	t.Run("Submitting a discrete trading market change without discrete trading mode fails", testNewDiscreteTradingMarketChangeSubmissionWithoutDiscreteTradingModeFails)
	t.Run("Submitting a discrete trading market change with discrete trading mode succeeds", testNewDiscreteTradingMarketChangeSubmissionWithDiscreteTradingModeSucceeds)
	t.Run("Submitting a discrete trading market change without duration fails", testNewDiscreteTradingMarketChangeSubmissionWithWrongDurationFails)
	t.Run("Submitting a discrete trading market change with duration succeeds", testNewDiscreteTradingMarketChangeSubmissionWithRightDurationSucceeds)
	t.Run("Submitting a simple risk parameters change without simple risk parameters fails", testNewSimpleRiskParametersChangeSubmissionWithoutSimpleRiskParametersFails)
	t.Run("Submitting a simple risk parameters change with simple risk parameters succeeds", testNewSimpleRiskParametersChangeSubmissionWithSimpleRiskParametersSucceeds)
	t.Run("Submitting a simple risk parameters change with min move down fails", testNewSimpleRiskParametersChangeSubmissionWithPositiveMinMoveDownFails)
	t.Run("Submitting a simple risk parameters change with min move down succeeds", testNewSimpleRiskParametersChangeSubmissionWithNonPositiveMinMoveDownSucceeds)
	t.Run("Submitting a simple risk parameters change with max move up fails", testNewSimpleRiskParametersChangeSubmissionWithNegativeMaxMoveUpFails)
	t.Run("Submitting a simple risk parameters change with max move up succeeds", testNewSimpleRiskParametersChangeSubmissionWithNonNegativeMaxMoveUpSucceeds)
	t.Run("Submitting a simple risk parameters change with wrong probability of trading fails", testNewSimpleRiskParametersChangeSubmissionWithWrongProbabilityOfTradingFails)
	t.Run("Submitting a simple risk parameters change with right probability of trading succeeds", testNewSimpleRiskParametersChangeSubmissionWithRightProbabilityOfTradingSucceeds)
	t.Run("Submitting a log normal risk parameters change without log normal risk parameters fails", testNewLogNormalRiskParametersChangeSubmissionWithoutLogNormalRiskParametersFails)
	t.Run("Submitting a log normal risk parameters change with log normal risk parameters succeeds", testNewLogNormalRiskParametersChangeSubmissionWithLogNormalRiskParametersSucceeds)
	t.Run("Submitting a log normal risk parameters change with params fails", testNewLogNormalRiskParametersChangeSubmissionWithoutParamsFails)
	t.Run("Submitting a log normal risk parameters change with params succeeds", testNewLogNormalRiskParametersChangeSubmissionWithParamsSucceeds)
	t.Run("Submitting a new market without liquidity commitment fails", testNewMarketSubmissionWithoutLiquidityCommitmentFails)
	t.Run("Submitting a new market with liquidity commitment succeeds", testNewMarketSubmissionWithLiquidityCommitmentSucceeds)
	t.Run("Submitting a new market without commitment amount fails", testNewMarketSubmissionWithoutCommitmentAmountFails)
	t.Run("Submitting a new market with commitment amount succeeds", testNewMarketSubmissionWithCommitmentAmountSucceeds)
	t.Run("Submitting a new market without fee fails", testNewMarketSubmissionWithoutFeeFails)
	t.Run("Submitting a new market with wrong fee fails", testNewMarketSubmissionWithWrongFeeFails)
	t.Run("Submitting a new market with non-positive fee fails", testNewMarketSubmissionWithNonPositiveFeeFails)
	t.Run("Submitting a new market with right fee succeeds", testNewMarketSubmissionWithRightFeeSucceeds)
	t.Run("Submitting a new market with buy side and no orders fails", testNewMarketSubmissionWithBuySideAndNoOrdersFails)
	t.Run("Submitting a new market with buy side and orders succeeds", testNewMarketSubmissionWithBuySideAndOrdersSucceeds)
	t.Run("Submitting a new market with sell side and no orders fails", testNewMarketSubmissionWithSellSideAndNoOrdersFails)
	t.Run("Submitting a new market with sell side and orders succeeds", testNewMarketSubmissionWithSellSideAndOrdersSucceeds)
	t.Run("Submitting a new market with buy side and wrong order reference fails", testNewMarketSubmissionWithBuySideAndWrongOrderReferenceFails)
	t.Run("Submitting a new market with buy side and right order reference succeeds", testNewMarketSubmissionWithBuySideAndRightOrderReferenceSucceeds)
	t.Run("Submitting a new market with sell side and wrong order reference fails", testNewMarketSubmissionWithSellSideAndWrongOrderReferenceFails)
	t.Run("Submitting a new market with sell side and right order reference succeeds", testNewMarketSubmissionWithSellSideAndRightOrderReferenceSucceeds)
	t.Run("Submitting a new market with buy side and no order proportion fails", testNewMarketSubmissionWithBuySideAndNoOrderProportionFails)
	t.Run("Submitting a new market with buy side and order proportion fails", testNewMarketSubmissionWithBuySideAndOrderProportionSucceeds)
	t.Run("Submitting a new market with sell side and no order proportion fails", testNewMarketSubmissionWithSellSideAndNoOrderProportionFails)
	t.Run("Submitting a new market with sell side and order proportion fails", testNewMarketSubmissionWithSellSideAndOrderProportionSucceeds)
	t.Run("Submitting a new market with buy side and best ask reference fails", testNewMarketSubmissionWithBuySideAndBestAskReferenceFails)
	t.Run("Submitting a new market with buy side and best bid reference succeeds", testNewMarketSubmissionWithBuySideAndBestBidReferenceSucceeds)
	t.Run("Submitting a new market with buy side and best bid reference and positive offset fails", testNewMarketSubmissionWithBuySideAndBestBidReferenceAndPositiveOffsetFails)
	t.Run("Submitting a new market with buy side and best bid reference and non positive offset succeeds", testNewMarketSubmissionWithBuySideAndBestBidReferenceAndNonPositiveOffsetSucceeds)
	t.Run("Submitting a new market with buy side and mid reference and non-negative offset fails", testNewMarketSubmissionWithBuySideAndMidReferenceAndNonNegativeOffsetFails)
	t.Run("Submitting a new market with buy side and mid reference and negative offset succeeds", testNewMarketSubmissionWithBuySideAndMidReferenceAndNegativeOffsetSucceeds)
	t.Run("Submitting a new market with sell side and best bid reference fails", testNewMarketSubmissionWithSellSideAndBestBidReferenceFails)
	t.Run("Submitting a new market with sell side and best ask reference succeeds", testNewMarketSubmissionWithSellSideAndBestAskReferenceSucceeds)
	t.Run("Submitting a new market with sell side and best ask reference and negative offset fails", testNewMarketSubmissionWithSellSideAndBestAskReferenceAndNegativeOffsetFails)
	t.Run("Submitting a new market with sell side and best ask reference and non negative offset succeeds", testNewMarketSubmissionWithSellSideAndBestAskReferenceAndNonNegativeOffsetSucceeds)
	t.Run("Submitting a new market with sell side and mid reference and non-positive offset fails", testNewMarketSubmissionWithSellSideAndMidReferenceAndNonPositiveOffsetFails)
	t.Run("Submitting a new market with sell side and mid reference and positive offset succeeds", testNewMarketSubmissionWithSellSideAndMidReferenceAndPositiveOffsetSucceeds)
}

func testProposalSubmissionWithoutTermsFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{})

	assert.Contains(t, err.Get("proposal_submission.terms"), commands.ErrIsRequired)
}

func testNilProposalSubmissionFails(t *testing.T) {
	err := checkProposalSubmission(nil)

	assert.Contains(t, err.Get("proposal_submission"), commands.ErrIsRequired)
}

func testProposalSubmissionWithNonPositiveClosingTimestampFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 as closing timestamp",
			value: 0,
		}, {
			msg:   "with negative closing timestamp",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					ClosingTimestamp: tc.value,
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.closing_timestamp"), commands.ErrMustBePositive)
		})
	}
}

func testProposalSubmissionWithPositiveClosingTimestampSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp: RandomPositiveI64(),
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.closing_timestamp"), commands.ErrMustBePositive)
}

func testProposalSubmissionWithNonPositiveEnactmentTimestampFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 as closing timestamp",
			value: 0,
		}, {
			msg:   "with negative closing timestamp",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					EnactmentTimestamp: tc.value,
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.enactment_timestamp"), commands.ErrMustBePositive)
		})
	}
}

func testProposalSubmissionWithPositiveEnactmentTimestampSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			EnactmentTimestamp: RandomPositiveI64(),
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.enactment_timestamp"), commands.ErrMustBePositive)
}

func testProposalSubmissionWithNegativeValidationTimestampFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ValidationTimestamp: RandomNegativeI64(),
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.validation_timestamp"), commands.ErrMustBePositiveOrZero)
}

func testProposalSubmissionWithPositiveValidationTimestampSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ValidationTimestamp: RandomPositiveI64(),
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.validation_timestamp"), commands.ErrIsRequired)
}

func testProposalSubmissionWithClosingTimestampAfterEnactmentTimestampFails(t *testing.T) {
	closingTime := RandomPositiveI64()
	enactmentTime := RandomPositiveI64Before(closingTime)
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:   closingTime,
			EnactmentTimestamp: enactmentTime,
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.closing_timestamp"),
		errors.New("cannot be after enactment time"),
	)
}

func testProposalSubmissionWithClosingTimestampBeforeEnactmentTimestampSucceeds(t *testing.T) {
	enactmentTime := RandomPositiveI64()
	closingTime := RandomPositiveI64Before(enactmentTime)

	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:   closingTime,
			EnactmentTimestamp: enactmentTime,
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.closing_timestamp"),
		errors.New("cannot be after enactment time"),
	)
}

func testProposalSubmissionWithClosingTimestampAtEnactmentTimestampSucceeds(t *testing.T) {
	enactmentTime := RandomPositiveI64()

	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:   enactmentTime,
			EnactmentTimestamp: enactmentTime,
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.closing_timestamp"),
		errors.New("cannot be after enactment time"),
	)
}

func testProposalSubmissionWithValidationTimestampAfterClosingTimestampFails(t *testing.T) {
	validationTime := RandomPositiveI64()
	closingTime := RandomPositiveI64Before(validationTime)
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:    closingTime,
			ValidationTimestamp: validationTime,
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.validation_timestamp"),
		errors.New("cannot be after or equal to closing time"),
	)
}

func testProposalSubmissionWithValidationTimestampAtClosingTimestampFails(t *testing.T) {
	validationTime := RandomPositiveI64()

	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:    validationTime,
			ValidationTimestamp: validationTime,
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.validation_timestamp"),
		errors.New("cannot be after or equal to closing time"),
	)

}

func testProposalSubmissionWithValidationTimestampBeforeClosingTimestampSucceeds(t *testing.T) {
	closingTime := RandomPositiveI64()
	validationTime := RandomPositiveI64Before(closingTime)

	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			ClosingTimestamp:    closingTime,
			ValidationTimestamp: validationTime,
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.validation_timestamp"),
		errors.New("cannot be after or equal to closing time"),
	)
}

func testProposalSubmissionWithoutChangeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change"), commands.ErrIsRequired)
}

func testAssetChangeSubmissionWithoutNewsAssetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset"), commands.ErrIsRequired)
}

func testAssetChangeSubmissionWithoutChangesFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes"), commands.ErrIsRequired)
}

func testAssetChangeSubmissionWithoutSourceFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithoutBuiltInAssetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_BuiltinAsset{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithoutNameFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Name: "",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.name"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithNameSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Name: "My built-in asset",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.name"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithoutSymbolFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Symbol: "",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.symbol"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithSymbolSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Symbol: "My symbol",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.symbol"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithoutDecimalsFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Decimals: 0,
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.decimals"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithDecimalsSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Decimals: RandomPositiveU64(),
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.decimals"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithoutTotalSupplyFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						TotalSupply: "",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.total_supply"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithTotalSupplySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						TotalSupply: "10000",
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.total_supply"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.total_supply"), commands.ErrIsNotValidNumber)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.total_supply"), commands.ErrMustBePositive)
}

func testBuiltInAssetChangeSubmissionWithNaNTotalSupplyFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value string
		error error
	}{
		{
			msg:   "with not-a-number value",
			value: "hello",
			error: commands.ErrIsNotValidNumber,
		}, {
			msg:   "with value of 0",
			value: "0",
			error: commands.ErrMustBePositive,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewAsset{
						NewAsset: &typespb.NewAsset{
							Changes: &typespb.AssetDetails{
								TotalSupply: tc.value,
								Source: &typespb.AssetDetails_BuiltinAsset{
									BuiltinAsset: &typespb.BuiltinAsset{},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.total_supply"), tc.error)
		})
	}
}

func testBuiltInAssetChangeSubmissionWithoutMaxFaucetAmountMintFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{
								MaxFaucetAmountMint: "",
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint"), commands.ErrIsRequired)
}

func testBuiltInAssetChangeSubmissionWithMaxFaucetAmountMintSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_BuiltinAsset{
							BuiltinAsset: &typespb.BuiltinAsset{
								MaxFaucetAmountMint: "10000",
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint"), commands.ErrIsNotValidNumber)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint"), commands.ErrMustBePositive)
}

func testBuiltInAssetChangeSubmissionWithNaNMaxFaucetAmountMintFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value string
		error error
	}{
		{
			msg:   "with not-a-number value",
			value: "hello",
			error: commands.ErrIsNotValidNumber,
		}, {
			msg:   "with value of 0",
			value: "0",
			error: commands.ErrMustBePositive,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewAsset{
						NewAsset: &typespb.NewAsset{
							Changes: &typespb.AssetDetails{
								Source: &typespb.AssetDetails_BuiltinAsset{
									BuiltinAsset: &typespb.BuiltinAsset{
										MaxFaucetAmountMint: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.builtin_asset.max_faucet_amount_mint"), tc.error)
		})
	}
}

func testERC20AssetChangeSubmissionWithoutErc20AssetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_Erc20{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.erc20"), commands.ErrIsRequired)
}

func testErc20AssetChangeSubmissionWithoutContractAddressFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_Erc20{
							Erc20: &typespb.ERC20{
								ContractAddress: "",
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.erc20.contract_address"), commands.ErrIsRequired)
}

func testErc20AssetChangeSubmissionWithContractAddressSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewAsset{
				NewAsset: &typespb.NewAsset{
					Changes: &typespb.AssetDetails{
						Source: &typespb.AssetDetails_Erc20{
							Erc20: &typespb.ERC20{
								ContractAddress: "My address",
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.changes.source.erc20.contract_address"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithoutNetworkParameterFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.update_network_parameter"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithoutChangesFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &typespb.UpdateNetworkParameter{},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.update_network_parameter.changes"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithoutKeyFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &typespb.UpdateNetworkParameter{
					Changes: &typespb.NetworkParameter{
						Key: "",
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.update_network_parameter.changes.key"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithKeySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &typespb.UpdateNetworkParameter{
					Changes: &typespb.NetworkParameter{
						Key: "My key",
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.update_network_parameter.changes.key"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithoutValueFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &typespb.UpdateNetworkParameter{
					Changes: &typespb.NetworkParameter{
						Value: "",
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.update_network_parameter.changes.value"), commands.ErrIsRequired)
}

func testNetworkParameterChangeSubmissionWithValueSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_UpdateNetworkParameter{
				UpdateNetworkParameter: &typespb.UpdateNetworkParameter{
					Changes: &typespb.NetworkParameter{
						Value: "My value",
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.update_network_parameter.changes.value"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutNewMarketFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutChangesFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutDecimalPlacesFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.decimal_places"), commands.ErrMustBePositive)
}

func testNewMarketChangeSubmissionWithDecimalPlacesAboveOrEqualTo150Fails(t *testing.T) {
	testCases := []struct {
		msg   string
		value uint64
	}{
		{
			msg:   "equal to 150",
			value: 150,
		}, {
			msg:   "above 150",
			value: 1000,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								DecimalPlaces: tc.value,
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.decimal_places"), commands.ErrMustBeLessThan150)
		})
	}
}

func testNewMarketChangeSubmissionWithDecimalPlacesBelow150Succeeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						DecimalPlaces: RandomPositiveU64Before(150),
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.decimal_places"), commands.ErrMustBePositive)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.decimal_places"), commands.ErrMustBeLessThan150)
}

func testNewMarketChangeSubmissionWithoutLiquidityMonitoringSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithLiquidityMonitoringSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters"), commands.ErrIsRequired)
}

func testLiquidityMonitoringChangeSubmissionWithWrongTriggeringRatioFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with probability of -1",
			value: -1,
		}, {
			msg:   "with probability of 2",
			value: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
									TriggeringRatio: tc.value,
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.triggering_ratio"),
				errors.New("should be between 0 (inclusive) and 1 (inclusive)"))
		})
	}
}

func testLiquidityMonitoringChangeSubmissionWithRightTriggeringRatioSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with ratio of 0",
			value: 0,
		}, {
			msg:   "with ratio of 0.5",
			value: 0.5,
		}, {
			msg:   "with ratio of 1",
			value: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
									TriggeringRatio: tc.value,
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters"),
				errors.New("should be between 0 (inclusive) and 1 (inclusive)"))
		})
	}
}

func testLiquidityMonitoringChangeSubmissionWithoutTargetStakeParametersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters"), commands.ErrIsRequired)
}

func testLiquidityMonitoringChangeSubmissionWithTargetStakeParametersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
							TargetStakeParameters: &typespb.TargetStakeParameters{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters"), commands.ErrIsRequired)
}

func testLiquidityMonitoringChangeSubmissionWithNonPositiveTimeWindowFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with ratio of 0",
			value: 0,
		}, {
			msg:   "with ratio of -1",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
									TargetStakeParameters: &typespb.TargetStakeParameters{
										TimeWindow: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.time_window"), commands.ErrMustBePositive)
		})
	}
}

func testLiquidityMonitoringChangeSubmissionWithPositiveTimeWindowSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
							TargetStakeParameters: &typespb.TargetStakeParameters{
								TimeWindow: RandomPositiveI64(),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.time_window"), commands.ErrMustBePositive)
}

func testLiquidityMonitoringChangeSubmissionWithNonPositiveScalingFactorFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with ratio of 0",
			value: 0,
		}, {
			msg:   "with ratio of -1.5",
			value: -1.5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
									TargetStakeParameters: &typespb.TargetStakeParameters{
										ScalingFactor: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.scaling_factor"), commands.ErrMustBePositive)
		})
	}
}

func testLiquidityMonitoringChangeSubmissionWithPositiveScalingFactorSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						LiquidityMonitoringParameters: &typespb.LiquidityMonitoringParameters{
							TargetStakeParameters: &typespb.TargetStakeParameters{
								ScalingFactor: 1.5,
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.liquidity_monitoring_parameters.target_stake_parameters.scaling_factor"), commands.ErrMustBePositive)
}

func testPriceMonitoringChangeSubmissionWithoutTriggersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers"), commands.ErrIsRequired)
}

func testPriceMonitoringChangeSubmissionWithTriggersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{},
								{},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers"), commands.ErrIsRequired)
}

func testPriceMonitoringChangeSubmissionWithoutTriggerHorizonFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{},
								{},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.horizon"), commands.ErrMustBePositive)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.horizon"), commands.ErrMustBePositive)
}

func testPriceMonitoringChangeSubmissionWithTriggerHorizonSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{
									Horizon: RandomPositiveI64(),
								},
								{
									Horizon: RandomPositiveI64(),
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.horizon"), commands.ErrMustBePositive)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.horizon"), commands.ErrMustBePositive)
}

func testPriceMonitoringChangeSubmissionWithWrongTriggerProbabilityFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with probability of -1",
			value: -1,
		}, {
			msg:   "with probability of 0",
			value: 0,
		}, {
			msg:   "with probability of 1",
			value: 1,
		}, {
			msg:   "with probability of 2",
			value: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
									Triggers: []*typespb.PriceMonitoringTrigger{
										{
											Probability: tc.value,
										},
										{
											Probability: tc.value,
										},
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.probability"),
				errors.New("should be between 0 (exclusive) and 1 (exclusive)"))
			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.probability"),
				errors.New("should be between 0 (exclusive) and 1 (exclusive)"))
		})
	}
}

func testPriceMonitoringChangeSubmissionWithRightTriggerProbabilitySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{
									Probability: 0.01,
								},
								{
									Probability: 0.9,
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.probability"),
		errors.New("should be between 0 (exclusive) and 1 (exclusive)"))
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.probability"),
		errors.New("should be between 0 (exclusive) and 1 (exclusive)"))
}

func testPriceMonitoringChangeSubmissionWithoutTriggerAuctionExtensionFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{},
								{},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.auction_extension"), commands.ErrMustBePositive)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.auction_extension"), commands.ErrMustBePositive)
}

func testPriceMonitoringChangeSubmissionWithTriggerAuctionExtensionSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{
							Triggers: []*typespb.PriceMonitoringTrigger{
								{
									AuctionExtension: RandomPositiveI64(),
								},
								{
									AuctionExtension: RandomPositiveI64(),
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.0.auction_extension"), commands.ErrMustBePositive)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters.triggers.1.auction_extension"), commands.ErrMustBePositive)
}

func testNewMarketChangeSubmissionWithoutPriceMonitoringSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithPriceMonitoringSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						PriceMonitoringParameters: &typespb.PriceMonitoringParameters{},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.price_monitoring_parameters"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutInstrumentNameFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Name: "",
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.name"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithInstrumentNameSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Name: "My name",
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.name"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutInstrumentCodeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Code: "",
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.code"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithInstrumentCodeSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Code: "My code",
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.code"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithoutProductFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product"), commands.ErrIsRequired)
}

func testNewMarketChangeSubmissionWithProductSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutFutureFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFutureSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutSettlementAssetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									SettlementAsset: "",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.settlement_asset"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithSettlementAssetSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									SettlementAsset: "BTC",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.settlement_asset"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutQuoteNameFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									QuoteName: "",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.quote_name"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithQuoteNameSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									QuoteName: "BTC",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.quote_name"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutMaturityFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									Maturity: "",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithMaturitySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									Maturity: "2020-10-22T12:00:00Z",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithWrongMaturityDateFormatFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									Maturity: "2020/10/25",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity"), commands.ErrMustBeValidDate)
}

func testNewFutureMarketChangeSubmissionWithRightMaturityDateFormatSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									Maturity: "2020-10-22T12:00:00Z",
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.maturity"), commands.ErrMustBeValidDate)
}

func testNewFutureMarketChangeSubmissionWithoutOracleSpecFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithOracleSpecSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutPubKeysFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										PubKeys: []string{},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithWrongPubKeysFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value []string
	}{
		{
			msg:   "with empty pub-keys",
			value: []string{"0xDEADBEEF", ""},
		}, {
			msg:   "with blank pub-keys",
			value: []string{"0xDEADBEEF", " "},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								Instrument: &typespb.InstrumentConfiguration{
									Product: &typespb.InstrumentConfiguration_Future{
										Future: &typespb.FutureProduct{
											OracleSpec: &oraclespb.OracleSpecConfiguration{
												PubKeys: tc.value,
											},
										},
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys.1"), commands.ErrIsNotValid)
		})
	}
}

func testNewFutureMarketChangeSubmissionWithPubKeysSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										PubKeys: []string{"0xDEADBEEF", "0xCAFEDUDE"},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys.0"), commands.ErrIsNotValid)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.pub_keys.1"), commands.ErrIsNotValid)
}

func testNewFutureMarketChangeSubmissionWithoutFiltersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFiltersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutKeyFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{}, {},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key"), commands.ErrIsNotValid)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key"), commands.ErrIsNotValid)
}

func testNewFutureMarketChangeSubmissionWithFilterWithKeySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{},
											}, {
												Key: &oraclespb.PropertyKey{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key"), commands.ErrIsNotValid)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key"), commands.ErrIsNotValid)
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutKeyNameFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{
													Name: "",
												},
											}, {
												Key: &oraclespb.PropertyKey{
													Name: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key.name"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key.name"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithKeyNameSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{
													Name: "key1",
												},
											}, {
												Key: &oraclespb.PropertyKey{
													Name: "key2",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key.name"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key.name"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutKeyTypeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Key: &oraclespb.PropertyKey{
													Type: oraclespb.PropertyKey_TYPE_UNSPECIFIED,
												},
											}, {
												Key: &oraclespb.PropertyKey{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key.type"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key.type"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithKeyTypeSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value oraclespb.PropertyKey_Type
	}{
		{
			msg:   "with EMPTY",
			value: oraclespb.PropertyKey_TYPE_EMPTY,
		}, {
			msg:   "with INTEGER",
			value: oraclespb.PropertyKey_TYPE_INTEGER,
		}, {
			msg:   "with STRING",
			value: oraclespb.PropertyKey_TYPE_STRING,
		}, {
			msg:   "with BOOLEAN",
			value: oraclespb.PropertyKey_TYPE_BOOLEAN,
		}, {
			msg:   "with DECIMAL",
			value: oraclespb.PropertyKey_TYPE_DECIMAL,
		}, {
			msg:   "with TIMESTAMP",
			value: oraclespb.PropertyKey_TYPE_TIMESTAMP,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								Instrument: &typespb.InstrumentConfiguration{
									Product: &typespb.InstrumentConfiguration_Future{
										Future: &typespb.FutureProduct{
											OracleSpec: &oraclespb.OracleSpecConfiguration{
												Filters: []*oraclespb.Filter{
													{
														Key: &oraclespb.PropertyKey{
															Type: tc.value,
														},
													}, {
														Key: &oraclespb.PropertyKey{
															Type: tc.value,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.key.type"), commands.ErrIsRequired)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.key.type"), commands.ErrIsRequired)
		})
	}
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutConditionsSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Conditions: []*oraclespb.Condition{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutConditionOperatorFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Conditions: []*oraclespb.Condition{
													{
														Operator: oraclespb.Condition_OPERATOR_UNSPECIFIED,
													},
													{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.0.operator"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.1.operator"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithConditionOperatorSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value oraclespb.Condition_Operator
	}{
		{
			msg:   "with EQUALS",
			value: oraclespb.Condition_OPERATOR_EQUALS,
		}, {
			msg:   "with GREATER_THAN",
			value: oraclespb.Condition_OPERATOR_GREATER_THAN,
		}, {
			msg:   "with GREATER_THAN_OR_EQUAL",
			value: oraclespb.Condition_OPERATOR_GREATER_THAN_OR_EQUAL,
		}, {
			msg:   "with LESS_THAN",
			value: oraclespb.Condition_OPERATOR_LESS_THAN,
		}, {
			msg:   "with LESS_THAN_OR_EQUAL",
			value: oraclespb.Condition_OPERATOR_LESS_THAN_OR_EQUAL,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								Instrument: &typespb.InstrumentConfiguration{
									Product: &typespb.InstrumentConfiguration_Future{
										Future: &typespb.FutureProduct{
											OracleSpec: &oraclespb.OracleSpecConfiguration{
												Filters: []*oraclespb.Filter{
													{
														Conditions: []*oraclespb.Condition{
															{
																Operator: tc.value,
															},
															{
																Operator: tc.value,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.0.operator"), commands.ErrIsRequired)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.1.conditions.0.operator"), commands.ErrIsRequired)
		})
	}
}

func testNewFutureMarketChangeSubmissionWithFilterWithoutConditionValueFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Conditions: []*oraclespb.Condition{
													{
														Value: "",
													},
													{
														Value: "",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.0.value"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.1.value"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithFilterWithConditionValueSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpec: &oraclespb.OracleSpecConfiguration{
										Filters: []*oraclespb.Filter{
											{
												Conditions: []*oraclespb.Condition{
													{
														Value: "value 1",
													},
													{
														Value: "value 2",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.0.value"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec.filters.0.conditions.1.value"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutOracleSpecBindingFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithOracleSpecBindingSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpecBinding: &typespb.OracleSpecToFutureBinding{},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithoutSettlementPricePropertyFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpecBinding: &typespb.OracleSpecToFutureBinding{
										SettlementPriceProperty: "",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding.settlement_price_property"), commands.ErrIsRequired)
}

func testNewFutureMarketChangeSubmissionWithSettlementPricePropertySucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						Instrument: &typespb.InstrumentConfiguration{
							Product: &typespb.InstrumentConfiguration_Future{
								Future: &typespb.FutureProduct{
									OracleSpecBinding: &typespb.OracleSpecToFutureBinding{
										SettlementPriceProperty: "My property",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.instrument.product.future.oracle_spec_binding.settlement_price_property"), commands.ErrIsRequired)
}

func testNewContinuousTradingMarketChangeSubmissionWithoutContinuousTradingModeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						TradingMode: &typespb.NewMarketConfiguration_Continuous{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.continuous"), commands.ErrIsRequired)
}

func testNewContinuousTradingMarketChangeSubmissionWithContinuousTradingModeSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						TradingMode: &typespb.NewMarketConfiguration_Continuous{
							Continuous: &typespb.ContinuousTrading{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.continuous"), commands.ErrIsRequired)
}

func testNewDiscreteTradingMarketChangeSubmissionWithoutDiscreteTradingModeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						TradingMode: &typespb.NewMarketConfiguration_Discrete{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.discrete"), commands.ErrIsRequired)
}

func testNewDiscreteTradingMarketChangeSubmissionWithDiscreteTradingModeSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						TradingMode: &typespb.NewMarketConfiguration_Discrete{
							Discrete: &typespb.DiscreteTrading{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.discrete"), commands.ErrIsRequired)
}

func testNewDiscreteTradingMarketChangeSubmissionWithWrongDurationFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with duration of 0",
			value: 0,
		}, {
			msg:   "with duration under 0",
			value: -1,
		}, {
			msg:   "with duration of 2592000000000000",
			value: 2592000000000000,
		}, {
			msg:   "with duration above 2592000000000000",
			value: 2592000000000001,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								TradingMode: &typespb.NewMarketConfiguration_Discrete{
									Discrete: &typespb.DiscreteTrading{
										DurationNs: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.discrete.duration_ns"),
				errors.New("should be between 0 (excluded) and 2592000000000000 (excluded)"))
		})
	}
}

func testNewDiscreteTradingMarketChangeSubmissionWithRightDurationSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						TradingMode: &typespb.NewMarketConfiguration_Discrete{
							Discrete: &typespb.DiscreteTrading{
								DurationNs: RandomPositiveI64Before(2592000000000000 - 1),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.trading_mode.discrete.duration_ns"),
		errors.New("should be between 0 (excluded) and 2592000000000000 (excluded)"))
}

func testNewSimpleRiskParametersChangeSubmissionWithoutSimpleRiskParametersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_Simple{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple"), commands.ErrIsRequired)
}

func testNewSimpleRiskParametersChangeSubmissionWithSimpleRiskParametersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_Simple{
							Simple: &typespb.SimpleModelParams{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple"), commands.ErrIsRequired)
}

func testNewSimpleRiskParametersChangeSubmissionWithPositiveMinMoveDownFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_Simple{
							Simple: &typespb.SimpleModelParams{
								MinMoveDown: 1,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.min_move_down"), commands.ErrMustBeNegativeOrZero)
}

func testNewSimpleRiskParametersChangeSubmissionWithNonPositiveMinMoveDownSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with min move down of 0",
			value: 0,
		}, {
			msg:   "with min move down of -1",
			value: -1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								RiskParameters: &typespb.NewMarketConfiguration_Simple{
									Simple: &typespb.SimpleModelParams{
										MinMoveDown: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.min_move_down"), commands.ErrMustBeNegativeOrZero)
		})
	}
}

func testNewSimpleRiskParametersChangeSubmissionWithNegativeMaxMoveUpFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_Simple{
							Simple: &typespb.SimpleModelParams{
								MaxMoveUp: -1,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.max_move_up"), commands.ErrMustBePositiveOrZero)
}

func testNewSimpleRiskParametersChangeSubmissionWithNonNegativeMaxMoveUpSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with max move up of 0",
			value: 0,
		}, {
			msg:   "with max move up of 1",
			value: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								RiskParameters: &typespb.NewMarketConfiguration_Simple{
									Simple: &typespb.SimpleModelParams{
										MaxMoveUp: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.max_move_up"), commands.ErrMustBePositiveOrZero)
		})
	}
}

func testNewSimpleRiskParametersChangeSubmissionWithWrongProbabilityOfTradingFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with probability of trading of -1",
			value: -1,
		}, {
			msg:   "with probability of trading of 2",
			value: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								RiskParameters: &typespb.NewMarketConfiguration_Simple{
									Simple: &typespb.SimpleModelParams{
										ProbabilityOfTrading: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.probability_of_trading"),
				errors.New("should be between 0 (inclusive) and 1 (inclusive)"))
		})
	}
}

func testNewSimpleRiskParametersChangeSubmissionWithRightProbabilityOfTradingSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value float64
	}{
		{
			msg:   "with probability of trading of 0",
			value: 0,
		}, {
			msg:   "with probability of trading of 1",
			value: 1,
		}, {
			msg:   "with probability of trading of 0.5",
			value: 0.5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							Changes: &typespb.NewMarketConfiguration{
								RiskParameters: &typespb.NewMarketConfiguration_Simple{
									Simple: &typespb.SimpleModelParams{
										ProbabilityOfTrading: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.simple.probability_of_trading"),
				errors.New("should be between 0 (inclusive) and 1 (inclusive)"))
		})
	}
}

func testNewLogNormalRiskParametersChangeSubmissionWithoutLogNormalRiskParametersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_LogNormal{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal"), commands.ErrIsRequired)
}

func testNewLogNormalRiskParametersChangeSubmissionWithLogNormalRiskParametersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_LogNormal{
							LogNormal: &typespb.LogNormalRiskModel{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal"), commands.ErrIsRequired)
}

func testNewLogNormalRiskParametersChangeSubmissionWithoutParamsFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_LogNormal{
							LogNormal: &typespb.LogNormalRiskModel{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal.params"), commands.ErrIsRequired)
}

func testNewLogNormalRiskParametersChangeSubmissionWithParamsSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					Changes: &typespb.NewMarketConfiguration{
						RiskParameters: &typespb.NewMarketConfiguration_LogNormal{
							LogNormal: &typespb.LogNormalRiskModel{
								Params: &typespb.LogNormalModelParams{},
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.changes.risk_parameters.log_normal.params"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithoutLiquidityCommitmentFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithLiquidityCommitmentSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithoutCommitmentAmountFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.commitment_amount"), commands.ErrMustBePositive)
}

func testNewMarketSubmissionWithCommitmentAmountSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						CommitmentAmount: RandomPositiveU64(),
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.commitment_amount"), commands.ErrMustBePositive)
}

func testNewMarketSubmissionWithoutFeeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Fee: "",
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.fee"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithWrongFeeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Fee: "no a valid fee",
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.fee"), commands.ErrIsNotValidNumber)
}

func testNewMarketSubmissionWithNonPositiveFeeFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Fee: "-1",
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.fee"), commands.ErrMustBePositiveOrZero)
}

func testNewMarketSubmissionWithRightFeeSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value uint64
	}{
		{
			msg:   "with 0 fee",
			value: 0,
		}, {
			msg:   "with positive fee",
			value: RandomPositiveU64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Fee: strconv.FormatUint(tc.value, 10),
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.fee"), commands.ErrIsNotValidNumber)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_market.liquidity_commitment.fee"), commands.ErrMustBePositive)
		})
	}
}

func testNewMarketSubmissionWithBuySideAndNoOrdersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithBuySideAndOrdersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{},
							{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithSellSideAndNoOrdersFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithSellSideAndOrdersSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{},
							{},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithBuySideAndWrongOrderReferenceFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Reference: 42,
							},
							{
								Reference: 21,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.0"), commands.ErrIsNotValid)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.1"), commands.ErrIsNotValid)
}

func testNewMarketSubmissionWithBuySideAndRightOrderReferenceSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.PeggedReference
	}{
		{
			msg:   "with MID",
			value: typespb.PeggedReference_PEGGED_REFERENCE_MID,
		}, {
			msg:   "with BEST_BID",
			value: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
		},
		{
			msg:   "with BEST_ASK",
			value: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Buys: []*typespb.LiquidityOrder{
									{
										Reference: tc.value,
									},
									{
										Reference: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.0"), commands.ErrIsNotValid)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.1"), commands.ErrIsNotValid)
		})
	}
}

func testNewMarketSubmissionWithSellSideAndWrongOrderReferenceFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Reference: 42,
							},
							{
								Reference: 21,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.0"), commands.ErrIsNotValid)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.1"), commands.ErrIsNotValid)
}

func testNewMarketSubmissionWithSellSideAndRightOrderReferenceSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.PeggedReference
	}{
		{
			msg:   "with MID",
			value: typespb.PeggedReference_PEGGED_REFERENCE_MID,
		}, {
			msg:   "with BEST_BID",
			value: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
		},
		{
			msg:   "with BEST_ASK",
			value: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Sells: []*typespb.LiquidityOrder{
									{
										Reference: tc.value,
									},
									{
										Reference: tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.0"), commands.ErrIsNotValid)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.1"), commands.ErrIsNotValid)
		})
	}
}

func testNewMarketSubmissionWithBuySideAndNoOrderProportionFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{},
							{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.0"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.1"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithBuySideAndOrderProportionSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Proportion: RandomPositiveU32(),
							},
							{
								Proportion: RandomPositiveU32(),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.0"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.1"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithSellSideAndNoOrderProportionFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{},
							{},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.proportion.0"), commands.ErrIsRequired)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.proportion.1"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithSellSideAndOrderProportionSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Proportion: RandomPositiveU32(),
							},
							{
								Proportion: RandomPositiveU32(),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.0"), commands.ErrIsRequired)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.proportion.1"), commands.ErrIsRequired)
}

func testNewMarketSubmissionWithBuySideAndBestAskReferenceFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.0"),
		errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.1"),
		errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
}

func testNewMarketSubmissionWithBuySideAndBestBidReferenceSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.0"),
		errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.reference.1"),
		errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
}

func testNewMarketSubmissionWithBuySideAndBestBidReferenceAndPositiveOffsetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
								Offset:    RandomPositiveI64(),
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
								Offset:    RandomPositiveI64(),
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.0"), commands.ErrMustBeNegativeOrZero)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.1"), commands.ErrMustBeNegativeOrZero)
}

func testNewMarketSubmissionWithBuySideAndBestBidReferenceAndNonPositiveOffsetSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 offset",
			value: 0,
		}, {
			msg:   "with negative offset",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Buys: []*typespb.LiquidityOrder{
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
										Offset:    RandomNegativeI64(),
									},
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
										Offset:    RandomNegativeI64(),
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.0"), commands.ErrMustBeNegativeOrZero)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.1"), commands.ErrMustBeNegativeOrZero)
		})
	}
}

func testNewMarketSubmissionWithBuySideAndMidReferenceAndNonNegativeOffsetFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 offset",
			value: 0,
		}, {
			msg:   "with positive offset",
			value: RandomPositiveI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Buys: []*typespb.LiquidityOrder{
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
										Offset:    tc.value,
									},
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
										Offset:    tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.0"), commands.ErrMustBeNegative)
			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.1"), commands.ErrMustBeNegative)
		})
	}
}

func testNewMarketSubmissionWithBuySideAndMidReferenceAndNegativeOffsetSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Buys: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
								Offset:    RandomNegativeI64(),
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
								Offset:    RandomNegativeI64(),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.0"), commands.ErrMustBeNegative)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.buys.offset.1"), commands.ErrMustBeNegative)
}

func testNewMarketSubmissionWithSellSideAndBestBidReferenceFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.0"),
		errors.New("cannot have a reference of type BEST_BID when on SELL side"))
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.1"),
		errors.New("cannot have a reference of type BEST_BID when on SELL side"))
}

func testNewMarketSubmissionWithSellSideAndBestAskReferenceSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.0"),
		errors.New("cannot have a reference of type BEST_BID when on SELL side"))
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.reference.1"),
		errors.New("cannot have a reference of type BEST_BID when on SELL side"))
}

func testNewMarketSubmissionWithSellSideAndBestAskReferenceAndNegativeOffsetFails(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
								Offset:    RandomNegativeI64(),
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
								Offset:    RandomNegativeI64(),
							},
						},
					},
				},
			},
		},
	})

	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.0"), commands.ErrMustBePositiveOrZero)
	assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.1"), commands.ErrMustBePositiveOrZero)
}

func testNewMarketSubmissionWithSellSideAndBestAskReferenceAndNonNegativeOffsetSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 offset",
			value: 0,
		}, {
			msg:   "with positive offset",
			value: RandomPositiveI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Sells: []*typespb.LiquidityOrder{
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
										Offset:    tc.value,
									},
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
										Offset:    tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.0"), commands.ErrMustBePositiveOrZero)
			assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.1"), commands.ErrMustBePositiveOrZero)
		})
	}
}

func testNewMarketSubmissionWithSellSideAndMidReferenceAndNonPositiveOffsetFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 offset",
			value: 0,
		}, {
			msg:   "with negative offset",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkProposalSubmission(&commandspb.ProposalSubmission{
				Terms: &typespb.ProposalTerms{
					Change: &typespb.ProposalTerms_NewMarket{
						NewMarket: &typespb.NewMarket{
							LiquidityCommitment: &typespb.NewMarketCommitment{
								Sells: []*typespb.LiquidityOrder{
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
										Offset:    tc.value,
									},
									{
										Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
										Offset:    tc.value,
									},
								},
							},
						},
					},
				},
			})

			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.0"), commands.ErrMustBePositive)
			assert.Contains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.1"), commands.ErrMustBePositive)
		})
	}
}

func testNewMarketSubmissionWithSellSideAndMidReferenceAndPositiveOffsetSucceeds(t *testing.T) {
	err := checkProposalSubmission(&commandspb.ProposalSubmission{
		Terms: &typespb.ProposalTerms{
			Change: &typespb.ProposalTerms_NewMarket{
				NewMarket: &typespb.NewMarket{
					LiquidityCommitment: &typespb.NewMarketCommitment{
						Sells: []*typespb.LiquidityOrder{
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
								Offset:    RandomPositiveI64(),
							},
							{
								Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
								Offset:    RandomPositiveI64(),
							},
						},
					},
				},
			},
		},
	})

	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.0"), commands.ErrMustBePositive)
	assert.NotContains(t, err.Get("proposal_submission.terms.change.new_asset.liquidity_commitment.sells.offset.1"), commands.ErrMustBePositive)
}

func checkProposalSubmission(cmd *commandspb.ProposalSubmission) commands.Errors {
	err := commands.CheckProposalSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
