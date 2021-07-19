package commands

import (
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
)

func CheckDelegateSubmission(cmd *commandspb.DelegateSubmission) error {
	return checkDelegateSubmission(cmd).ErrorOrNil()
}

func checkDelegateSubmission(cmd *commandspb.DelegateSubmission) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("delegate_submission", ErrIsRequired)
	}

	if cmd.Amount <= 0 {
		errs.AddForProperty("delegate_submission.amount", ErrIsRequired)
	}

	if len(cmd.NodeId) <= 0 {
		errs.AddForProperty("delegate_submission.node_id", ErrIsRequired)
	}

	return errs
}

func CheckUndelegateAtEpochEndSubmission(cmd *commandspb.UndelegateAtEpochEndSubmission) error {
	return checkUndelegateAtEpochEndSubmission(cmd).ErrorOrNil()
}

func checkUndelegateAtEpochEndSubmission(cmd *commandspb.UndelegateAtEpochEndSubmission) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("undelegateAtEpochEnd_submission", ErrIsRequired)
	}

	if cmd.Amount <= 0 {
		errs.AddForProperty("undelegateAtEpochEnd_submission.amount", ErrIsRequired)
	}

	if len(cmd.NodeId) <= 0 {
		errs.AddForProperty("undelegateAtEpochEnd_submission.node_id", ErrIsRequired)
	}

	return errs
}
