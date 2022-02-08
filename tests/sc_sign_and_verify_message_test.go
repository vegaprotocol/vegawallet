package tests_test

import (
	"encoding/base64"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestSignMessage(t *testing.T) {
	// given
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	recoveryPhraseFilePath := NewFile(t, home, "recovery-phrase.txt", testRecoveryPhrase)
	walletName := vgrand.RandomStr(5)

	// when
	importWalletResp, err := WalletImport(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--recovery-phrase-file", recoveryPhraseFilePath,
		"--version", "2",
	})

	// then
	require.NoError(t, err)
	AssertImportWallet(t, importWalletResp).
		WithName(walletName).
		LocatedUnder(home)

	// given
	message := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")
	encodedMessage := base64.StdEncoding.EncodeToString(message)

	// when
	signResp, err := SignMessage(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", importWalletResp.Key.PublicKey,
		"--message", encodedMessage,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertSignMessage(t, signResp).
		WithSignature("StH82RHxjQ3yTeaSN25b6sJwAyLiq1CDvPWf0X4KIf/WTIjkunkWKn1Gq9ntCoGBfBZIyNfpPtGx0TSZsSrbCA==")

	// when
	verifyResp, err := VerifyMessage(t, []string{
		"--home", home,
		"--output", "json",
		"--pubkey", importWalletResp.Key.PublicKey,
		"--message", encodedMessage,
		"--signature", signResp.Signature,
	})

	// then
	require.NoError(t, err)
	AssertVerifyMessage(t, verifyResp).IsValid()
}

func TestSignMessageWithTaintedKey(t *testing.T) {
	// given
	home := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, home)
	recoveryPhraseFilePath := NewFile(t, home, "recovery-phrase.txt", testRecoveryPhrase)
	walletName := vgrand.RandomStr(5)

	// when
	importWalletResp, err := WalletImport(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--passphrase-file", passphraseFilePath,
		"--recovery-phrase-file", recoveryPhraseFilePath,
		"--version", "2",
	})

	// then
	require.NoError(t, err)
	AssertImportWallet(t, importWalletResp).
		WithName(walletName).
		WithPublicKey("b5fd9d3c4ad553cb3196303b6e6df7f484cf7f5331a572a45031239fd71ad8a0").
		LocatedUnder(home)

	// when
	err = KeyTaint(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", importWalletResp.Key.PublicKey,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)

	// given
	message := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")
	encodedMessage := base64.StdEncoding.EncodeToString(message)

	// when
	signResp, err := SignMessage(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", importWalletResp.Key.PublicKey,
		"--message", encodedMessage,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.EqualError(t, err, "couldn't sign message: public key is tainted")
	require.Nil(t, signResp)

	// when
	err = KeyUntaint(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", importWalletResp.Key.PublicKey,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)

	// when
	signResp, err = SignMessage(t, []string{
		"--home", home,
		"--output", "json",
		"--wallet", walletName,
		"--pubkey", importWalletResp.Key.PublicKey,
		"--message", encodedMessage,
		"--passphrase-file", passphraseFilePath,
	})

	// then
	require.NoError(t, err)
	AssertSignMessage(t, signResp).
		WithSignature("StH82RHxjQ3yTeaSN25b6sJwAyLiq1CDvPWf0X4KIf/WTIjkunkWKn1Gq9ntCoGBfBZIyNfpPtGx0TSZsSrbCA==")
}
