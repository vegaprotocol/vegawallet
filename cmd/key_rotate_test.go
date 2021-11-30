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

func TestRotateKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testRotateKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testRotateKeyFlagsMissingWalletFails)
	t.Run("Missing new public key fails", testRotateKeyFlagsMissingPublicKeyFails)
	t.Run("Missing current public key fails", testRotateKeyFlagsMissingCurrentPublicKeyFails)
	t.Run("Missing tx height fails", testRotateKeyFlagsMissingTxBlockHeightFails)
	t.Run("Missing target height fails", testRotateKeyFlagsMissingTargetBlockHeightFails)
	t.Run("Validate fails when target height is less then tx height", testRotateKeyFlagsTargetFailsWhenBlockHeightIsLessThanTXHeight)
}

func testRotateKeyFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	currentPubKey := vgrand.RandomStr(20)
	pubKey := vgrand.RandomStr(20)
	txBlockHeight := uint64(20)
	targetBlockHeight := uint64(25)

	f := &cmd.RotateKeyFlags{
		Wallet:            walletName,
		PassphraseFile:    passphraseFilePath,
		CurrentPubKey:     currentPubKey,
		PublicKey:         pubKey,
		TxBlockHeight:     txBlockHeight,
		TargetBlockHeight: targetBlockHeight,
	}

	expectedReq := &wallet.RotateKeyRequest{
		Wallet:            walletName,
		Passphrase:        passphrase,
		CurrentPublicKey:  currentPubKey,
		PublicKey:         pubKey,
		TxBlockHeight:     txBlockHeight,
		TargetBlockHeight: targetBlockHeight,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testRotateKeyFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testRotateKeyFlagsMissingTxBlockHeightFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.TxBlockHeight = 0

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("tx-height"))
	assert.Nil(t, req)
}

func testRotateKeyFlagsMissingTargetBlockHeightFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.TargetBlockHeight = 0

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("target-height"))
	assert.Nil(t, req)
}

func testRotateKeyFlagsMissingPublicKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.PublicKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	assert.Nil(t, req)
}

func testRotateKeyFlagsMissingCurrentPublicKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.CurrentPubKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("current-pubkey"))
	assert.Nil(t, req)
}

func testRotateKeyFlagsTargetFailsWhenBlockHeightIsLessThanTXHeight(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newRotateKeyFlags(t, testDir)
	f.TxBlockHeight = 25
	f.TargetBlockHeight = 20

	// when
	req, err := f.Validate()

	// then
	assert.EqualError(t, err, "--target-height flag must be greater than --tx-height")
	assert.Nil(t, req)
}

func newRotateKeyFlags(t *testing.T, testDir string) *cmd.RotateKeyFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)
	currentPubKey := vgrand.RandomStr(20)
	txBlockHeight := uint64(20)
	targetBlockHeight := uint64(25)

	return &cmd.RotateKeyFlags{
		Wallet:            walletName,
		PublicKey:         pubKey,
		CurrentPubKey:     currentPubKey,
		PassphraseFile:    passphraseFilePath,
		TxBlockHeight:     txBlockHeight,
		TargetBlockHeight: targetBlockHeight,
	}
}
