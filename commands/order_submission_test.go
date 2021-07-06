package commands_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckOrderSubmission(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilOrderSubmissionFails)
	t.Run("Submitting an empty order fails", testEmptyOrderSubmissionFails)
	t.Run("Submitting an order without market ID fails", testOrderSubmissionWithoutMarketIDFails)
	t.Run("Submitting an order with unspecified side fails", testOrderSubmissionWithUnspecifiedSideFails)
	t.Run("Submitting an order with undefined side fails", testOrderSubmissionWithUndefinedSideFails)
	t.Run("Submitting an order with unspecified type fails", testOrderSubmissionWithUnspecifiedTypeFails)
	t.Run("Submitting an order with undefined type fails", testOrderSubmissionWithUndefinedTypeFails)
	t.Run("Submitting an order with NETWORK type fails", testOrderSubmissionWithNetworkTypeFails)
	t.Run("Submitting an order with undefined time in force fails", testOrderSubmissionWithUndefinedTimeInForceFails)
	t.Run("Submitting an order with unspecified time in force fails", testOrderSubmissionWithUnspecifiedTimeInForceFails)
	t.Run("Submitting an order with non-positive size fails", testOrderSubmissionWithNonPositiveSizeFails)
	t.Run("Submitting an order with GTT and non-positive expiration date fails", testOrderSubmissionWithGTTAndNonPositiveExpirationDateFails)
	t.Run("Submitting an order without GTT and expiration date fails", testOrderSubmissionWithoutGTTAndExpirationDateFails)
	t.Run("Submitting an order with MARKET type and price fails", testOrderSubmissionWithMarketTypeAndPriceFails)
	t.Run("Submitting an order with MARKET type and wrong time in force fails", testOrderSubmissionWithMarketTypeAndWrongTimeInForceFails)
	t.Run("Submitting an order with LIMIT type and no price fails", testOrderSubmissionWithLimitTypeAndNoPriceFails)
	t.Run("Submitting a pegged order with LIMIT type and no price succeeds", testPeggedOrderSubmissionWithLimitTypeAndNoPriceSucceeds)
	t.Run("Submitting a pegged order with undefined time in force fails", testPeggedOrderSubmissionWithUndefinedReferenceFails)
	t.Run("Submitting a pegged order with unspecified time in force fails", testPeggedOrderSubmissionWithUnspecifiedReferenceFails)
	t.Run("Submitting a pegged order without LIMIT type fails", testPeggedOrderSubmissionWithoutLimitTypeFails)
	t.Run("Submitting a pegged order with LIMIT type succeeds", testPeggedOrderSubmissionWithLimitTypeSucceeds)
	t.Run("Submitting a pegged order with wrong time in force fails", testPeggedOrderSubmissionWithWrongTimeInForceFails)
	t.Run("Submitting a pegged order with right time in force succeeds", testPeggedOrderSubmissionWithRightTimeInForceSucceeds)
	t.Run("Submitting a pegged order with side buy and best ask reference fails", testPeggedOrderSubmissionWithSideBuyAndBestAskReferenceFails)
	t.Run("Submitting a pegged order with side buy and best bid reference succeeds", testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceSucceeds)
	t.Run("Submitting a pegged order with side buy and best bid reference and positive offset fails", testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceAndPositiveOffsetFails)
	t.Run("Submitting a pegged order with side buy and best bid reference and non positive offset succeeds", testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceAndNonPositiveOffsetSucceeds)
	t.Run("Submitting a pegged order with side buy and mid reference and non-negative offset fails", testPeggedOrderSubmissionWithSideBuyAndMidReferenceAndNonNegativeOffsetFails)
	t.Run("Submitting a pegged order with side buy and mid reference and negative offset succeeds", testPeggedOrderSubmissionWithSideBuyAndMidReferenceAndNegativeOffsetSucceeds)
	t.Run("Submitting a pegged order with side sell and best bid reference fails", testPeggedOrderSubmissionWithSideSellAndBestBidReferenceFails)
	t.Run("Submitting a pegged order with side sell and best ask reference succeeds", testPeggedOrderSubmissionWithSideSellAndBestAskReferenceSucceeds)
	t.Run("Submitting a pegged order with side sell and best ask reference and negative offset fails", testPeggedOrderSubmissionWithSideSellAndBestAskReferenceAndNegativeOffsetFails)
	t.Run("Submitting a pegged order with side sell and best ask reference and non negative offset succeeds", testPeggedOrderSubmissionWithSideSellAndBestAskReferenceAndNonNegativeOffsetSucceeds)
	t.Run("Submitting a pegged order with side sell and mid reference and non-positive offset fails", testPeggedOrderSubmissionWithSideSellAndMidReferenceAndNonPositiveOffsetFails)
	t.Run("Submitting a pegged order with side sell and mid reference and positive offset succeeds", testPeggedOrderSubmissionWithSideSellAndMidReferenceAndPositiveOffsetSucceeds)
}

func testEmptyOrderSubmissionFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{})

	assert.Error(t, err)
}

func testNilOrderSubmissionFails(t *testing.T) {
	err := checkOrderSubmission(nil)

	assert.Contains(t, err.Get("order_submission"), commands.ErrIsRequired)
}

func testOrderSubmissionWithoutMarketIDFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		MarketId: "",
	})

	assert.Contains(t, err.Get("order_submission.market_id"), commands.ErrIsRequired)
}

func testOrderSubmissionWithUnspecifiedSideFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_UNSPECIFIED,
	})

	assert.Contains(t, err.Get("order_submission.side"), commands.ErrIsRequired)
}

func testOrderSubmissionWithUndefinedSideFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side(-42),
	})

	assert.Contains(t, err.Get("order_submission.side"), commands.ErrIsNotValid)
}

func testOrderSubmissionWithUnspecifiedTypeFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type: typespb.Order_TYPE_UNSPECIFIED,
	})

	assert.Contains(t, err.Get("order_submission.type"), commands.ErrIsRequired)
}

func testOrderSubmissionWithUndefinedTypeFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type: typespb.Order_Type(-42),
	})

	assert.Contains(t, err.Get("order_submission.type"), commands.ErrIsNotValid)
}

func testOrderSubmissionWithNetworkTypeFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type: typespb.Order_TYPE_NETWORK,
	})

	assert.Contains(t, err.Get("order_submission.type"), commands.ErrIsUnauthorised)
}

func testOrderSubmissionWithUnspecifiedTimeInForceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		TimeInForce: typespb.Order_TIME_IN_FORCE_UNSPECIFIED,
	})

	assert.Contains(t, err.Get("order_submission.time_in_force"), commands.ErrIsRequired)
}

func testOrderSubmissionWithUndefinedTimeInForceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		TimeInForce: typespb.Order_TimeInForce(-42),
	})

	assert.Contains(t, err.Get("order_submission.time_in_force"), commands.ErrIsNotValid)
}

func testOrderSubmissionWithNonPositiveSizeFails(t *testing.T) {
	// FIXME(big int) doesn't test negative numbers since it's an unsigned int
	// 	but that will definitely be needed when moving to big int.
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Size: 0,
	})

	assert.Contains(t, err.Get("order_submission.size"), commands.ErrMustBePositive)
}

func testOrderSubmissionWithGTTAndNonPositiveExpirationDateFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value int64
	}{
		{
			msg:   "with 0 as expiration date",
			value: 0,
		}, {
			msg:   "with negative expiration date",
			value: RandomNegativeI64(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				TimeInForce: typespb.Order_TIME_IN_FORCE_GTT,
				ExpiresAt:   tc.value,
			})

			assert.Contains(t, err.Get("order_submission.expires_at"), commands.ErrMustBePositive)
		})
	}
}

func testOrderSubmissionWithoutGTTAndExpirationDateFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.Order_TimeInForce
	}{
		{
			msg:   "with GTC",
			value: typespb.Order_TIME_IN_FORCE_GTC,
		}, {
			msg:   "with IOC",
			value: typespb.Order_TIME_IN_FORCE_IOC,
		}, {
			msg:   "with FOK",
			value: typespb.Order_TIME_IN_FORCE_FOK,
		}, {
			msg:   "with GFA",
			value: typespb.Order_TIME_IN_FORCE_GFA,
		}, {
			msg:   "with GFN",
			value: typespb.Order_TIME_IN_FORCE_GFN,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				TimeInForce: tc.value,
				ExpiresAt:   RandomI64(),
			})

			assert.Contains(t, err.Get("order_submission.expires_at"), errors.New("is only available when the time in force is of type GTT"))
		})
	}
}

