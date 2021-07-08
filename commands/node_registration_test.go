package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
	"github.com/stretchr/testify/assert"
)

func TestCheckNodeRegistration(t *testing.T) {
	t.Run("Submitting a nil command fails", testNilNodeRegistrationFails)
	t.Run("Submitting a node registration without pub key fails", testNodeRegistrationWithoutPubKeyFails)
	t.Run("Submitting a node registration with pub key succeeds", testNodeRegistrationWithPubKeySucceeds)
	t.Run("Submitting a node registration without chain pub key fails", testNodeRegistrationWithoutChainPubKeyFails)
	t.Run("Submitting a node registration with chain pub key succeeds", testNodeRegistrationWithChainPubKeySucceeds)
}

func testNilNodeRegistrationFails(t *testing.T) {
	err := checkNodeRegistration(nil)

	assert.Error(t, err)
}

func testNodeRegistrationWithoutPubKeyFails(t *testing.T) {
	err := checkNodeRegistration(&commandspb.NodeRegistration{})
	assert.Contains(t, err.Get("node_registration.pub_key"), commands.ErrIsRequired)
}

func testNodeRegistrationWithPubKeySucceeds(t *testing.T) {
	err := checkNodeRegistration(&commandspb.NodeRegistration{
		PubKey: []byte("0xDEADBEEF"),
	})
	assert.NotContains(t, err.Get("node_registration.pub_key"), commands.ErrIsRequired)
}

func testNodeRegistrationWithoutChainPubKeyFails(t *testing.T) {
	err := checkNodeRegistration(&commandspb.NodeRegistration{})
	assert.Contains(t, err.Get("node_registration.chain_pub_key"), commands.ErrIsRequired)
}

func testNodeRegistrationWithChainPubKeySucceeds(t *testing.T) {
	err := checkNodeRegistration(&commandspb.NodeRegistration{
		ChainPubKey: []byte("0xDEADBEEF"),
	})
	assert.NotContains(t, err.Get("node_registration.chain_pub_key"), commands.ErrIsRequired)
}

func checkNodeRegistration(cmd *commandspb.NodeRegistration) commands.Errors {
	err := commands.CheckNodeRegistration(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}
