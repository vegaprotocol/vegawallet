package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testCreateWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testCreateWalletFlagsMissingWalletFails)
}

func testCreateWalletFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	walletName := vgrand.RandomStr(10)
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	f := &cmd.CreateWalletFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.CreateWalletRequest{
		Wallet:     walletName,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testCreateWalletFlagsMissingWalletFails(t *testing.T) {
	// given
	f := newCreateWalletFlags(t)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func newCreateWalletFlags(t *testing.T) *cmd.CreateWalletFlags {
	t.Helper()
	return &cmd.CreateWalletFlags{
		Wallet:         vgrand.RandomStr(10),
		PassphraseFile: "/some/fake/path",
	}
}
