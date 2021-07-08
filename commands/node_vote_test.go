package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckNodeVote(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilNodeVoteFails)
	t.Run("Submitting a node vote without pub key fails", testNodeVoteWithoutPubKeyFails)
	t.Run("Submitting a node vote with pub key succeeds", testNodeVoteWithPubKeySucceeds)
	t.Run("Submitting a node vote without reference fails", testNodeVoteWithoutReferenceFails)
	t.Run("Submitting a node vote with reference succeeds", testNodeVoteWithReferenceSucceeds)
}

func testNilNodeVoteFails(t *testing.T) {
	err := checkNodeVote(nil)

	assert.Error(t, err)
}

func testNodeVoteWithoutPubKeyFails(t *testing.T) {
	err := checkNodeVote(&commandspb.NodeVote{})
	assert.Contains(t, err.Get("node_vote.pub_key"), commands.ErrIsRequired)
}

func testNodeVoteWithPubKeySucceeds(t *testing.T) {
	err := checkNodeVote(&commandspb.NodeVote{
		PubKey: []byte("0xDEADBEEF"),
	})
	assert.NotContains(t, err.Get("node_vote.pub_key"), commands.ErrIsRequired)
}

func testNodeVoteWithoutReferenceFails(t *testing.T) {
	err := checkNodeVote(&commandspb.NodeVote{})
	assert.Contains(t, err.Get("node_vote.reference"), commands.ErrIsRequired)
}

func testNodeVoteWithReferenceSucceeds(t *testing.T) {
	err := checkNodeVote(&commandspb.NodeVote{
		Reference: "my ref",
	})
	assert.NotContains(t, err.Get("node_vote.reference"), commands.ErrIsRequired)
}

func checkNodeVote(cmd *commandspb.NodeVote) commands.Errors {
	err := commands.CheckNodeVote(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
