package wallet_test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"testing"

	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallet/mocks"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnotateKey(t *testing.T) {
	t.Run("Annotating an existing key succeeds", testAnnotatingKeySucceeds)
}

func testAnnotatingKeySucceeds(t *testing.T) {
	tcs := []struct {
		name     string
		metadata []wallet.Meta
	}{
		{
			name: "with metadata",
			metadata: []wallet.Meta{
				{Key: "name", Value: "my-wallet"},
				{Key: "role", Value: "validation"},
			},
		}, {
			name:     "without metadata",
			metadata: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// given
			w := newWalletWithKey(t)
			kp := w.ListKeyPairs()[0]

			req := &wallet.AnnotateKeyRequest{
				Wallet:     w.Name(),
				PubKey:     kp.PublicKey(),
				Metadata:   tc.metadata,
				Passphrase: "passphrase",
			}

			// setup
			store := handlerMocks(tt)
			store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
			store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)
			store.EXPECT().SaveWallet(w, req.Passphrase).Times(1).Return(nil)

			// when
			err := wallet.AnnotateKey(store, req)

			// then
			require.NoError(tt, err)
			assert.Equal(tt, req.Metadata, w.ListKeyPairs()[0].Meta())
		})
	}
}

func TestGenerateKey(t *testing.T) {
	t.Run("Generating keys in non-existing wallet succeeds", testGenerateKeyInNonExistingWalletSucceeds)
	t.Run("Generating keys in existing wallet succeeds", testGenerateKeyInExistingWalletSucceeds)
}

