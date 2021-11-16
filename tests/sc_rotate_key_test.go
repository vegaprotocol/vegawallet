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
	resp, err := KeyRotate(t, append(cmd,
		"--pubkey", generateKeyResp1.Key.KeyPair.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	))

	// then
	require.NoError(t, err)
	AssertKeyRotate(t, resp)
}

func TestRotateKeyFailsOnTainedPublicKey(t *testing.T) {
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
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	err = KeyTaint(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
	))

	// then
	require.NoError(t, err)

	// when
	resp, err := KeyRotate(t, append(cmd,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	))

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrPubKeyIsTainted)
}

func TestRotateKeyFailsInIsolatedWallet(t *testing.T) {
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
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	isolateKeyResp, err := KeyIsolate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertIsolateKey(t, isolateKeyResp).
		WithSpecialName(walletName, generateKeyResp.Key.KeyPair.PublicKey).
		LocatedUnder(home)

	// when
	resp, err := KeyRotate(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", isolateKeyResp.Wallet,
		"--passphrase-file", passphraseFilePath,
		"--pubkey", generateKeyResp.Key.KeyPair.PublicKey,
		"--tx-height", "20",
		"--target-height", "25",
	})

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrCantRotateKeyInIsolatedWallet)
}

func TestRotateKeyFailsOnNonExitingKey(t *testing.T) {
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
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	resp, err := KeyRotate(t, append(cmd,
		"--pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "25",
	))

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrPubKeyDoesNotExist)
}

func TestRotateKeyFailsOnNonExitingWallet(t *testing.T) {
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
	resp, err := KeyRotate(t, append(cmd,
		"--pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "25",
	))

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, wallet.ErrWalletDoesNotExists)
}

func TestRotateKeyFailsWhenTargetHeighIsLessnThanTxHeight(t *testing.T) {
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
		WithMeta(map[string]string{"name": "key-1", "role": "validation"}).
		LocatedUnder(home)

	// when
	resp, err := KeyRotate(t, append(cmd,
		"--pubkey", "nonexisting",
		"--tx-height", "20",
		"--target-height", "19",
	))

	// then
	require.Nil(t, resp)
	require.ErrorIs(t, err, flags.FlagRequireLessThanFlagError("tx-height", "target-height"))
}