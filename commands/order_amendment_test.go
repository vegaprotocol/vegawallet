package commands_test

import (
	"testing"
	"time"

	"code.vegaprotocol.io/go-wallet/commands"
	"code.vegaprotocol.io/go-wallet/internal/proto"
	typespb "code.vegaprotocol.io/go-wallet/internal/proto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckOrderAmendment(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilOrderAmendmentFails)
	t.Run("amend order price - success", testAmendOrderJustPriceSuccess)
	t.Run("amend order reduce - success", testAmendOrderJustReduceSuccess)
	t.Run("amend order increase - success", testAmendOrderJustIncreaseSuccess)
	t.Run("amend order expiry - success", testAmendOrderJustExpirySuccess)
	t.Run("amend order tif - success", testAmendOrderJustTIFSuccess)
	t.Run("amend order expiry before creation time - success", testAmendOrderPastExpiry)
	t.Run("amend order empty - fail", testAmendOrderEmptyFail)
	t.Run("amend order empty - fail", testAmendEmptyFail)
	t.Run("amend order invalid expiry type - fail", testAmendOrderInvalidExpiryFail)
	t.Run("amend order tif to GFA - fail", testAmendOrderToGFA)
	t.Run("amend order tif to GFN - fail", testAmendOrderToGFN)
}

func testNilOrderAmendmentFails(t *testing.T) {
	err := checkOrderAmendment(nil)

	assert.Contains(t, err.Get("order_amendment"), commands.ErrIsRequired)
}

func testAmendOrderJustPriceSuccess(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:  "orderid",
		MarketId: "marketid",
		Price:    &typespb.Price{Value: 1000},
	}
	err := checkOrderAmendment(arg)

	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderJustReduceSuccess(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:   "orderid",
		MarketId:  "marketid",
		SizeDelta: -10,
	}
	err := checkOrderAmendment(arg)
	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderJustIncreaseSuccess(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:   "orderid",
		MarketId:  "marketid",
		SizeDelta: 10,
	}
	err := checkOrderAmendment(arg)
	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderJustExpirySuccess(t *testing.T) {
	now := time.Now()
	expires := now.Add(-2 * time.Hour)
	arg := &commandspb.OrderAmendment{
		OrderId:   "orderid",
		MarketId:  "marketid",
		ExpiresAt: &proto.Timestamp{Value: expires.UnixNano()},
	}
	err := checkOrderAmendment(arg)
	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderJustTIFSuccess(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:     "orderid",
		MarketId:    "marketid",
		TimeInForce: proto.Order_TIME_IN_FORCE_GTC,
	}
	err := checkOrderAmendment(arg)
	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderEmptyFail(t *testing.T) {
	arg := &commandspb.OrderAmendment{}
	err := checkOrderAmendment(arg)
	assert.Error(t, err)

	arg2 := &commandspb.OrderAmendment{
		OrderId:  "orderid",
		MarketId: "marketid",
	}
	err = checkOrderAmendment(arg2)
	assert.Error(t, err)
}

func testAmendEmptyFail(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:  "orderid",
		MarketId: "marketid",
	}
	err := checkOrderAmendment(arg)
	assert.Error(t, err)
}

func testAmendOrderInvalidExpiryFail(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:     "orderid",
		TimeInForce: proto.Order_TIME_IN_FORCE_GTC,
		ExpiresAt:   &proto.Timestamp{Value: 10},
	}
	err := checkOrderAmendment(arg)
	assert.Error(t, err)

	arg.TimeInForce = proto.Order_TIME_IN_FORCE_FOK
	err = checkOrderAmendment(arg)
	assert.Error(t, err)

	arg.TimeInForce = proto.Order_TIME_IN_FORCE_IOC
	err = checkOrderAmendment(arg)
	assert.Error(t, err)
}

/*
 * Sending an old expiry date is OK and should not be rejected here.
 * The validation should take place inside the core
 */
func testAmendOrderPastExpiry(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:     "orderid",
		MarketId:    "marketid",
		TimeInForce: proto.Order_TIME_IN_FORCE_GTT,
		ExpiresAt:   &proto.Timestamp{Value: 10},
	}
	err := checkOrderAmendment(arg)
	assert.NoError(t, err.ErrorOrNil())
}

func testAmendOrderToGFN(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:     "orderid",
		TimeInForce: proto.Order_TIME_IN_FORCE_GFN,
		ExpiresAt:   &proto.Timestamp{Value: 10},
	}
	err := checkOrderAmendment(arg)
	assert.Error(t, err)
}

func testAmendOrderToGFA(t *testing.T) {
	arg := &commandspb.OrderAmendment{
		OrderId:     "orderid",
		TimeInForce: proto.Order_TIME_IN_FORCE_GFA,
		ExpiresAt:   &proto.Timestamp{Value: 10},
	}
	err := checkOrderAmendment(arg)
	assert.Error(t, err)
}

func checkOrderAmendment(cmd *commandspb.OrderAmendment) commands.Errors {
	err := commands.CheckOrderAmendment(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
