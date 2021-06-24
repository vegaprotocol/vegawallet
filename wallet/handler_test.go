package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/require"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	walletpb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/wallet/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testHandler struct {
	*wallet.Handler
	ctrl  *gomock.Controller
	store *mockedStore
}

func getTestHandler(t *testing.T) *testHandler {
	ctrl := gomock.NewController(t)
	store := newMockedStore()

	h := wallet.NewHandler(store)
	return &testHandler{
		Handler: h,
		ctrl:    ctrl,
		store:   store,
	}
}

func TestHandler(t *testing.T) {
	t.Run("Creating a wallet succeeds", testHandlerCreatingWalletSucceeds)
	t.Run("Recreating a wallet with same name fails", testHandlerRecreatingWalletWithSameNameFails)
	t.Run("Recreating a wallet with same name and different passphrase fails", testHandlerRecreatingWalletWithSameNameButDifferentPassphraseFails)
	t.Run("Login to existing wallet succeeds", testHandlerLoginToExistingWalletSucceeds)
	t.Run("Login to non-existing wallet succeeds", testHandlerLoginToNonExistingWalletFails)
	t.Run("Generating new key pair succeeds", testHandlerGeneratingNewKeyPairSucceeds)
	t.Run("Generating new key pair with invalid name fails", testHandlerGeneratingNewKeyPairWithInvalidNameFails)
	t.Run("Generating new key pair without wallet fails", testHandlerGeneratingNewKeyPairWithoutWalletFails)
	t.Run("Listing public keys with invalid name fails", testHandlerListingPublicKeysWithInvalidNameFails)
	t.Run("Listing public keys without wallet fails", testHandlerListingPublicKeysWithoutWalletFails)
	t.Run("Getting public key succeeds", testHandlerGettingPublicKeySucceeds)
	t.Run("Getting public key without wallet fails", testHandlerGettingPublicKeyWithoutWalletFails)
	t.Run("Getting public key with invalid name fails", testHandlerGettingPublicKeyWithInvalidNameFails)
	t.Run("Getting non-existing public key fails", testGettingNonExistingPublicKeyFails)
	t.Run("Tainting key pair succeeds", testHandlerTaintingKeyPairSucceeds)
	t.Run("Tainting key pair with invalid name fails", testHandlerTaintingKeyPairWithInvalidNameFails)
	t.Run("Tainting key pair without wallet fails", testHandlerTaintingKeyPairWithoutWalletFails)
	t.Run("Tainting key pair that is already tainted fails", testHandlerTaintingKeyThatIsAlreadyTaintedFails)
	t.Run("Updating key pair meta succeeds", testHandlerUpdatingKeyPairMetaSucceeds)
	t.Run("Updating key pair meta with invalid passphrase fails", testHandlerUpdatingKeyPairMetaWithInvalidPassphraseFails)
	t.Run("Updating key pair meta with invalid name fails", testHandlerUpdatingKeyPairMetaWithInvalidNameFails)
	t.Run("Updating key pair meta without wallet fails", testHandlerUpdatingKeyPairMetaWithoutWalletFails)
	t.Run("Updating key pair meta with non-existing public key fails", testHandlerUpdatingKeyPairMetaWithNonExistingPublicKeyFails)
	t.Run("Get wallet path succeeds", testHandlerGettingWalletPathSucceeds)
	t.Run("Signing transaction request (v2) succeeds", testHandlerSigningTxV2Succeeds)
	t.Run("Signing transaction request (v2) with tainted key fails", testHandlerSigningTxV2WithTaintedKeyFails)
}

func testHandlerCreatingWalletSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
}

func testHandlerRecreatingWalletWithSameNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	err = h.CreateWallet(name, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletAlreadyExists.Error())
}

func testHandlerRecreatingWalletWithSameNameButDifferentPassphraseFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"
	othPassphrase := "different-passphrase"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	err = h.CreateWallet(name, othPassphrase)

	// then
	require.Error(t, err)
}

func testHandlerLoginToExistingWalletSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// then
	err = h.LoginWallet(name, passphrase)

	require.NoError(t, err)
}

func testHandlerLoginToNonExistingWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.LoginWallet(name, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
}

func testHandlerGeneratingNewKeyPairSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)
	
	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keys, err := h.ListPublicKeys(name)

	// then
	require.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, key, keys[0].Key)
	assert.False(t, keys[0].Tainted)
}

func testHandlerGeneratingNewKeyPairWithInvalidNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	otherName := "bad name"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(otherName, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerGeneratingNewKeyPairWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"
	passphrase := "Th1isisasecurep@ssphraseinnit"

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerListingPublicKeysWithInvalidNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	otherName := "bad name"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.ListPublicKeys(otherName)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerListingPublicKeysWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"

	// when
	key, err := h.ListPublicKeys(name)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Empty(t, key)
}

func testHandlerGettingPublicKeyWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"

	// when
	key, err := h.GetPublicKey(name, name)

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

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(name, key)

	require.NoError(t, err)
	assert.Equal(t, key, keyPair.Key)
}

func testHandlerGettingPublicKeyWithInvalidNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	otherName := "bad name"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(otherName, key)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Nil(t, keyPair)
}

func testGettingNonExistingPublicKeyFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(name, "non-existing-pub-key")
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Nil(t, keyPair)
}

func testHandlerTaintingKeyPairSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(name, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// when
	err = h.TaintKey(name, key, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, key).Tainted)
}

func testHandlerTaintingKeyPairWithInvalidNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	otherName := "other name"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(name, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// when
	err = h.TaintKey(otherName, key, passphrase)

	// then
	assert.Error(t, err)
}

func testHandlerUpdatingKeyPairMetaWithNonExistingPublicKeyFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	err = h.TaintKey(name, "non-existing-pub-key", passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
}

func testHandlerTaintingKeyPairWithoutWalletFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.TaintKey(name, "non-existing-pub-key", passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
}

func testHandlerTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	keyPair, err := h.GetPublicKey(name, key)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.False(t, keyPair.Tainted)

	// when
	err = h.TaintKey(name, key, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, key).Tainted)

	// when
	err = h.TaintKey(name, key, passphrase)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyAlreadyTainted.Error())
}

func testHandlerUpdatingKeyPairMetaSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	err = h.UpdateMeta(name, key, passphrase, meta)

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

	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	err = h.UpdateMeta(name, key, othPassphrase, meta)

	// then
	assert.Error(t, err)
	assert.Len(t, h.store.GetKey(name, key).Meta, 0)
}

func testHandlerUpdatingKeyPairMetaWithInvalidNameFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"
	otherName := "other name"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	key, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, key)

	// when
	err = h.UpdateMeta(otherName, key, passphrase, meta)

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

	pubKey := "non-existing-public-key"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	err := h.UpdateMeta(name, pubKey, passphrase, meta)

	// then
	assert.Error(t, err)
}

func testHandlerGettingWalletPathSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	name := "jeremy"

	// when
	path, err := h.GetWalletPath(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, path)
}

func testHandlerSigningTxV2Succeeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	passphrase := "Th1isisasecurep@ssphraseinnit"
	name := "jeremy"

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	pubKey, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// given
	req := &walletpb.SubmitTransactionRequest{
		PubKey: pubKey,
		Command: &walletpb.SubmitTransactionRequest_OrderCancellation{
			OrderCancellation: &commandspb.OrderCancellation{},
		},
	}

	// when
	tx, err := h.SignTxV2(name, req, 42)

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

	// when
	err := h.CreateWallet(name, passphrase)

	// then
	require.NoError(t, err)

	// when
	pubKey, err := h.SecureGenerateKeyPair(name, passphrase)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, pubKey)

	// when
	err = h.TaintKey(name, pubKey, passphrase)

	// then
	require.NoError(t, err)
	assert.True(t, h.store.GetKey(name, pubKey).Tainted)

	// given
	req := &walletpb.SubmitTransactionRequest{
		PubKey: pubKey,
		Command: &walletpb.SubmitTransactionRequest_OrderCancellation{
			OrderCancellation: &commandspb.OrderCancellation{},
		},
	}

	// when
	tx, err := h.SignTxV2(name, req, 42)

	// then
	assert.Error(t, err)
	assert.Nil(t, tx)
}
