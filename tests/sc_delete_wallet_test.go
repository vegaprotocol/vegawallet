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
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	walletName := vgrand.RandomStr(5)

	// when
	createWalletResp, err := WalletCreate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertCreateWallet(t, createWalletResp).
		WithName(walletName).
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
	assert.NoFileExists(t, createWalletResp.Wallet.FilePath)
}

func TestDeleteNonExistingWallet(t *testing.T) {
	home := t.TempDir()

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
