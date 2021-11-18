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

const mnemonic = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"

func TestImportWalletFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testImportWalletFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testImportWalletFlagsMissingWalletFails)
	t.Run("Missing mnemonic file fails", testImportWalletFlagsMissingMnemonicFileFails)
}

func testImportWalletFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	mnemonicFilePath := NewFile(t, testDir, "mnemonic.txt", mnemonic)
	walletName := vgrand.RandomStr(10)

	f := &cmd.ImportWalletFlags{
		Wallet:         walletName,
		MnemonicFile:   mnemonicFilePath,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.ImportWalletRequest{
		Wallet:     walletName,
		Mnemonic:   mnemonic,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testImportWalletFlagsMissingWalletFails(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := newImportWalletFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testImportWalletFlagsMissingMnemonicFileFails(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := newImportWalletFlags(t, testDir)
	f.MnemonicFile = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("mnemonic-file"))
	assert.Nil(t, req)
}

func newImportWalletFlags(t *testing.T, testDir string) *cmd.ImportWalletFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	NewFile(t, testDir, "mnemonic.txt", mnemonic)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.ImportWalletFlags{
		Wallet:         walletName,
		MnemonicFile:   pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
