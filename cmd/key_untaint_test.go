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

func TestUntaintKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testUntaintKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testUntaintKeyFlagsMissingWalletFails)
	t.Run("Missing public key fails", testUntaintKeyFlagsMissingPubKeyFails)
}

func testUntaintKeyFlagsValidFlagsSucceeds(t *testing.T) {
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	f := &cmd.UntaintKeyFlags{
		Wallet:         walletName,
		PubKey:         pubKey,
		PassphraseFile: passphraseFilePath,
	}

	expectedReq := &wallet.UntaintKeyRequest{
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

func testUntaintKeyFlagsMissingWalletFails(t *testing.T) {
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	f := newUntaintKeyFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testUntaintKeyFlagsMissingPubKeyFails(t *testing.T) {
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	f := newUntaintKeyFlags(t, testDir)
	f.PubKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	assert.Nil(t, req)
}

func newUntaintKeyFlags(t *testing.T, testDir string) *cmd.UntaintKeyFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.UntaintKeyFlags{
		Wallet:         walletName,
		PubKey:         pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
