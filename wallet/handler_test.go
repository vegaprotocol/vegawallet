package wallet_test

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"code.vegaprotocol.io/go-wallet/wallet/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testHandler struct {
	*wallet.Handler
	ctrl    *gomock.Controller
	auth    *mocks.MockAuth
	rootDir string
}

func getTestHandler(t *testing.T) *testHandler {
	ctrl := gomock.NewController(t)
	auth := mocks.NewMockAuth(ctrl)
	rootPath := rootDir()
	fsutil.EnsureDir(rootPath)
	wallet.EnsureBaseFolder(rootPath)

	h := wallet.NewHandler(zap.NewNop(), auth, rootPath)
	return &testHandler{
		Handler: h,
		ctrl:    ctrl,
		auth:    auth,
		rootDir: rootPath,
	}
}

func TestHandler(t *testing.T) {
	t.Run("create a wallet success then login", testHandlerCreateWalletThenLogin)
	t.Run("create a wallet failure - already exists", testHandlerCreateWalletFailureAlreadyExists)
	t.Run("login failure on non wallet", testHandlerLoginFailureOnNonCreatedWallet)
	t.Run("revoke token success", testHandlerRevokeTokenSuccess)
	t.Run("revoke token failure", testHandlerRevokeTokenFailure)
	t.Run("generate keypair success and list public keys", testVerifyTokenSuccess)
	t.Run("generate keypair failure - invalid token", testVerifyTokenInvalidToken)
	t.Run("generate keypair failure - wallet not found", testVerifyTokenWalletNotFound)
	t.Run("list public key failure - invalid token", testListPubInvalidToken)
	t.Run("list public key failure - wallet not found", testListPubWalletNotFound)
	t.Run("get public key failure - success", testGetPubSuccess)
	t.Run("get public key failure - wallet not found", testGetPubWalletNotFound)
	t.Run("get public key failure - invalid token", testGetPubInvalidToken)
	t.Run("get public key failure - key not found", testGetPubKeyNotFound)
	t.Run("sign tx - success", testSignTxSuccess)
	t.Run("sign tx - failure key tainted", testSignTxFailure)
	t.Run("taint key - success", testTaintKeySuccess)
	t.Run("taint key failure - invalid token", testTaintKeyInvalidToken)
	t.Run("taint key failure - wallet not found", testTaintKeyWalletNotFound)
	t.Run("taint key failure - already tainted", testTaintKeyAlreadyFailAlreadyTainted)
	t.Run("update meta failure - pub key does not exists", testTaintKeyPubKeyDoesNotExists)
	t.Run("update meta - success", testUpdateMetaSuccess)
	t.Run("update meta - failure invalid passphrase", testUpdateMetaFailureInvalidPassphrase)
	t.Run("update meta failure - invalid token", testUpdateMetaInvalidToken)
	t.Run("update meta taint key failure - wallet not found", testUpdateMetaWalletNotFound)
	t.Run("update meta failure - pub key does not exists", testUpdateMetaPubKeyDoesNotExists)
}

