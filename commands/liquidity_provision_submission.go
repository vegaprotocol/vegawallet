package commands

import (
	"errors"
	"fmt"
	"strconv"

	typespb "code.vegaprotocol.io/go-wallet/internal/proto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
)

var (
	ErrOrderInShapeWithoutReference         = errors.New("order in shape without reference")
	ErrOrderInShapeWithoutProportion        = errors.New("order in shape without a proportion")
	ErrOrderInBuySideShapeWithBestAskPrice  = errors.New("order in buy side shape with best ask price reference")
	ErrOrderInBuySideShapeOffsetSup0        = errors.New("order in buy side shape offset must be <= 0")
	ErrOrderInBuySideShapeOffsetSupEq0      = errors.New("order in buy side shape offset must be < 0")
	ErrOrderInSellSideShapeOffsetInf0       = errors.New("order in sell shape offset must be >= 0")
	ErrOrderInSellSideShapeWithBestBidPrice = errors.New("order in sell side shape with best bid price reference")
	ErrOrderInSellSideShapeOffsetInfEq0     = errors.New("order in sell shape offset must be > 0")
)

func CheckLiquidityProvisionSubmission(cmd *commandspb.LiquidityProvisionSubmission) error {
	return checkLiquidityProvisionSubmission(cmd).ErrorOrNil()
}

func checkLiquidityProvisionSubmission(cmd *commandspb.LiquidityProvisionSubmission) Errors {
	var errs = NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("liquidity_provision_submission", ErrIsRequired)
	}

	if len(cmd.MarketId) <= 0 {
		errs.AddForProperty("liquidity_provision_submission.market_id", ErrIsRequired)
	}

	// if the commitment amount is 0, then the command should be interpreted as
	// a cancellation of the liquidity provision. As a result, the validation
	// shouldn't be made on the rest of the field.
	// However, since the user might by sending an blank command to probe the
	// validation, we want to return as many error message as possible.
	// A cancellation is only valid if a market is specified, and the commitment is
	// 0. In any case the core will consider that as a cancellation, so we return
	// the error that we go from the market id check.
	if cmd.CommitmentAmount == 0 {
		return errs
	}

	if len(cmd.Fee) <= 0 {
		errs.AddForProperty("liquidity_provision_submission.fee", ErrIsRequired)
	} else {
		if fee, err := strconv.ParseFloat(cmd.Fee, 64); err != nil {
			errs.AddForProperty(
				"liquidity_provision_submission.fee",
				ErrIsNotValid,
			)
		} else if fee < 0 {
			errs.AddForProperty("liquidity_provision_submission.fee", ErrMustBePositive)
		}

	}

	errs.Merge(checkLiquidityProvisionShape(cmd.Buys, typespb.Side_SIDE_BUY))
	errs.Merge(checkLiquidityProvisionShape(cmd.Sells, typespb.Side_SIDE_SELL))

	return errs
}

func checkLiquidityProvisionShape(
	orders []*typespb.LiquidityOrder, side typespb.Side,
) Errors {
	var (
		errs           = NewErrors()
		shapeSideField = "liquidity_provision_submission.buys"
	)
	if side == typespb.Side_SIDE_SELL {
		shapeSideField = "liquidity_provision_submission.sells"
	}

	if len(orders) <= 0 {
		errs.AddForProperty(shapeSideField, errors.New("empty shape"))
		return errs

	}

	for idx, order := range orders {
		if order.Reference == typespb.PeggedReference_PEGGED_REFERENCE_UNSPECIFIED {
			errs.AddForProperty(
				fmt.Sprintf("%v.%d.reference", shapeSideField, idx),
				ErrOrderInShapeWithoutReference,
			)
		}
		if order.Proportion == 0 {
			errs.AddForProperty(
				fmt.Sprintf("%v.%d.proportion", shapeSideField, idx),
				ErrOrderInShapeWithoutProportion,
			)
		}

		if side == typespb.Side_SIDE_BUY {
			switch order.Reference {
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK:
				errs.AddForProperty(
					fmt.Sprintf("%v.%d.reference", shapeSideField, idx),
					ErrOrderInBuySideShapeWithBestAskPrice,
				)
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID:
				if order.Offset > 0 {
					errs.AddForProperty(
						fmt.Sprintf("%v.%d.offset", shapeSideField, idx),
						ErrOrderInBuySideShapeOffsetSup0,
					)
				}
			case typespb.PeggedReference_PEGGED_REFERENCE_MID:
				if order.Offset >= 0 {
					errs.AddForProperty(
						fmt.Sprintf("%v.%d.offset", shapeSideField, idx),
						ErrOrderInBuySideShapeOffsetSupEq0,
					)
				}
			}
		} else {
			switch order.Reference {
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_ASK:
				if order.Offset < 0 {
					errs.AddForProperty(
						fmt.Sprintf("%v.%d.offset", shapeSideField, idx),
						ErrOrderInSellSideShapeOffsetInf0,
					)
				}
			case typespb.PeggedReference_PEGGED_REFERENCE_BEST_BID:
				errs.AddForProperty(
					fmt.Sprintf("%v.%d.offset", shapeSideField, idx),
					ErrOrderInSellSideShapeWithBestBidPrice,
				)
			case typespb.PeggedReference_PEGGED_REFERENCE_MID:
				if order.Offset <= 0 {
					errs.AddForProperty(
						fmt.Sprintf("%v.%d.offset", shapeSideField, idx),
						ErrOrderInSellSideShapeOffsetInfEq0,
					)
				}
			}
		}
	}
	return errs
}