func testOrderSubmissionWithMarketTypeAndPriceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type:  typespb.Order_TYPE_MARKET,
		Price: RandomPositiveU64(),
	})

	assert.Contains(t, err.Get("order_submission.price"), errors.New("is unavailable when the order is of type MARKET"))
}

func testOrderSubmissionWithMarketTypeAndWrongTimeInForceFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.Order_TimeInForce
	}{
		{
			msg:   "with GTC",
			value: typespb.Order_TIME_IN_FORCE_GTC,
		}, {
			msg:   "with GTT",
			value: typespb.Order_TIME_IN_FORCE_GTT,
		}, {
			msg:   "with GFA",
			value: typespb.Order_TIME_IN_FORCE_GFA,
		}, {
			msg:   "with GFN",
			value: typespb.Order_TIME_IN_FORCE_GFN,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Type:        typespb.Order_TYPE_MARKET,
				TimeInForce: tc.value,
			})

			assert.Contains(t, err.Get("order_submission.time_in_force"), errors.New("is expected to be of type FOK or IOC when order is of type MARKET"))
		})
	}
}

func testOrderSubmissionWithLimitTypeAndNoPriceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type: typespb.Order_TYPE_LIMIT,
	})

	assert.Contains(t, err.Get("order_submission.price"), errors.New("is required when the order is of type LIMIT"))
}

func testPeggedOrderSubmissionWithLimitTypeAndNoPriceSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type:        typespb.Order_TYPE_LIMIT,
		PeggedOrder: &typespb.PeggedOrder{},
	})

	assert.NotContains(t, err.Get("order_submission.price"), errors.New("is required when the order is of type LIMIT"))
}

func testPeggedOrderSubmissionWithUnspecifiedReferenceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED,
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.reference"), commands.ErrIsRequired)
}

func testPeggedOrderSubmissionWithUndefinedReferenceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference(-42),
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.reference"), commands.ErrIsNotValid)
}

func testPeggedOrderSubmissionWithoutLimitTypeFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.Order_Type
	}{
		{
			msg:   "with MARKET",
			value: typespb.Order_TYPE_MARKET,
		}, {
			msg:   "with NETWORK",
			value: typespb.Order_TYPE_NETWORK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Type:        tc.value,
				PeggedOrder: &typespb.PeggedOrder{},
			})

			assert.Contains(t, err.Get("order_submission.type"), errors.New("is expected to be an order of type LIMIT when the order is pegged"))
		})
	}
}

func testPeggedOrderSubmissionWithLimitTypeSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Type:        typespb.Order_TYPE_LIMIT,
		PeggedOrder: &typespb.PeggedOrder{},
	})

	assert.NotContains(t, err.Get("order_submission.type"), errors.New("is expected to be an order of type LIMIT when the order is pegged"))
}

func testPeggedOrderSubmissionWithWrongTimeInForceFails(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.Order_TimeInForce
	}{
		{
			msg:   "with IOC",
			value: typespb.Order_TIME_IN_FORCE_IOC,
		}, {
			msg:   "with FOK",
			value: typespb.Order_TIME_IN_FORCE_FOK,
		}, {
			msg:   "with GFA",
			value: typespb.Order_TIME_IN_FORCE_GFA,
		}, {
			msg:   "with GFN",
			value: typespb.Order_TIME_IN_FORCE_GFN,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				TimeInForce: tc.value,
				PeggedOrder: &typespb.PeggedOrder{},
			})

			assert.Contains(t, err.Get("order_submission.time_in_force"), errors.New("is expected to have a time in force of type GTT or GTC when the order is pegged"))
		})
	}
}

