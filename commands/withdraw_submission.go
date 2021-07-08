package commands

import (
	"errors"

	typespb "code.vegaprotocol.io/go-wallet/internal/proto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
)

func CheckWithdrawSubmission(cmd *commandspb.WithdrawSubmission) error {
	return checkWithdrawSubmission(cmd).ErrorOrNil()
}

func checkWithdrawSubmission(cmd *commandspb.WithdrawSubmission) Errors {
	var errs = NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("withdraw_submission", ErrIsRequired)
	}

	if cmd.Amount <= 0 {
		errs.AddForProperty("withdraw_submission.amount", ErrIsRequired)
	}

	if len(cmd.Asset) <= 0 {
		errs.AddForProperty("withdraw_submission.asset", ErrIsRequired)
	}

	if cmd.Ext != nil {
		errs.Merge(checkWithdrawExt(cmd.Ext))
	}

	return errs
}

func checkWithdrawExt(wext *typespb.WithdrawExt) Errors {
	var errs = NewErrors()
	switch v := wext.Ext.(type) {
	case *typespb.WithdrawExt_Erc20:
		if len(v.Erc20.ReceiverAddress) <= 0 {
			errs.AddForProperty(
				"withdraw_ext.erc20.received_address",
				ErrIsRequired,
			)
		}
	default:
		errs.AddForProperty("withdraw_ext.ext", errors.New("unsupported withdraw extended details"))
	}
	return errs
}
