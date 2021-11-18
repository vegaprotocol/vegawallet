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

func TestIsolateKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testIsolateKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testIsolateKeyFlagsMissingWalletFails)
	t.Run("Missing public key fails", testIsolateKeyFlagsMissingPubKeyFails)
}

func testIsolateKeyFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	f := &cmd.IsolateKeyFlags{
		Wallet:         walletName,
		PubKey:         pubKey,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.IsolateKeyRequest{
		Wallet:     walletName,
		PubKey:     pubKey,
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testIsolateKeyFlagsMissingWalletFails(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := newIsolateKeyFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testIsolateKeyFlagsMissingPubKeyFails(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := newIsolateKeyFlags(t, testDir)
	f.PubKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	assert.Nil(t, req)
}

func newIsolateKeyFlags(t *testing.T, testDir string) *cmd.IsolateKeyFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.IsolateKeyFlags{
		Wallet:         walletName,
		PubKey:         pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
