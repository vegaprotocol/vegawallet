package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestTaintKeys(t *testing.T) {
	// given
	home := NewTempDir(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	walletName := vgrand.RandomStr(5)

	cmd := []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	}

	// when
	generateKeyResp1, err := KeyGenerate(t, append(cmd,
		"--meta", "name:key-1,role:validation",
	))

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp1).
		WithWalletCreation().
		WithName(walletName).
		WithVersion(2).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	err = KeyTaint(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.PublicKey,
	))

	// then
	require.NoError(t, err)

	// when
	descResp, err := KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).WithTainted(true)

	// when
	err = KeyUntaint(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.PublicKey,
	))

	// then
	require.NoError(t, err)

	// when
	descResp, err = KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).WithTainted(false)
}
