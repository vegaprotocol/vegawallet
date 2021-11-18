package tests_test

import (
	"os"
	"path/filepath"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameWallet(t *testing.T) {
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
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithWalletCreation().
		WithName(walletName).
		WithVersion(2).
		LocatedUnder(home)

	// given
	newWalletName := vgrand.RandomStr(5)
	currentDir := filepath.Dir(generateKeyResp.Wallet.FilePath)
	newPath := filepath.Join(currentDir, newWalletName)

	// when
	err = os.Rename(generateKeyResp.Wallet.FilePath, newPath)

	// then
	require.NoError(t, err)

	// when
	listKeysResp, err := KeyList(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", newWalletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listKeysResp)
	require.Len(t, listKeysResp.Keys, 1)
	assert.Equal(t, listKeysResp.Keys[0].PublicKey, generateKeyResp.Key.PublicKey)
}