func testGenerateKeyInNonExistingWalletSucceeds(t *testing.T) {
	// given
	req := &wallet.GenerateKeyRequest{
		Wallet: "my-wallet",
		Metadata: []wallet.Meta{
			{Key: "name", Value: "my-wallet"},
			{Key: "role", Value: "validation"},
		},
		Passphrase: "passphrase",
	}

	// setup
	var generatedWallet wallet.Wallet
	captureWallet := func(w wallet.Wallet, passphrase string) error {
		generatedWallet = w
		return nil
	}
	fakePath := "/path/to/wallets/my-wallet"
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWalletPath(req.Wallet).Times(1).Return(fakePath)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(1).DoAndReturn(captureWallet)

	// when
	resp, err := wallet.GenerateKey(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	// verify generated wallet
	assert.Equal(t, req.Wallet, generatedWallet.Name())
	assert.Len(t, generatedWallet.ListKeyPairs(), 1)
	keyPair := generatedWallet.ListKeyPairs()[0]
	assert.Equal(t, req.Metadata, keyPair.Meta())
	// verify response
	assert.Equal(t, req.Wallet, resp.Wallet.Name)
	assert.Equal(t, fakePath, resp.Wallet.FilePath)
	assert.NotEmpty(t, resp.Wallet.Mnemonic)
	assert.Equal(t, keyPair.PublicKey(), resp.Key.KeyPair.PublicKey)
	assert.Equal(t, keyPair.PrivateKey(), resp.Key.KeyPair.PrivateKey)
	assert.Equal(t, keyPair.AlgorithmName(), resp.Key.Algorithm.Name)
	assert.Equal(t, keyPair.AlgorithmVersion(), resp.Key.Algorithm.Version)
	assert.Equal(t, keyPair.Meta(), resp.Key.Meta)
}

func testGenerateKeyInExistingWalletSucceeds(t *testing.T) {
	// given
	w := newWallet(t)
	req := &wallet.GenerateKeyRequest{
		Wallet: w.Name(),
		Metadata: []wallet.Meta{
			{Key: "name", Value: "my-wallet"},
			{Key: "role", Value: "validation"},
		},
		Passphrase: "passphrase",
	}

	// setup
	fakePath := "/path/to/wallets/my-wallet"
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)
	store.EXPECT().GetWalletPath(req.Wallet).Times(1).Return(fakePath)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(1).Return(nil)

	// when
	resp, err := wallet.GenerateKey(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	// verify updated wallet
	assert.Equal(t, req.Wallet, w.Name())
	require.Len(t, w.ListKeyPairs(), 1)
	keyPair := w.ListKeyPairs()[0]
	assert.Equal(t, req.Metadata, keyPair.Meta())
	// verify response
	assert.Equal(t, req.Wallet, resp.Wallet.Name)
	assert.Equal(t, fakePath, resp.Wallet.FilePath)
	assert.Empty(t, resp.Wallet.Mnemonic)
	assert.Equal(t, keyPair.PublicKey(), resp.Key.KeyPair.PublicKey)
	assert.Equal(t, keyPair.PrivateKey(), resp.Key.KeyPair.PrivateKey)
	assert.Equal(t, keyPair.AlgorithmName(), resp.Key.Algorithm.Name)
	assert.Equal(t, keyPair.AlgorithmVersion(), resp.Key.Algorithm.Version)
	assert.Equal(t, keyPair.Meta(), resp.Key.Meta)
}

func TestTaintKey(t *testing.T) {
	t.Run("Tainting key succeeds", testTaintingKeySucceeds)
	t.Run("Tainting key of non-existing wallet fails", testTaintingKeyOfNonExistingWalletFails)
}

func testTaintingKeySucceeds(t *testing.T) {
	// given
	w := newWalletWithKey(t)
	kp := w.ListKeyPairs()[0]

	req := &wallet.TaintKeyRequest{
		Wallet:     w.Name(),
		PubKey:     kp.PublicKey(),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)
	store.EXPECT().SaveWallet(w, req.Passphrase).Times(1).Return(nil)

	// when
	err := wallet.TaintKey(store, req)

	// then
	require.NoError(t, err)
	assert.True(t, w.ListKeyPairs()[0].IsTainted())
}

func testTaintingKeyOfNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.TaintKeyRequest{
		Wallet:     vgrand.RandomStr(5),
		PubKey:     vgrand.RandomStr(25),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(0)

	// when
	err := wallet.TaintKey(store, req)

	// then
	require.Error(t, err)
}

func TestUntaintKey(t *testing.T) {
	t.Run("Untainting key succeeds", testUntaintingKeySucceeds)
	t.Run("Untainting key of non-existing wallet fails", testUntaintingKeyOfNonExistingWalletFails)
}

func testUntaintingKeySucceeds(t *testing.T) {
	// given
	w := newWalletWithKey(t)
	kp := w.ListKeyPairs()[0]
	err := w.TaintKey(kp.PublicKey())
	if err != nil {
		t.Fatalf("couldn't taint key: %v", err)
	}

	req := &wallet.UntaintKeyRequest{
		Wallet:     w.Name(),
		PubKey:     kp.PublicKey(),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)
	store.EXPECT().SaveWallet(w, req.Passphrase).Times(1).Return(nil)

	// when
	err = wallet.UntaintKey(store, req)

	// then
	require.NoError(t, err)
	assert.False(t, w.ListKeyPairs()[0].IsTainted())
}

func testUntaintingKeyOfNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.UntaintKeyRequest{
		Wallet:     vgrand.RandomStr(5),
		PubKey:     vgrand.RandomStr(25),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(0)

	// when
	err := wallet.UntaintKey(store, req)

	// then
	require.Error(t, err)
}

func TestIsolateKey(t *testing.T) {
	t.Run("Isolating key succeeds", testIsolatingKeySucceeds)
	t.Run("Isolating key of non-existing wallet fails", testIsolatingKeyOfNonExistingWalletFails)
}

func testIsolatingKeySucceeds(t *testing.T) {
	// given
	w := newWalletWithKey(t)
	kp := w.ListKeyPairs()[0]
	expectedResp := &wallet.IsolateKeyResponse{
		Wallet:   fmt.Sprintf("%s.%s.isolated", w.Name(), kp.PublicKey()[0:8]),
		FilePath: vgrand.RandomStr(10),
	}
	req := &wallet.IsolateKeyRequest{
		Wallet:     w.Name(),
		PubKey:     kp.PublicKey(),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(1).Return(nil)
	store.EXPECT().GetWalletPath(gomock.Any()).Times(1).Return(expectedResp.FilePath)

	// when
	resp, err := wallet.IsolateKey(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedResp, resp)
}

func testIsolatingKeyOfNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.IsolateKeyRequest{
		Wallet:     vgrand.RandomStr(5),
		PubKey:     vgrand.RandomStr(25),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(0)

	// when
	resp, err := wallet.IsolateKey(store, req)

	// then
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestListKeys(t *testing.T) {
	t.Run("List keys succeeds", testListKeysSucceeds)
	t.Run("List keys of non-existing wallet fails", testListKeysOfNonExistingWalletFails)
}

func testListKeysSucceeds(t *testing.T) {
	// given
	w := newWallet(t)
	keyCount := 3
	expectedKeys := &wallet.ListKeysResponse{
		Keys: make([]wallet.NamedPubKey, 0, keyCount),
	}
	for i := 0; i < keyCount; i++ {
		keyName := vgrand.RandomStr(5)
		kp, err := w.GenerateKeyPair([]wallet.Meta{{Key: "name", Value: keyName}})
		if err != nil {
			t.Fatalf("couldn't generate key: %v", err)
		}
		expectedKeys.Keys = append(expectedKeys.Keys, wallet.NamedPubKey{
			Name:      keyName,
			PublicKey: kp.PublicKey(),
		})
	}

	req := &wallet.ListKeysRequest{
		Wallet:     w.Name(),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)

	// when
	resp, err := wallet.ListKeys(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedKeys, resp)
}

func testListKeysOfNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.ListKeysRequest{
		Wallet:     vgrand.RandomStr(5),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)

	// when
	resp, err := wallet.ListKeys(store, req)

	// then
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetWalletInfo(t *testing.T) {
	t.Run("Get wallet info succeeds", testGetWalletInfoSucceeds)
	t.Run("Get wallet info of non-existing wallet fails", testGetWalletInfoOfNonExistingWalletFails)
}

func testGetWalletInfoSucceeds(t *testing.T) {
	// given
	w := newWallet(t)
	expectedKeys := &wallet.GetWalletInfoResponse{
		Type:    w.Type(),
		Version: w.Version(),
		ID:      w.ID(),
	}

	req := &wallet.GetWalletInfoRequest{
		Wallet:     w.Name(),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)

	// when
	resp, err := wallet.GetWalletInfo(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedKeys, resp)
}

func testGetWalletInfoOfNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.GetWalletInfoRequest{
		Wallet:     vgrand.RandomStr(5),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)

	// when
	resp, err := wallet.GetWalletInfo(store, req)

	// then
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestImportWalletSucceeds(t *testing.T) {
	// given
	walletName := vgrand.RandomStr(5)
	mnemonic := "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"

	req := &wallet.ImportWalletRequest{
		Wallet:     walletName,
		Mnemonic:   mnemonic,
		Version:    2,
		Passphrase: "passphrase",
	}

	expectedResp := &wallet.ImportWalletResponse{
		Name:     walletName,
		FilePath: vgrand.RandomStr(5),
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().SaveWallet(gomock.Any(), req.Passphrase).Times(1).Return(nil)
	store.EXPECT().GetWalletPath(req.Wallet).Times(1).Return(expectedResp.FilePath)

	// when
	resp, err := wallet.ImportWallet(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedResp, resp)
}

func TestListWalletsSucceeds(t *testing.T) {
	// given
	w1 := newWallet(t)
	w2 := newWallet(t)
	w3 := newWallet(t)
	walletNames := []string{w1.Name(), w2.Name(), w3.Name()}
	sort.Strings(walletNames)
	expectedResp := &wallet.ListWalletsResponse{
		Wallets: walletNames,
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().ListWallets().Times(1).Return(walletNames, nil)

	// when
	resp, err := wallet.ListWallets(store)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedResp, resp)
}

func TestSignMessage(t *testing.T) {
	t.Run("Sign message succeeds", testSignMessageSucceeds)
	t.Run("Sign message of non-existing wallet fails", testSignMessageWithNonExistingWalletFails)
}

func testSignMessageSucceeds(t *testing.T) {
	// given
	w := importWalletWithKey(t)
	kp := w.ListKeyPairs()[0]

	expectedKeys := &wallet.SignMessageResponse{
		Base64: "StH82RHxjQ3yTeaSN25b6sJwAyLiq1CDvPWf0X4KIf/WTIjkunkWKn1Gq9ntCoGBfBZIyNfpPtGx0TSZsSrbCA==",
		Bytes:  []byte{0x4a, 0xd1, 0xfc, 0xd9, 0x11, 0xf1, 0x8d, 0xd, 0xf2, 0x4d, 0xe6, 0x92, 0x37, 0x6e, 0x5b, 0xea, 0xc2, 0x70, 0x3, 0x22, 0xe2, 0xab, 0x50, 0x83, 0xbc, 0xf5, 0x9f, 0xd1, 0x7e, 0xa, 0x21, 0xff, 0xd6, 0x4c, 0x88, 0xe4, 0xba, 0x79, 0x16, 0x2a, 0x7d, 0x46, 0xab, 0xd9, 0xed, 0xa, 0x81, 0x81, 0x7c, 0x16, 0x48, 0xc8, 0xd7, 0xe9, 0x3e, 0xd1, 0xb1, 0xd1, 0x34, 0x99, 0xb1, 0x2a, 0xdb, 0x8},
	}

	req := &wallet.SignMessageRequest{
		Wallet:     w.Name(),
		PubKey:     kp.PublicKey(),
		Message:    []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit."),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)

	// when
	resp, err := wallet.SignMessage(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, expectedKeys, resp)
}

func testSignMessageWithNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.SignMessageRequest{
		Wallet:     vgrand.RandomStr(5),
		Passphrase: "passphrase",
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)

	// when
	resp, err := wallet.SignMessage(store, req)

	// then
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestRotateKey(t *testing.T) {
	t.Run("Rotate key succeeds", testRotateKeySucceeds)
	t.Run("Rotate key with non existing wallet fails", testRotateWithNonExistingWalletFails)
	t.Run("Rotate key with non existing public key fails", testRotateKeyWithNonExistingPublicKeyFails)
}

func testRotateKeySucceeds(t *testing.T) {
	// given
	w := importWalletWithKey(t)
	kp := w.ListKeyPairs()[0]

	masterKeyPair, err := w.GetMasterKeyPair()
	require.NoError(t, err)

	req := &wallet.RotateKeyRequest{
		Wallet:            w.Name(),
		Passphrase:        "passphrase",
		NewPublicKey:      kp.PublicKey(),
		TXBlockHeight:     20,
		TargetBlockHeight: 25,
	}

	expectedNewPubHash := hex.EncodeToString(vgcrypto.Hash([]byte(req.NewPublicKey)))

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)

	// when
	resp, err := wallet.RotateKey(store, req)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, masterKeyPair.PublicKey(), resp.MasterPublicKey)
	require.Equal(t, kp.PublicKey(), resp.NewPublicKey)

	transactionRaw, err := base64.StdEncoding.DecodeString(resp.Base64Transaction)
	require.NoError(t, err)

	transaction := &commandspb.Transaction{}
	err = proto.Unmarshal(transactionRaw, transaction)
	require.NoError(t, err)

	inputData := &commandspb.InputData{}
	err = proto.Unmarshal(transaction.InputData, inputData)
	require.NoError(t, err)

	keyRotate, ok := inputData.Command.(*commandspb.InputData_KeyRotateSubmission)
	require.True(t, ok)
	require.NotNil(t, keyRotate)

	require.Equal(t, uint64(req.TXBlockHeight), inputData.BlockHeight)
	require.Equal(t, uint64(kp.Index()), keyRotate.KeyRotateSubmission.KeyNumber)
	require.Equal(t, uint64(req.TargetBlockHeight), keyRotate.KeyRotateSubmission.TargetBlock)
	require.Equal(t, expectedNewPubHash, keyRotate.KeyRotateSubmission.NewPubKeyHash)
}

func testRotateWithNonExistingWalletFails(t *testing.T) {
	// given
	req := &wallet.RotateKeyRequest{
		Wallet:            vgrand.RandomStr(5),
		Passphrase:        "passphrase",
		NewPublicKey:      "nonexisting",
		TXBlockHeight:     20,
		TargetBlockHeight: 25,
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(false)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(0)

	// when
	resp, err := wallet.RotateKey(store, req)

	// then
	require.Error(t, err)
	assert.Nil(t, resp)
}

func testRotateKeyWithNonExistingPublicKeyFails(t *testing.T) {
	// given
	w := importWalletWithKey(t)

	req := &wallet.RotateKeyRequest{
		Wallet:            w.Name(),
		Passphrase:        "passphrase",
		NewPublicKey:      "nonexisting",
		TXBlockHeight:     20,
		TargetBlockHeight: 25,
	}

	// setup
	store := handlerMocks(t)
	store.EXPECT().WalletExists(req.Wallet).Times(1).Return(true)
	store.EXPECT().GetWallet(req.Wallet, req.Passphrase).Times(1).Return(w, nil)

	// when
	resp, err := wallet.RotateKey(store, req)

	// then
	require.Nil(t, resp)
	require.Error(t, err)
}

func newWalletWithKey(t *testing.T) *wallet.HDWallet {
	t.Helper()
	return newWalletWithKeys(t, 1)
}

func newWalletWithKeys(t *testing.T, n int) *wallet.HDWallet {
	t.Helper()
	w := newWallet(t)

	for i := 0; i < n; i++ {
		if _, err := w.GenerateKeyPair(nil); err != nil {
			t.Fatalf("couldn't generate key: %v", err)
		}
	}
	return w
}

func importWalletWithKey(t *testing.T) *wallet.HDWallet {
	t.Helper()
	w, err := wallet.ImportHDWallet(
		vgrand.RandomStr(5),
		"swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render",
		2,
	)
	if err != nil {
		t.Fatalf("couldn't import wallet: %v", err)
	}

	if _, err := w.GenerateKeyPair(nil); err != nil {
		t.Fatalf("couldn't generate key: %v", err)
	}
	return w
}

func newWallet(t *testing.T) *wallet.HDWallet {
	t.Helper()
	w, _, err := wallet.NewHDWallet(vgrand.RandomStr(5))
	if err != nil {
		t.Fatalf("couldn't create HD wallet: %v", err)
	}
	return w
}

func handlerMocks(t *testing.T) *mocks.MockStore {
	t.Helper()
	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)
	return store
}
