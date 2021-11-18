package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteWallet(t *testing.T) {
	// given
	home := NewTempDir(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	walletName := vgrand.RandomStr(5)

	// when
	generateKeyResp, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-1,role:validation",
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithWalletCreation().
		WithName(walletName).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	err = WalletDelete(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--force",
	})

	// then
	require.NoError(t, err)
	assert.NoFileExists(t, generateKeyResp.Wallet.FilePath)
}

func TestDeleteNonExistingWallet(t *testing.T) {
	home := NewTempDir(t)

	// when
	err := WalletDelete(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", vgrand.RandomStr(5),
		"--force",
	})

	// then
	require.ErrorIs(t, err, wallet.ErrWalletDoesNotExists)
}
