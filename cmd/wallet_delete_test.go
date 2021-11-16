package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testDeleteWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testDeleteWalletFlagsMissingWalletFails)
}

func testDeleteWalletFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	walletName := vgrand.RandomStr(10)

	f := &cmd.DeleteWalletFlags{
		Wallet: walletName,
	}

	// when
	err := f.Validate()

	// then
	require.NoError(t, err)
}

func testDeleteWalletFlagsMissingWalletFails(t *testing.T) {
	// given
	f := newDeleteWalletFlags(t)
	f.Wallet = ""

	// when
	err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
}

func newDeleteWalletFlags(t *testing.T) *cmd.DeleteWalletFlags {
	t.Helper()

	walletName := vgrand.RandomStr(10)

	return &cmd.DeleteWalletFlags{
		Wallet: walletName,
	}
}
