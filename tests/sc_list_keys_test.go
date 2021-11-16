package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListKeys(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	walletName := vgrand.RandomStr(5)

	// when
	generateKeyResp1, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-1,role:validation",
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp1).
		WithWalletCreation().
		WithName(walletName).
		WithVersion(2).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
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
	assert.Equal(t, listKeysResp1.Keys[0].PublicKey, generateKeyResp1.Key.KeyPair.PublicKey)

	// when
	generateKeyResp2, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp2).
		WithoutWalletCreation().
		WithName(walletName).
		WithMeta(map[string]string{"name": DefaultMetaName(t, walletName, 2)}).
		LocatedUnder(home)

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
	assert.Equal(t, listKeysResp2.Keys[0].PublicKey, generateKeyResp1.Key.KeyPair.PublicKey)
	assert.Equal(t, listKeysResp2.Keys[1].PublicKey, generateKeyResp2.Key.KeyPair.PublicKey)
}
