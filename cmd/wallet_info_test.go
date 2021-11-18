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

func TestGetWalletInfoFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testGetWalletInfoFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testGetWalletInfoFlagsMissingWalletFails)
}

func testGetWalletInfoFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.GetWalletInfoFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.GetWalletInfoRequest{
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

func testGetWalletInfoFlagsMissingWalletFails(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := newGetWalletInfoFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func newGetWalletInfoFlags(t *testing.T, testDir string) *cmd.GetWalletInfoFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.GetWalletInfoFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}
}
