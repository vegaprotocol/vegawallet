package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	"github.com/stretchr/testify/assert"
)

func TestCheckOracleDataSubmission(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilOracleDataSubmissionFails)
	t.Run("Submitting an oracle data without payload fails", testOracleDataSubmissionWithoutPayloadFails)
	t.Run("Submitting an oracle data with payload succeeds", testOracleDataSubmissionWithPayloadSucceeds)
	t.Run("Submitting an oracle data without source fails", testOracleDataSubmissionWithoutSourceFails)
	t.Run("Submitting an oracle data with invalid source fails", testOracleDataSubmissionWithInvalidSourceFails)
	t.Run("Submitting an oracle data with source succeeds", testOracleDataSubmissionWithSourceSucceeds)
}

func testNilOracleDataSubmissionFails(t *testing.T) {
	err := checkOracleDataSubmission(nil)

	assert.Contains(t, err.Get("oracle_data_submission"), commands.ErrIsRequired)
}

func testOracleDataSubmissionWithoutPayloadFails(t *testing.T) {
	err := checkOracleDataSubmission(&commandspb.OracleDataSubmission{})
	assert.Contains(t, err.Get("oracle_data_submission.payload"), commands.ErrIsRequired)
}

func testOracleDataSubmissionWithPayloadSucceeds(t *testing.T) {
	err := checkOracleDataSubmission(&commandspb.OracleDataSubmission{
		Payload: []byte("0xDEADBEEF"),
	})
	assert.NotContains(t, err.Get("oracle_data_submission.payload"), commands.ErrIsRequired)
}


func testOracleDataSubmissionWithoutSourceFails(t *testing.T) {
	err := checkOracleDataSubmission(&commandspb.OracleDataSubmission{})
	assert.Contains(t, err.Get("oracle_data_submission.source"), commands.ErrIsRequired)
}

func testOracleDataSubmissionWithInvalidSourceFails(t *testing.T) {
	err := checkOracleDataSubmission(&commandspb.OracleDataSubmission{
		Source: commandspb.OracleDataSubmission_OracleSource(-42),
	})
	assert.Contains(t, err.Get("oracle_data_submission.source"), commands.ErrIsNotValid)
}

func testOracleDataSubmissionWithSourceSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value commandspb.OracleDataSubmission_OracleSource
	}{
		{
			msg:   "with Open Oracle source",
			value: commandspb.OracleDataSubmission_ORACLE_SOURCE_OPEN_ORACLE,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkOracleDataSubmission(&commandspb.OracleDataSubmission{
				Source: tc.value,
			})
			assert.NotContains(t, err.Get("oracle_data_submission.source"), commands.ErrIsRequired)
			assert.NotContains(t, err.Get("oracle_data_submission.source"), commands.ErrIsNotValid)
		})
	}
}

func checkOracleDataSubmission(cmd *commandspb.OracleDataSubmission) commands.Errors {
	err := commands.CheckOracleDataSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
