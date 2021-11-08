package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestAnnotateKey(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

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
	err = KeyAnnotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:prefer-this-name",
	})

	// then
	require.NoError(t, err)

	// when
	err = KeyAnnotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--passphrase-file", passphraseFilePath,
		"--clear",
	})

	// then
	require.NoError(t, err)

	// when
	err = KeyAnnotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-1,role:validation",
	})

	// then
	require.NoError(t, err)
}
