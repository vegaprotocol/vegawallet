package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestNilLiquidityProvisionSubmissionFails(t *testing.T) {
	err := commands.CheckLiquidityProvisionSubmission(nil)

	assert.Error(t, err)
}

func TestLiquidityProvisionSubmission(t *testing.T) {
	var cases = []struct {
		lp        commandspb.LiquidityProvisionSubmission
		errString string
	}{
		{
			// this is a valid cancellation.
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 0,
				MarketId:         "okmarketid",
			},
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "abcd",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.fee (is not a valid value)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "-1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.fee (must be positive)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.market_id (is required)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.sells.0.proportion (order in shape without a proportion)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10},
				},
			},
			errString: "liquidity_provision_submission.buys.0.proportion (order in shape without a proportion)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Proportion: 1}, // no offset is ok
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Proportion: 1}, // no offset is ok
				},
			},
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{},
			},
			errString: "liquidity_provision_submission.buys (empty shape)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells:            []*typespb.LiquidityOrder{},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.sells (empty shape)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: 10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: -10, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.buys.0.reference (order in buy side shape with best ask price reference), liquidity_provision_submission.sells.0.offset (order in sell side shape with best bid price reference)",
		},
		{
			lp: commandspb.LiquidityProvisionSubmission{
				CommitmentAmount: 100,
				Fee:              "0.1",
				MarketId:         "okmarketid",
				Sells: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID, Offset: 0, Proportion: 1},
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID, Offset: 0, Proportion: 1},
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK, Offset: -10, Proportion: 1},
				},
				Buys: []*typespb.LiquidityOrder{
					{Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID, Offset: 0, Proportion: 1},
				},
			},
			errString: "liquidity_provision_submission.buys.0.offset (order in buy side shape offset must be < 0), liquidity_provision_submission.sells.0.offset (order in sell shape offset must be > 0), liquidity_provision_submission.sells.1.offset (order in sell side shape with best bid price reference), liquidity_provision_submission.sells.2.offset (order in sell shape offset must be >= 0)",
		},
		{
			// this is considered as an invalid cancellation, as a market id is missing.
			lp:        commandspb.LiquidityProvisionSubmission{},
			errString: "liquidity_provision_submission.market_id (is required)",
		},
	}

	for _, c := range cases {
		err := commands.CheckLiquidityProvisionSubmission(&c.lp)
		if len(c.errString) <= 0 {
			assert.NoError(t, err)
			continue
		}

		assert.Error(t, err)
		assert.EqualError(t, err, c.errString)
	}
}
