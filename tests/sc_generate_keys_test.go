package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndListKeys(t *testing.T) {
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
	generateKeyResp1, err := KeyGenerate(t, append(cmd,
		"--meta", "name:key-1,role:validation",
	))

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp1).
		WithWalletCreation().
		WithName(walletName).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	descResp, err := KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		WithAlgorithm("vega/ed25519", 1).
		WithTainted(false)

	// when
	generateKeyResp2, err := KeyGenerate(t, cmd)

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp2).
		WithoutWalletCreation().
		WithName(walletName).
		WithMeta(map[string]string{"name": DefaultMetaName(t, walletName, 2)}).
		LocatedUnder(home)

	// when
	descResp, err = KeyDescribe(t, append(cmd,
		"--pubkey", generateKeyResp2.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).
		WithMeta(map[string]string{"name": DefaultMetaName(t, walletName, 2)}).
		WithAlgorithm("vega/ed25519", 1).
		WithTainted(false)
}
