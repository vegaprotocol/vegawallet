package commands

import (
	"errors"

	typespb "code.vegaprotocol.io/go-wallet/internal/proto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
)

func CheckOrderAmendment(cmd *commandspb.OrderAmendment) error {
	return checkOrderAmendment(cmd).ErrorOrNil()
}

func checkOrderAmendment(cmd *commandspb.OrderAmendment) Errors {
	var (
		errs       = NewErrors()
		isAmending bool
	)

	if cmd == nil {
		return errs.FinalAddForProperty("order_amendment", ErrIsRequired)
	}

	if len(cmd.OrderId) <= 0 {
		errs.AddForProperty("order_amendment.order_id", ErrIsRequired)
	}

	if len(cmd.MarketId) <= 0 {
		errs.AddForProperty("order_amendment.market_id", ErrIsRequired)
	}

	// Check we are not trying to amend to a GFA
	if cmd.TimeInForce == typespb.Order_TIME_IN_FORCE_GFA {
		errs.AddForProperty("order_amendment.time_in_force", ErrCannotAmendToGFA)
	}

	// Check we are not trying to amend to a GFN
	if cmd.TimeInForce == typespb.Order_TIME_IN_FORCE_GFN {
		errs.AddForProperty("order_amendment.time_in_force", ErrCannotAmendToGFN)
	}

	if cmd.Price != nil {
		isAmending = true
		if cmd.Price.Value == 0 {
			errs.AddForProperty("order_amendment.price", ErrCannotAmendToGFN)
		}
	}

	if cmd.SizeDelta != 0 {
		isAmending = true
	}

	if cmd.TimeInForce == typespb.Order_TIME_IN_FORCE_GTT {
		isAmending = true
		if cmd.ExpiresAt == nil {
			errs.AddForProperty(
				"order_amendment.time_in_force", ErrGTTOrderWithNoExpiry)
		}
	}

	if cmd.TimeInForce != typespb.Order_TIME_IN_FORCE_UNSPECIFIED {
		isAmending = true
		if _, ok := typespb.Order_TimeInForce_name[int32(cmd.TimeInForce)]; !ok {
			errs.AddForProperty("order_amendment.time_in_force", ErrIsNotValid)
		}
	}

	if cmd.PeggedReference != typespb.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED {
		isAmending = true
		if _, ok := typespb.PeggedReference_name[int32(cmd.PeggedReference)]; !ok {
			errs.AddForProperty("order_amendment.pegged_reference", ErrIsNotValid)
		}
	}

	if cmd.ExpiresAt != nil && cmd.ExpiresAt.Value > 0 {
		isAmending = true
		if cmd.TimeInForce != typespb.Order_TIME_IN_FORCE_GTT &&
			cmd.TimeInForce != typespb.Order_TIME_IN_FORCE_UNSPECIFIED {
			errs.AddForProperty(
				"order_amendment.expires_at", ErrNonGTTOrderWithExpiry)
		}
	}

	if cmd.PeggedOffset != nil {
		isAmending = true
	}

	if !isAmending {
		errs.Add(errors.New("order_amendment does not amend anything"))
	}

	return errs
}
