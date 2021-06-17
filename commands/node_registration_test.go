package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	"github.com/stretchr/testify/assert"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
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
