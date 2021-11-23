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
	t.Run("Unsupported log level fails", testRunServiceFlagsUnsupportedLogLevelFails)
	t.Run("No browser without console proxy fails", testRunServiceFlagsNoBrowserWithoutConsoleProxyFails)
}

func testRunServiceFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	networkName := vgrand.RandomStr(10)

	f := &cmd.RunServiceFlags{
		Network:      networkName,
		StartConsole: true,
		NoBrowser:    true,
		LogLevel:     "debug",
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

func testRunServiceFlagsUnsupportedLogLevelFails(t *testing.T) {
	// given
	f := newRunServiceFlags(t)
	f.LogLevel = vgrand.RandomStr(2)

	// when
	err := f.Validate()

	// then
	assert.ErrorIs(t, err, cmd.NewUnsupportedFlagValueError(f.LogLevel))
}

func testRunServiceFlagsNoBrowserWithoutConsoleProxyFails(t *testing.T) {
	// given
	f := newRunServiceFlags(t)
	f.StartConsole = false
	f.NoBrowser = true

	// when
	err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.ParentFlagMustBeSpecifiedError("no-browser", "console-proxy"))
}

func newRunServiceFlags(t *testing.T) *cmd.RunServiceFlags {
	t.Helper()

	networkName := vgrand.RandomStr(10)

	return &cmd.RunServiceFlags{
		Network:      networkName,
		StartConsole: true,
		NoBrowser:    true,
		LogLevel:     "debug",
	}
}
