package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

/**********************************************************************************/
/*                                   DELEGATION                                   */
/**********************************************************************************/
func TestSubmittingNoDelegateCommandFails(t *testing.T) {
	err := checkDelegateSubmission(nil)

	assert.Contains(t, err.Get("delegate_submission"), commands.ErrIsRequired)
}

func TestSubmittingNoDelegateNodeIdFails(t *testing.T) {
	cmd := &commandspb.DelegateSubmission{
		Amount: 1000,
	}
	err := checkDelegateSubmission(cmd)

	assert.Contains(t, err.Get("delegate_submission.node_id"), commands.ErrIsRequired)
}

func TestSubmittingNoDelegateAmountFails(t *testing.T) {
	cmd := &commandspb.DelegateSubmission{
		NodeId: "TestingNodeID",
	}
	err := checkDelegateSubmission(cmd)

	assert.Contains(t, err.Get("delegate_submission.amount"), commands.ErrIsRequired)
}

func checkDelegateSubmission(cmd *commandspb.DelegateSubmission) commands.Errors {
	err := commands.CheckDelegateSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}
	return e
}

/**********************************************************************************/
/*                                  UNDELEGATION                                  */
/**********************************************************************************/
func TestSubmittingNoUndelegateCommandFails(t *testing.T) {
	err := checkUndelegateAtEpochEndSubmission(nil)

	assert.Contains(t, err.Get("undelegateAtEpochEnd_submission"), commands.ErrIsRequired)
}

func TestSubmittingNoUndelegateNodeIdFails(t *testing.T) {
	cmd := &commandspb.UndelegateAtEpochEndSubmission{
		Amount: 1000,
	}
	err := checkUndelegateAtEpochEndSubmission(cmd)

	assert.Contains(t, err.Get("undelegateAtEpochEnd_submission.node_id"), commands.ErrIsRequired)
}

func TestSubmittingNoUndelegateAtEpochEndAmountFails(t *testing.T) {
	cmd := &commandspb.UndelegateAtEpochEndSubmission{
		NodeId: "TestingNodeID",
	}
	err := checkUndelegateAtEpochEndSubmission(cmd)

	assert.Contains(t, err.Get("undelegateAtEpochEnd_submission.amount"), commands.ErrIsRequired)
}

func checkUndelegateAtEpochEndSubmission(cmd *commandspb.UndelegateAtEpochEndSubmission) commands.Errors {
	err := commands.CheckUndelegateAtEpochEndSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}
	return e
}
