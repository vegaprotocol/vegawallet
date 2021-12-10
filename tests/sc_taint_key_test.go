package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestTaintKeys(t *testing.T) {
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
	err = KeyTaint(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", createWalletResp.Key.PublicKey,
	})

	// then
	require.NoError(t, err)

	// when
	descResp, err := KeyDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", createWalletResp.Key.PublicKey,
	})

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).WithTainted(true)

	// when
	err = KeyUntaint(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", createWalletResp.Key.PublicKey,
	})

	// then
	require.NoError(t, err)

	// when
	descResp, err = KeyDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", createWalletResp.Key.PublicKey,
	})

	// then
	require.NoError(t, err)
	AssertDescribeKey(t, descResp).WithTainted(false)
}
