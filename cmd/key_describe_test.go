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

func TestDescribeKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testKeyDescribeValidFlagsSucceeds)
	t.Run("Missing wallet fails", testKeyMissingWalletFails)
	t.Run("Missing public key fails", testKeyMissingPublicKeyFails)
}

func testKeyDescribeValidFlagsSucceeds(t *testing.T) {
	// given
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(10)

	f := &cmd.DescribeKeyFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
		PubKey:         pubKey,
	}

	expectedReq := &wallet.DescribeKeyRequest{
		Wallet:     walletName,
		Passphrase: passphrase,
		PubKey:     pubKey,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testKeyMissingWalletFails(t *testing.T) {
	// given
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	pubKey := vgrand.RandomStr(10)

	f := &cmd.DescribeKeyFlags{
		PassphraseFile: passphraseFilePath,
		PubKey:         pubKey,
	}

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	require.Nil(t, req)
}

func testKeyMissingPublicKeyFails(t *testing.T) {
	// given
	testDir, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	// given
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	walletName := vgrand.RandomStr(10)

	f := &cmd.DescribeKeyFlags{
		PassphraseFile: passphraseFilePath,
		Wallet:         walletName,
	}

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	require.Nil(t, req)
}