func testPeggedOrderSubmissionWithRightTimeInForceSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value typespb.Order_TimeInForce
	}{
		{
			msg:   "with GTC",
			value: typespb.Order_TIME_IN_FORCE_GTC,
		}, {
			msg:   "with GTT",
			value: typespb.Order_TIME_IN_FORCE_GTT,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				TimeInForce: tc.value,
				PeggedOrder: &typespb.PeggedOrder{},
			})

			assert.NotContains(t, err.Get("order_submission.time_in_force"), errors.New("is expected to have a time in force of type GTT or GTC when the order is pegged"))
		})
	}
}

func testPeggedOrderSubmissionWithSideBuyAndBestAskReferenceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_BUY,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.reference"), errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
}

func testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_BUY,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
		},
	})

	assert.NotContains(t, err.Get("order_submission.pegged_order.reference"), errors.New("cannot have a reference of type BEST_ASK when on BUY side"))
}

func testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceAndPositiveOffsetFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_BUY,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
			Offset:    RandomPositiveI64(),
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be negative or zero"))
}

func testPeggedOrderSubmissionWithSideBuyAndBestBidReferenceAndNonPositiveOffsetSucceeds(t *testing.T) {
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
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Side: typespb.Side_SIDE_BUY,
				PeggedOrder: &typespb.PeggedOrder{
					Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
					Offset:    tc.value,
				},
			})

			assert.NotContains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be negative or zero"))
		})
	}
}

func testPeggedOrderSubmissionWithSideBuyAndMidReferenceAndNonNegativeOffsetFails(t *testing.T) {
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
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Side: typespb.Side_SIDE_BUY,
				PeggedOrder: &typespb.PeggedOrder{
					Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
					Offset:    tc.value,
				},
			})

			assert.Contains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be negative"))
		})
	}
}

func testPeggedOrderSubmissionWithSideBuyAndMidReferenceAndNegativeOffsetSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_BUY,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
			Offset:    RandomNegativeI64(),
		},
	})

	assert.NotContains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be negative"))
}

func testPeggedOrderSubmissionWithSideSellAndBestBidReferenceFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_SELL,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID,
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.reference"), errors.New("cannot have a reference of type BEST_BID when on SELL side"))
}

func testPeggedOrderSubmissionWithSideSellAndBestAskReferenceSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_SELL,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
		},
	})

	assert.NotContains(t, err.Get("order_submission.pegged_order.reference"), errors.New("cannot have a reference of type BEST_BID when on SELL side"))
}

func testPeggedOrderSubmissionWithSideSellAndBestAskReferenceAndNegativeOffsetFails(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_SELL,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
			Offset:    RandomNegativeI64(),
		},
	})

	assert.Contains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be positive or zero"))
}

func testPeggedOrderSubmissionWithSideSellAndBestAskReferenceAndNonNegativeOffsetSucceeds(t *testing.T) {
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
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Side: typespb.Side_SIDE_SELL,
				PeggedOrder: &typespb.PeggedOrder{
					Reference: typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK,
					Offset:    tc.value,
				},
			})

			assert.NotContains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be positive or zero"))
		})
	}
}

func testPeggedOrderSubmissionWithSideSellAndMidReferenceAndNonPositiveOffsetFails(t *testing.T) {
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
			err := checkOrderSubmission(&commandspb.OrderSubmission{
				Side: typespb.Side_SIDE_SELL,
				PeggedOrder: &typespb.PeggedOrder{
					Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
					Offset:    tc.value,
				},
			})

			assert.Contains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be positive"))
		})
	}
}

func testPeggedOrderSubmissionWithSideSellAndMidReferenceAndPositiveOffsetSucceeds(t *testing.T) {
	err := checkOrderSubmission(&commandspb.OrderSubmission{
		Side: typespb.Side_SIDE_SELL,
		PeggedOrder: &typespb.PeggedOrder{
			Reference: typespb.PeggedReference_PEGGED_REFERENCE_MID,
			Offset:    RandomPositiveI64(),
		},
	})

	assert.NotContains(t, err.Get("order_submission.pegged_order.offset"), errors.New("must be positive"))
}

func checkOrderSubmission(cmd *commandspb.OrderSubmission) commands.Errors {
	err := commands.CheckOrderSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
