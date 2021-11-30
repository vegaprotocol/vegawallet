package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteNetworkFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testDeleteNetworkFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testDeleteNetworkFlagsMissingNetworkFails)
}

func testDeleteNetworkFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	walletName := vgrand.RandomStr(10)

	f := &cmd.DeleteNetworkFlags{
		Network: walletName,
		Force:   true,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, f.Network, req.Name)
}

func testDeleteNetworkFlagsMissingNetworkFails(t *testing.T) {
	// given
	f := newDeleteNetworkFlags(t)
	f.Network = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("network"))
	assert.Nil(t, req)
}

func newDeleteNetworkFlags(t *testing.T) *cmd.DeleteNetworkFlags {
	t.Helper()

	walletName := vgrand.RandomStr(10)

	return &cmd.DeleteNetworkFlags{
		Network: walletName,
	}
}
