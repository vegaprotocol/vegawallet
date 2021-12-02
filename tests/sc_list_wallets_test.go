package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListWallets(t *testing.T) {
	// given
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	walletName1 := "a" + vgrand.RandomStr(5)

	// when
	createWalletResp1, err := WalletCreate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName1,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertCreateWallet(t, createWalletResp1).
		WithName(walletName1).
		LocatedUnder(home)

	// when
	listWalletsResp1, err := WalletList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listWalletsResp1)
	require.Len(t, listWalletsResp1.Wallets, 1)
	assert.Equal(t, listWalletsResp1.Wallets[0], createWalletResp1.Wallet.Name)

	// given
	walletName2 := "b" + vgrand.RandomStr(5)

	// when
	createWalletResp2, err := WalletCreate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName2,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertCreateWallet(t, createWalletResp2).
		WithName(walletName2).
		LocatedUnder(home)

	// when
	listWalletsResp2, err := WalletList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listWalletsResp2)
	require.Len(t, listWalletsResp2.Wallets, 2)
	assert.Equal(t, listWalletsResp2.Wallets[0], createWalletResp1.Wallet.Name)
	assert.Equal(t, listWalletsResp2.Wallets[1], createWalletResp2.Wallet.Name)
}
