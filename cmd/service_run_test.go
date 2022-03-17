package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunServiceFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testRunServiceFlagsValidFlagsSucceeds)
	t.Run("Missing network fails", testRunServiceFlagsMissingNetworkFails)
	t.Run("No browser without console nor token dApp fails", testRunServiceFlagsNoBrowserWithoutConsoleNorTokenDAppFails)
}

func testRunServiceFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	networkName := vgrand.RandomStr(10)

	f := &cmd.RunServiceFlags{
		Network:       networkName,
		WithConsole:   true,
		WithTokenDApp: true,
		NoBrowser:     true,
	}

	// when
	err := f.Validate()

	// then
	require.NoError(t, err)
}

func testRunServiceFlagsMissingNetworkFails(t *testing.T) {
	// given
	f := newRunServiceFlags(t)
	f.Network = ""

	// when
	err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("network"))
}

func testRunServiceFlagsNoBrowserWithoutConsoleNorTokenDAppFails(t *testing.T) {
	// given
	f := newRunServiceFlags(t)
	f.WithConsole = false
	f.WithTokenDApp = false
	f.NoBrowser = true

	// when
	err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.OneOfParentsFlagMustBeSpecifiedError("no-browser", "with-console", "with-token-dapp"))
}

func newRunServiceFlags(t *testing.T) *cmd.RunServiceFlags {
	t.Helper()

	networkName := vgrand.RandomStr(10)

	return &cmd.RunServiceFlags{
		Network:       networkName,
		WithConsole:   true,
		WithTokenDApp: true,
		NoBrowser:     true,
	}
}
