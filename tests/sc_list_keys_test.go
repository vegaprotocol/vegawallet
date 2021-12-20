package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListKeys(t *testing.T) {
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
	listKeysResp1, err := KeyList(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listKeysResp1)
	require.Len(t, listKeysResp1.Keys, 1)
	assert.Equal(t, listKeysResp1.Keys[0].PublicKey, createWalletResp.Key.PublicKey)

	// when
	generateKeyResp, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithMeta(map[string]string{"name": DefaultMetaName(t, walletName, 2)})

	// when
	listKeysResp2, err := KeyList(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listKeysResp2)
	require.Len(t, listKeysResp2.Keys, 2)
	assert.Equal(t, listKeysResp2.Keys[0].PublicKey, createWalletResp.Key.PublicKey)
	assert.Equal(t, listKeysResp2.Keys[1].PublicKey, generateKeyResp.PublicKey)
}
