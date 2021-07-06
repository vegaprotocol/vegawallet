package wallet_test

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerDeprecated(t *testing.T) {
	t.Run("sign tx - success", testSignTxSuccess)
	t.Run("sign tx - failure key tainted", testSignTxFailure)
}

func testSignTxSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	err := h.CreateWallet(name, passphrase)
	require.NoError(t, err)

	key, err := h.SecureGenerateKeyPair(name, passphrase)
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	message := "hello world."

	keyBytes, _ := hex.DecodeString(key)

	signedBundle, err := h.SignTx(name, base64.StdEncoding.EncodeToString([]byte(message)), key, 42)
	require.NoError(t, err)

	// verify signature then
	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	require.NoError(t, err)

	v, err := alg.Verify(keyBytes, signedBundle.Tx, signedBundle.Sig.Sig)
	require.NoError(t, err)
	assert.True(t, v)
}

func testSignTxFailure(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"

	err := h.CreateWallet(name, passphrase)
	require.NoError(t, err)

	key, err := h.SecureGenerateKeyPair(name, passphrase)
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// taint the key
	err = h.TaintKey(name, key, passphrase)
	require.NoError(t, err)

	message := "hello world."
	_, err = h.SignTx(name, base64.StdEncoding.EncodeToString([]byte(message)), key, 42)
	assert.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())
}
