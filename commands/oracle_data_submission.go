package commands

import commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

func CheckOracleDataSubmission(cmd *commandspb.OracleDataSubmission) error {
	return checkOracleDataSubmission(cmd).ErrorOrNil()
}

func checkOracleDataSubmission(cmd *commandspb.OracleDataSubmission) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("oracle_data_submission", ErrIsRequired)
	}

	if len(cmd.Payload) == 0 {
		errs.AddForProperty("oracle_data_submission.payload", ErrIsRequired)
	}

	if cmd.Source == commandspb.OracleDataSubmission_ORACLE_SOURCE_UNSPECIFIED {
		errs.AddForProperty("oracle_data_submission.source", ErrIsRequired)
	}
	if _, ok := commandspb.OracleDataSubmission_OracleSource_name[int32(cmd.Source)]; !ok {
		errs.AddForProperty("oracle_data_submission.source", ErrIsNotValid)
	}

	return errs
}
