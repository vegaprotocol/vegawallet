package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestSubmittingNilVoteFails(t *testing.T) {
	err := checkVoteSubmission(nil)

	assert.Contains(t, err.Get("vote_submission"), commands.ErrIsRequired)
}

func TestVoteSubmission(t *testing.T) {
	var cases = []struct {
		vote      commandspb.VoteSubmission
		errString string
	}{
		{
			vote: commandspb.VoteSubmission{
				Value:      typespb.Vote_VALUE_YES,
				ProposalId: "OKPROPOSALID",
			},
		},
		{
			vote: commandspb.VoteSubmission{
				ProposalId: "OKPROPOSALID",
			},
			errString: "vote_submission.value (is required)",
		},
		{
			vote: commandspb.VoteSubmission{
				Value:      typespb.Vote_Value(-42),
				ProposalId: "OKPROPOSALID",
			},
			errString: "vote_submission.value (is not a valid value)",
		},
		{
			vote: commandspb.VoteSubmission{
				Value: typespb.Vote_VALUE_NO,
			},
			errString: "vote_submission.proposal_id (is required)",
		},
		{
			vote:      commandspb.VoteSubmission{},
			errString: "vote_submission.proposal_id (is required), vote_submission.value (is required)",
		},
	}

	for _, c := range cases {
		err := commands.CheckVoteSubmission(&c.vote)
		if len(c.errString) <= 0 {
			assert.NoError(t, err)
			continue
		}
		assert.Error(t, err)
		assert.EqualError(t, err, c.errString)
	}
}

func checkVoteSubmission(cmd *commandspb.VoteSubmission) commands.Errors {
	err := commands.CheckVoteSubmission(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