func testHandlerCreateWalletThenLogin(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	h.auth.EXPECT().NewSession(gomock.Any()).Times(2).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	tok, err = h.LoginWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testHandlerCreateWalletFailureAlreadyExists(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	// create the wallet once.
	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// try to create it again
	tok, err = h.CreateWallet("jeremy", "We can use a d1fferent passphrase yo!")
	assert.EqualError(t, err, wallet.ErrWalletAlreadyExists.Error())
	assert.Empty(t, tok)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testHandlerLoginFailureOnNonCreatedWallet(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	tok, err := h.LoginWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, tok)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testHandlerRevokeTokenSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	h.auth.EXPECT().Revoke(gomock.Any()).Times(1).
		Return(nil)
	err = h.RevokeToken(tok)
	assert.NoError(t, err)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testHandlerRevokeTokenFailure(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	h.auth.EXPECT().Revoke(gomock.Any()).Times(1).
		Return(errors.New("bad token"))
	err = h.RevokeToken(tok)
	assert.EqualError(t, err, "bad token")

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testVerifyTokenSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(2).
		Return("jeremy", nil)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// now make sure we have the new key saved
	keys, err := h.ListPublicKeys(tok)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0].Pub)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testVerifyTokenInvalidToken(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("", errors.New("bad token"))

	key, err := h.GenerateKeypair("yolo token", "whatever")
	assert.EqualError(t, err, "bad token")
	assert.Empty(t, key)

	assert.NoError(t, os.RemoveAll(h.rootDir))

}

// this should never happend but beeeh....
func testVerifyTokenWalletNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("jeremy", nil)

	key, err := h.GenerateKeypair("yolo token", "whatever")
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testListPubInvalidToken(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("", errors.New("bad token"))

	key, err := h.ListPublicKeys("yolo token")
	assert.EqualError(t, err, "bad token")
	assert.Empty(t, key)

	assert.NoError(t, os.RemoveAll(h.rootDir))

}

// this should never happend but beeeh....
func testListPubWalletNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("jeremy", nil)

	key, err := h.ListPublicKeys("yolo token")
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testGetPubWalletNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("jeremy", nil)

	key, err := h.GetPublicKey("yolo token", "1122aabb")
	assert.Empty(t, key)
	assert.Equal(t, err, wallet.ErrWalletDoesNotExists)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testGetPubSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("walletname", "s1cur@Epassphrase")
	if err != nil {
		t.Fatal("Failed to create wallet")
	}
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("walletname", nil)

	// Create a keypair to be ignored, *not* returned later
	pubKey, err := h.GenerateKeypair(tok, "s1cur@Epassphrase")
	assert.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// Create a keypair to be retrieved
	pubKey, err = h.GenerateKeypair(tok, "s1cur@Epassphrase")
	assert.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// now make sure we have the new key saved
	key, err := h.GetPublicKey(tok, pubKey)
	if err != nil {
		t.Fatal("Failed to get public key")
	}
	assert.Equal(t, pubKey, key.Pub)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testGetPubInvalidToken(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	h.auth.EXPECT().VerifyToken(gomock.Any()).Times(1).
		Return("", errors.New("bad token"))

	key, err := h.GetPublicKey("yolo token", "1122aabb")
	assert.EqualError(t, err, "bad token")
	assert.Empty(t, key)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testGetPubKeyNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("walletname", "s1cur@Epassphrase")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("walletname", nil)

	// Create a keypair to be ignored, *not* returned later
	pubKey, err := h.GenerateKeypair(tok, "s1cur@Epassphrase")
	assert.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// now make sure this key is not returned
	key, err := h.GetPublicKey(tok, "nonexistantpubkey")
	assert.Nil(t, key)
	assert.Equal(t, err, wallet.ErrPubKeyDoesNotExist)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testSignTxSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	message := "hello world."

	keyBytes, _ := hex.DecodeString(key)

	signedBundle, err := h.SignTx(tok, base64.StdEncoding.EncodeToString([]byte(message)), key)
	assert.NoError(t, err)

	// verify signature then
	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	assert.NoError(t, err)

	v, err := alg.Verify(keyBytes, signedBundle.Tx, signedBundle.Sig.Sig)
	assert.NoError(t, err)
	assert.True(t, v)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testSignTxFailure(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// taint the key
	err = h.TaintKey(tok, key, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)

	message := "hello world."
	_, err = h.SignTx(tok, base64.StdEncoding.EncodeToString([]byte(message)), key)
	assert.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testTaintKeySuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// taint the key
	err = h.TaintKey(tok, key, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)

	// now make sure we have the new key saved
	keys, err := h.ListPublicKeys(tok)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0].Pub)
	assert.True(t, keys[0].Tainted)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testTaintKeyInvalidToken(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("", errors.New("invalid token"))

	// taint the key
	err := h.TaintKey("some token", "some key", "Th1isisasecurep@ssphraseinnit")
	assert.EqualError(t, err, "invalid token")

	assert.NoError(t, os.RemoveAll(h.rootDir))

}
func testTaintKeyPubKeyDoesNotExists(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// taint the key
	err = h.TaintKey(tok, "some key", "Th1isisasecurep@ssphraseinnit")
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testTaintKeyWalletNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// taint the key
	err := h.TaintKey("some token", "some key", "Th1isisasecurep@ssphraseinnit")
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testTaintKeyAlreadyFailAlreadyTainted(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// taint the key
	err = h.TaintKey(tok, key, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)

	// taint the key again which produce an error
	err = h.TaintKey(tok, key, "Th1isisasecurep@ssphraseinnit")
	assert.Error(t, err, wallet.ErrPubKeyAlreadyTainted)

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testUpdateMetaSuccess(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// add meta
	err = h.UpdateMeta(tok, key, "Th1isisasecurep@ssphraseinnit", []wallet.Meta{wallet.Meta{Key: "primary", Value: "yes"}})
	assert.NoError(t, err)

	// now make sure we have the new key saved
	keys, err := h.ListPublicKeys(tok)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0].Pub)
	assert.Len(t, keys[0].Meta, 1)
	assert.Equal(t, keys[0].Meta[0].Key, "primary")
	assert.Equal(t, keys[0].Meta[0].Value, "yes")

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testUpdateMetaFailureInvalidPassphrase(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	key, err := h.GenerateKeypair(tok, "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, key)

	// add meta
	err = h.UpdateMeta(tok, key, "this is the wrong passphrase", []wallet.Meta{wallet.Meta{Key: "primary", Value: "yes"}})
	assert.EqualError(t, err, "cipher: message authentication failed")
}

func testUpdateMetaInvalidToken(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("", errors.New("invalid token"))

	// taint the key
	err := h.UpdateMeta("some token", "some key", "Th1isisasecurep@ssphraseinnit", []wallet.Meta{})
	assert.EqualError(t, err, "invalid token")

	assert.NoError(t, os.RemoveAll(h.rootDir))

}

func testUpdateMetaPubKeyDoesNotExists(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// first create the wallet
	h.auth.EXPECT().NewSession(gomock.Any()).Times(1).
		Return("some fake token", nil)

	tok, err := h.CreateWallet("jeremy", "Th1isisasecurep@ssphraseinnit")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// update meta
	err = h.UpdateMeta(tok, "some key", "Th1isisasecurep@ssphraseinnit", []wallet.Meta{})
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())

	assert.NoError(t, os.RemoveAll(h.rootDir))
}

func testUpdateMetaWalletNotFound(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// then start the test
	h.auth.EXPECT().VerifyToken(gomock.Any()).AnyTimes().
		Return("jeremy", nil)

	// taint the key
	err := h.UpdateMeta("some token", "some key", "Th1isisasecurep@ssphraseinnit", []wallet.Meta{})
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())

	assert.NoError(t, os.RemoveAll(h.rootDir))
}
