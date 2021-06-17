package wallet_test

import (
	"errors"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/mocks"
	"github.com/stretchr/testify/require"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	walletpb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/wallet/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*wallet.Handler
	ctrl  *gomock.Controller
	auth  *mocks.MockAuth
	store *mockedStore
}

func getTestHandler(t *testing.T) *testHandler {
	ctrl := gomock.NewController(t)
	auth := mocks.NewMockAuth(ctrl)
	store := newMockedStore()

	h := wallet.NewHandler(auth, store)
	return &testHandler{
		Handler: h,
		ctrl:    ctrl,
		auth:    auth,
		store:   store,
	}
}

func TestHandler(t *testing.T) {
	t.Run("Creating a wallet succeeds", testHandlerCreatingWalletSucceeds)
	t.Run("Recreating a wallet with same name fails", testHandlerRecreatingWalletWithSameNameFails)
	t.Run("Recreating a wallet with same name and different passphrase fails", testHandlerRecreatingWalletWithSameNameButDifferentPassphraseFails)
	t.Run("Login to existing wallet succeeds", testHandlerLoginToExistingWalletSucceeds)
	t.Run("Login to non-existing wallet succeeds", testHandlerLoginToNonExistingWalletFails)
	t.Run("Revoking right token success", testHandlerRevokingRightTokenSucceeds)
	t.Run("Revoking invalid token fails", testHandlerRevokingInvalidTokenFails)
	t.Run("Generating new key pair succeeds", testHandlerGeneratingNewKeyPairSucceeds)
	t.Run("Generating new key pair with invalid token fails", testHandlerGeneratingNewKeyPairWithInvalidTokenFails)
	t.Run("Generating new key pair without wallet fails", testHandlerGeneratingNewKeyPairWithoutWalletFails)
	t.Run("Listing public keys with invalid token fails", testHandlerListingPublicKeysWithInvalidTokenFails)
	t.Run("Listing public keys without wallet fails", testHandlerListingPublicKeysWithoutWalletFails)
	t.Run("Getting public key succeeds", testHandlerGettingPublicKeySucceeds)
	t.Run("Getting public key without wallet fails", testHandlerGettingPublicKeyWithoutWalletFails)
	t.Run("Getting public key with invalid token fails", testHandlerGettingPublicKeyWithInvalidTokenFails)
	t.Run("Getting non-existing public key fails", testGettingNonExistingPublicKeyFails)
	t.Run("Tainting key pair succeeds", testHandlerTaintingKeyPairSucceeds)
	t.Run("Tainting key pair with invalid token fails", testHandlerTaintingKeyPairWithInvalidTokenFails)
	t.Run("Tainting key pair without wallet fails", testHandlerTaintingKeyPairWithoutWalletFails)
	t.Run("Tainting key pair that is already tainted fails", testHandlerTaintingKeyThatIsAlreadyTaintedFails)
	t.Run("Updating key pair meta succeeds", testHandlerUpdatingKeyPairMetaSucceeds)
	t.Run("Updating key pair meta with invalid passphrase fails", testHandlerUpdatingKeyPairMetaWithInvalidPassphraseFails)
	t.Run("Updating key pair meta with invalid token fails", testHandlerUpdatingKeyPairMetaWithInvalidTokenFails)
	t.Run("Updating key pair meta without wallet fails", testHandlerUpdatingKeyPairMetaWithoutWalletFails)
	t.Run("Updating key pair meta with non-existing public key fails", testHandlerUpdatingKeyPairMetaWithNonExisitingPublicKeyFails)
	t.Run("Get wallet path succeeds", testHandlerGettingWalletPathSucceeds)
	t.Run("Get wallet path with invalid token fails", testHandlerGettingWalletPathWithInvalidTokenFails)
	t.Run("Signing transaction request (v2) succeeds", testHandlerSigningTxV2Succeeds)
	t.Run("Signing transaction request (v2) with tainted key fails", testHandlerSigningTxV2WithTaintedKeyFails)
}

func testHandlerCreatingWalletSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(gomock.Any()).
		Return("some fake token", nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)
}

func testHandlerRecreatingWalletWithSameNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).
		Return("some fake token", nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().NewSession(name).
		Times(0)

	// when
	returnedToken, err = h.CreateWallet(name, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletAlreadyExists.Error())
	assert.Empty(t, returnedToken)
}

func testHandlerRecreatingWalletWithSameNameButDifferentPassphraseFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"
	othPassphrase := "different-passphrase"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).
		Return("some fake token", nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().NewSession(name).
		Times(0)

	// when
	returnedToken, err = h.CreateWallet(name, othPassphrase)

	// then
	require.Error(t, err)
	assert.Empty(t, returnedToken)
}

func testHandlerLoginToExistingWalletSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().NewSession(name).
		Return(token, nil)

	// then
	returnedToken, err = h.LoginWallet(name, passphrase)

	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)
}

func testHandlerLoginToNonExistingWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// setup
	h.auth.EXPECT().NewSession(gomock.Any()).
		Times(0)

	// when
	returnedToken, err := h.LoginWallet(name, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, returnedToken)
}

func testHandlerRevokingRightTokenSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().Revoke(token).Times(1).
		Return(nil)

	// when
	err = h.RevokeToken(returnedToken)

	// then
	require.NoError(t, err)
}

func testHandlerRevokingInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "bad token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().Revoke(othToken).Times(1).
		Return(errors.New(othToken))

	// when
	err = h.RevokeToken(othToken)

	// then
	assert.EqualError(t, err, othToken)
}

func testHandlerGeneratingNewKeyPairSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	keys, err := h.ListPublicKeys(returnedToken)

	// then
	require.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0].Pub)
	assert.False(t, keys[0].Tainted)
}

func testHandlerGeneratingNewKeyPairWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "bad token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(othToken).
		Return("", errors.New(othToken))

	// when
	key, err := h.GenerateKeypair(othToken, passphrase)

	// then
	assert.EqualError(t, err, othToken)
	assert.Empty(t, key)
}

func testHandlerGeneratingNewKeyPairWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(token, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerListingPublicKeysWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "bad token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(othToken).
		Return("", errors.New(othToken))

	// when
	key, err := h.ListPublicKeys(othToken)

	// then
	assert.EqualError(t, err, othToken)
	assert.Empty(t, key)
}

func testHandlerListingPublicKeysWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	key, err := h.ListPublicKeys(token)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerGettingPublicKeyWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	key, err := h.GetPublicKey(token, name)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerGettingPublicKeySucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	keyPair, err := h.GetPublicKey(token, key)

	require.NoError(t, err)
	assert.Equal(t, key, keyPair.Pub)
	assert.Empty(t, keyPair.Priv)
}

func testHandlerGettingPublicKeyWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "bad token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(othToken).
		Return("", errors.New(othToken))

	// when
	keyPair, err := h.GetPublicKey(othToken, key)

	// then
	assert.EqualError(t, err, othToken)
	assert.Nil(t, keyPair)
}

func testGettingNonExistingPublicKeyFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	keyPair, err := h.GetPublicKey(token, "nonexistantpubkey")
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Nil(t, keyPair)
}

func testHandlerTaintingKeyPairSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	keyPair, err := h.GetPublicKey(token, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.TaintKey(token, key, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, key).Tainted)
}

func testHandlerTaintingKeyPairWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "other token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	keyPair, err := h.GetPublicKey(token, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// setup
	h.auth.EXPECT().VerifyToken(othToken).
		Return("", errors.New(othToken))

	// when
	err = h.TaintKey(othToken, key, passphrase)

	// then
	assert.Error(t, err)
}

func testHandlerUpdatingKeyPairMetaWithNonExisitingPublicKeyFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	err = h.TaintKey(token, "non-existing-pub-key", passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
}

func testHandlerTaintingKeyPairWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	err := h.TaintKey(token, "non-existing-pub-key", passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
}

func testHandlerTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	keyPair, err := h.GetPublicKey(token, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.TaintKey(token, key, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, key).Tainted)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.TaintKey(token, key, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyAlreadyTainted.Error())
}

func testHandlerUpdatingKeyPairMetaSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.UpdateMeta(token, key, passphrase, meta)

	// then
	require.NoError(t, err)
	updatedKp := h.store.GetKey(name, key)
	assert.Len(t, updatedKp.Meta, 1)
	assert.Equal(t, updatedKp.Meta[0].Key, "primary")
	assert.Equal(t, updatedKp.Meta[0].Value, "yes")
}

func testHandlerUpdatingKeyPairMetaWithInvalidPassphraseFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	othPassphrase := "other-passphrase"
	name := "jeremy"
	token := "some fake token"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.UpdateMeta(token, key, othPassphrase, meta)

	// then
	assert.Error(t, err)
	assert.Len(t, h.store.GetKey(name, key).Meta, 0)
}

func testHandlerUpdatingKeyPairMetaWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	othToken := "other token"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	key, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return("", errors.New(othToken))

	// when
	err = h.UpdateMeta(token, key, passphrase, meta)

	// then
	assert.Error(t, err)
	assert.Len(t, h.store.GetKey(name, key).Meta, 0)
}

func testHandlerUpdatingKeyPairMetaWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"
	pubKey := "non-existing-public-key"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	err := h.UpdateMeta(token, pubKey, passphrase, meta)

	// then
	assert.Error(t, err)
}

func testHandlerGettingWalletPathSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return(name, nil)

	// when
	path, err := h.GetWalletPath(token)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, path)
}

func testHandlerGettingWalletPathWithInvalidTokenFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	token := "some fake token"

	// setup
	h.auth.EXPECT().VerifyToken(token).
		Return("", errors.New(token))

	// when
	path, err := h.GetWalletPath(token)

	// then
	assert.Error(t, err)
	assert.Empty(t, path)
}

func testHandlerSigningTxV2Succeeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	pubKey, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// given
	req := walletpb.SubmitTransactionRequest{
		PubKey: pubKey,
		Command: &walletpb.SubmitTransactionRequest_OrderCancellation{
			OrderCancellation: &commandspb.OrderCancellation{},
		},
	}

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	tx, err := h.SignTxV2(returnedToken, req, 42)

	// then
	require.NoError(t, err)
	assert.Equal(t, uint32(1), tx.Version)
	assert.NotEmpty(t, tx.From)
	assert.NotEmpty(t, tx.InputData)
	assert.NotNil(t, tx.Signature)
	key := h.store.GetKey(name, pubKey)
	assert.Equal(t, key.Algorithm.Version(), tx.Signature.Version)
	assert.Equal(t, key.Algorithm.Name(), tx.Signature.Algo)
	assert.NotEmpty(t, tx.Signature.Value)
}

func testHandlerSigningTxV2WithTaintedKeyFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	token := "some fake token"

	// setup
	h.auth.EXPECT().NewSession(name).Times(1).
		Return(token, nil)

	// when
	returnedToken, err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	assert.Equal(t, token, returnedToken)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	pubKey, err := h.GenerateKeypair(returnedToken, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	err = h.TaintKey(token, pubKey, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, pubKey).Tainted)

	// given
	req := walletpb.SubmitTransactionRequest{
		PubKey: pubKey,
		Command: &walletpb.SubmitTransactionRequest_OrderCancellation{
			OrderCancellation: &commandspb.OrderCancellation{},
		},
	}

	// setup
	h.auth.EXPECT().VerifyToken(returnedToken).
		Return(name, nil)

	// when
	tx, err := h.SignTxV2(returnedToken, req, 42)

	// then
	assert.Error(t, err)
	assert.Nil(t, tx)
}
