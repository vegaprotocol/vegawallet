package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestImportWalletV1(t *testing.T) {
	// given
	home := NewTempDir(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	mnemonicFilePath := NewFile(t, home, "mnemonic.txt", testMnemonic)

	walletName := vgrand.RandomStr(5)

	// when
	importWalletResp, err := WalletImport(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--mnemonic-file", mnemonicFilePath,
		"--version", "1",
	})

	// then
	require.NoError(t, err)
	AssertImportWallet(t, importWalletResp).
		WithName(walletName).
		LocatedUnder(home)

	// when
	walletInfoResp, err := WalletInfo(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertWalletInfo(t, walletInfoResp).
		IsHDWallet().
		WithVersion(1)

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
	require.Len(t, listKeysResp1.Keys, 0)

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
		WithoutWalletCreation().
		WithName(walletName).
		WithVersion(1).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		WithPublicKey("30ebce58d94ad37c4ff6a9014c955c20e12468da956163228cc7ec9b98d3a371").
		LocatedUnder(home)
}

func TestImportWalletV2(t *testing.T) {
	// given
	home := NewTempDir(t)

	_, passphraseFilePath := NewPassphraseFile(t, home)

	mnemonicFilePath := NewFile(t, home, "mnemonic.txt", testMnemonic)

	walletName := vgrand.RandomStr(5)

	// when
	importWalletResp, err := WalletImport(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--mnemonic-file", mnemonicFilePath,
		"--version", "2",
	})

	// then
	require.NoError(t, err)
	AssertImportWallet(t, importWalletResp).
		WithName(walletName).
		LocatedUnder(home)

	// when
	walletInfoResp, err := WalletInfo(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertWalletInfo(t, walletInfoResp).
		IsHDWallet().
		WithVersion(2)

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
	require.Len(t, listKeysResp1.Keys, 0)

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
		WithoutWalletCreation().
		WithName(walletName).
		WithVersion(2).
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		WithPublicKey("b5fd9d3c4ad553cb3196303b6e6df7f484cf7f5331a572a45031239fd71ad8a0").
		LocatedUnder(home)
}
