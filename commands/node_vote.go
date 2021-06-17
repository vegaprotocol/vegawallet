package commands

import commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

func CheckNodeVote(cmd *commandspb.NodeVote) error {
	return checkNodeVote(cmd).ErrorOrNil()
}

func checkNodeVote(cmd *commandspb.NodeVote) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("node_vote", ErrIsRequired)
	}

	if len(cmd.PubKey) == 0 {
		errs.AddForProperty("node_vote.pub_key", ErrIsRequired)
	}

	if len(cmd.Reference) == 0 {
		errs.AddForProperty("node_vote.reference", ErrIsRequired)
	}

	return errs
}
