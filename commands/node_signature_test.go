package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckNodeSignature(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilNodeSignatureFails)
	t.Run("Submitting a node signature without id fails", testNodeSignatureWithoutIDFails)
	t.Run("Submitting a node signature with id succeeds", testNodeSignatureWithIDSucceeds)
	t.Run("Submitting a node signature without sig fails", testNodeSignatureWithoutSigFails)
	t.Run("Submitting a node signature with sig succeeds", testNodeSignatureWithSigSucceeds)
	t.Run("Submitting a node signature without kind fails", testNodeSignatureWithoutKindFails)
	t.Run("Submitting a node signature with invalid kind fails", testNodeSignatureWithInvalidKindFails)
	t.Run("Submitting a node signature with kind succeeds", testNodeSignatureWithKindSucceeds)
}

func testNilNodeSignatureFails(t *testing.T) {
	err := checkNodeSignature(nil)

	assert.Error(t, err)
}

func testNodeSignatureWithoutIDFails(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{})
	assert.Contains(t, err.Get("node_signature.id"), commands.ErrIsRequired)
}

func testNodeSignatureWithIDSucceeds(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{
		Id: "My ID",
	})
	assert.NotContains(t, err.Get("node_signature.id"), commands.ErrIsRequired)
}

func testNodeSignatureWithoutSigFails(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{})
	assert.Contains(t, err.Get("node_signature.sig"), commands.ErrIsRequired)
}

func testNodeSignatureWithSigSucceeds(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{
		Sig: []byte("0xDEADBEEF"),
	})
	assert.NotContains(t, err.Get("node_signature.sig"), commands.ErrIsRequired)
}

func testNodeSignatureWithoutKindFails(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{})
	assert.Contains(t, err.Get("node_signature.kind"), commands.ErrIsRequired)
}

func testNodeSignatureWithInvalidKindFails(t *testing.T) {
	err := checkNodeSignature(&commandspb.NodeSignature{
		Kind: commandspb.NodeSignatureKind(-42),
	})
	assert.Contains(t, err.Get("node_signature.kind"), commands.ErrIsNotValid)
}

func testNodeSignatureWithKindSucceeds(t *testing.T) {
	testCases := []struct {
		msg   string
		value commandspb.NodeSignatureKind
	}{
		{
			msg:   "with new kind",
			value: commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_NEW,
		}, {
			msg:   "with withdrawal kind",
			value: commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_ASSET_WITHDRAWAL,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			err := checkNodeSignature(&commandspb.NodeSignature{
				Kind: tc.value,
			})
			assert.NotContains(t, err.Get("node_signature.kind"), commands.ErrIsRequired)
			assert.NotContains(t, err.Get("node_signature.kind"), commands.ErrIsNotValid)
		})
	}
}

func checkNodeSignature(cmd *commandspb.NodeSignature) commands.Errors {
	err := commands.CheckNodeSignature(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
