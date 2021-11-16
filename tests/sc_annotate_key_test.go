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

	cmd := []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	}

	// when
	generateKeyResp, err := KeyGenerate(t, append(cmd,
		"--meta", "name:key-1,role:validation",
	))

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithWalletCreation().
		WithName(walletName).
		WithVersion(2).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	err = KeyAnnotate(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--meta", "name:prefer-this-name",
	))

	// then
	require.NoError(t, err)

	// when
	descResp, err := KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": "prefer-this-name"})

	// when
	err = KeyAnnotate(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--clear",
	))

	// then
	require.NoError(t, err)

	// when
	descResp, err = KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{})

	// when
	err = KeyAnnotate(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--meta", "name:key-1,role:validation",
	))

	// then
	require.NoError(t, err)

	// when
	descResp, err = KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"})
}
