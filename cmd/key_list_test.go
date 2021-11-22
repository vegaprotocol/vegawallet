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

func TestListKeysFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testListKeysFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testListKeysFlagsMissingWalletFails)
}

func testListKeysFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.ListKeysRequest{
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

func testListKeysFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newListKeysFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func newListKeysFlags(t *testing.T, testDir string) *cmd.ListKeysFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	return &cmd.ListKeysFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
	}
}
