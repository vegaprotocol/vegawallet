package tests_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/require"
)

func TestRotateKeySucceeds(t *testing.T) {
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
	generateKeyResp, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-2,role:validation",
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithMeta(map[string]string{"name": "key-2", "role": "validation"})

	// when
	resp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--current-pubkey", createWalletResp.Key.PublicKey,
		"--new-pubkey", generateKeyResp.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.NoError(t, err)
	AssertKeyRotate(t, resp)
}

func TestRotateKeyFailsOnTaintedPublicKey(t *testing.T) {
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
	generateKeyResp, err := KeyGenerate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--meta", "name:key-2,role:validation",
	})

	// then
	require.NoError(t, err)
	AssertGenerateKey(t, generateKeyResp).
		WithMeta(map[string]string{"name": "key-2", "role": "validation"})

	// when
	err = KeyTaint(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", generateKeyResp.PublicKey,
	})

	// then
	require.NoError(t, err)

	// when
	resp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--current-pubkey", createWalletResp.Key.PublicKey,
		"--new-pubkey", generateKeyResp.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrPubKeyIsTainted)
}

func TestRotateKeyFailsInIsolatedWallet(t *testing.T) {
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
	isolateKeyResp, err := KeyIsolate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", createWalletResp.Key.PublicKey,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertIsolateKey(t, isolateKeyResp).
		WithSpecialName(walletName, createWalletResp.Key.PublicKey).
		LocatedUnder(home)

	// when
	resp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", isolateKeyResp.Wallet,
		"--passphrase-file", passphraseFilePath,
		"--new-pubkey", createWalletResp.Key.PublicKey,
		"--current-pubkey", "current-public-key",
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrCantRotateKeyInIsolatedWallet)
}

func TestRotateKeyFailsOnNonExitingNewPublicKey(t *testing.T) {
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
	KeyRotateResp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--current-pubkey", "current-public-key",
		"--new-pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, KeyRotateResp)
	require.ErrorIs(t, err, wallet.ErrPubKeyDoesNotExist)
}

func TestRotateKeyFailsOnNonExitingCurrentPublicKey(t *testing.T) {
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
	keyRotateResp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--current-pubkey", "nonexisting",
		"--new-pubkey", createWalletResp.Key.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, keyRotateResp)
	require.ErrorIs(t, err, wallet.ErrPubKeyDoesNotExist)
}

func TestRotateKeyFailsOnNonExitingWallet(t *testing.T) {
	// given
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	walletName := vgrand.RandomStr(5)

	// when
	keyRotateResp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--new-pubkey", "nonexisting",
		"--current-pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, keyRotateResp)
	require.ErrorIs(t, err, wallet.ErrWalletDoesNotExists)
}

func TestRotateKeyFailsWhenTargetHeighIsLessnThanTxHeight(t *testing.T) {
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
	keyRotateResp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--new-pubkey", "nonexisting",
		"--current-pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "19",
	})

	// then
	require.Nil(t, keyRotateResp)
	require.ErrorIs(t, err, flags.FlagRequireLessThanFlagError("tx-height", "target-height"))
}
