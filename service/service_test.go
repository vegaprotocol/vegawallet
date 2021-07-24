package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/internal/proto/api"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/mocks"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// this tests in general ensure request / response contracts are not broken for the service

const (
	TestMnemonic = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"
)

type testService struct {
	*service.Service

	ctrl        *gomock.Controller
	handler     *mocks.MockWalletHandler
	nodeForward *mocks.MockNodeForward
	auth        *mocks.MockAuth
}

func getTestService(t *testing.T) *testService {
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockWalletHandler(ctrl)
	auth := mocks.NewMockAuth(ctrl)
	nodeForward := mocks.NewMockNodeForward(ctrl)
	// no needs of the conf or path as we do not run an actual service
	s, _ := service.NewServiceWith(zap.NewNop(), nil, handler, auth, nodeForward, "v1.2.3", "abcdef12")
	return &testService{
		Service:     s,
		ctrl:        ctrl,
		handler:     handler,
		auth:        auth,
		nodeForward: nodeForward,
	}
}

func TestService(t *testing.T) {
	t.Run("create wallet ok", testServiceCreateWalletOK)
	t.Run("create wallet fail invalid request", testServiceCreateWalletFailInvalidRequest)
	t.Run("Importing a wallet succeeds", testServiceImportWalletOK)
	t.Run("Importing a wallet with and invalid request fails", testServiceImportWalletFailInvalidRequest)
	t.Run("login wallet ok", testServiceLoginWalletOK)
	t.Run("download wallet ok", testServiceDownloadWalletOK)
	t.Run("login wallet fail invalid request", testServiceLoginWalletFailInvalidRequest)
	t.Run("revoke token ok", testServiceRevokeTokenOK)
	t.Run("revoke token fail invalid request", testServiceRevokeTokenFailInvalidRequest)
	t.Run("gen keypair ok", testServiceGenKeypairOK)
	t.Run("gen keypair fail invalid request", testServiceGenKeypairFailInvalidRequest)
	t.Run("list keypair ok", testServiceListPublicKeysOK)
	t.Run("list keypair fail invalid request", testServiceListPublicKeysFailInvalidRequest)
	t.Run("get keypair ok", testServiceGetPublicKeyOK)
	t.Run("get keypair fail invalid request", testServiceGetPublicKeyFailInvalidRequest)
	t.Run("get keypair fail key not found", testServiceGetPublicKeyFailKeyNotFound)
	t.Run("get keypair fail misc error", testServiceGetPublicKeyFailMiscError)
	t.Run("taint ok", testServiceTaintOK)
	t.Run("taint fail invalid request", testServiceTaintFailInvalidRequest)
	t.Run("update meta", testServiceUpdateMetaOK)
	t.Run("update meta invalid request", testServiceUpdateMetaFailInvalidRequest)
	t.Run("Signing transaction succeeds", testSigningTransactionSucceeds)
	t.Run("Signing transaction with propagation succeeds", testSigningTransactionWithPropagationSucceeds)
	t.Run("Signing transaction with failed propagation fails", testSigningTransactionWithFailedPropagationFails)
	t.Run("Failed signing of transaction fails", testFailedSigningTransactionFails)
	t.Run("Signing transaction with invalid payload fails", testSigningTransactionWithInvalidPayloadFails)
	t.Run("Signing transaction without pub-key fails", testSigningTransactionWithoutPubKeyFails)
	t.Run("Signing transaction without command fails", testSigningTransactionWithoutCommandFails)
	t.Run("Signing anything succeeds", testSigningAnythingSucceeds)
	t.Run("Verifying anything succeeds", testVerifyingAnythingSucceeds)
	t.Run("Failed verification fails", testVerifyingAnythingFails)
}

func testServiceCreateWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().CreateWallet("jeremy", "oh yea?").Times(1).Return(TestMnemonic, nil)
	s.auth.EXPECT().NewSession("jeremy").Times(1).Return("this is a token", nil)

	payload := `{"wallet": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.CreateWallet(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceCreateWalletFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	payload := `{"wall": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.CreateWallet(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"wallet": "jeremy", "passrase": "oh yea?"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()

	s.CreateWallet(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceImportWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().ImportWallet("jeremy", "oh yea?", TestMnemonic).Times(1).Return(nil)
	s.auth.EXPECT().NewSession("jeremy").Times(1).Return("this is a token", nil)

	payload := fmt.Sprintf(`{"wallet": "jeremy", "passphrase": "oh yea?", "mnemonic": "%s"}`, TestMnemonic)
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.ImportWallet(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceImportWalletFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	payload := fmt.Sprintf(`{"wall": "jeremy", "passphrase": "oh yea?", "mnemonic": \"%s\"}`, TestMnemonic)
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.CreateWallet(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = fmt.Sprintf(`{"wallet": "jeremy", "password": "oh yea?", "mnemonic": \"%s\"}`, TestMnemonic)
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()

	s.ImportWallet(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = fmt.Sprintf(`{"wallet": "jeremy", "passphrase": "oh yea?", "little_words": \"%s\"}`, TestMnemonic)
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()

	s.CreateWallet(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceLoginWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	s.auth.EXPECT().NewSession("jeremy").Times(1).Return("this is a token", nil)

	payload := `{"wallet": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.Login(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceDownloadWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	payload := `{"wallet": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.handler.EXPECT().LoginWallet("jeremy", "oh yea?").Times(1).Return(nil)
	s.auth.EXPECT().NewSession("jeremy").Times(1).Return("this is a token", nil)

	s.Login(w, r, nil)

	resp := w.Result()
	var token struct {
		Token string
	}
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	raw, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	_ = json.Unmarshal(raw, &token)

	tmpFile, _ := ioutil.TempFile(".", "test-wallet")
	defer func() {
		name := tmpFile.Name()
		tmpFile.Close()
		os.Remove(name)
	}()

	s.auth.EXPECT().VerifyToken("this is a token").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().GetWalletPath("jeremy").Times(1).Return(tmpFile.Name(), nil)

	// now get the file:
	r = httptest.NewRequest(http.MethodGet, "scheme://host/path", bytes.NewBufferString(""))
	w = httptest.NewRecorder()

	s.DownloadWallet(token.Token, w, r, nil)
	resp = w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func testServiceLoginWalletFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	payload := `{"wall": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(0)
	s.auth.EXPECT().NewSession(gomock.Any()).Times(0)

	s.Login(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"wallet": "jeremy", "passrase": "oh yea?"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()

	s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(0)
	s.auth.EXPECT().NewSession(gomock.Any()).Times(0)

	s.Login(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceRevokeTokenOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.auth.EXPECT().Revoke(gomock.Any()).Times(1).Return(nil)
	s.handler.EXPECT().LogoutWallet("jeremy").Times(1)

	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.Revoke)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceRevokeTokenFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.Revoke)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.Revoke)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGenKeypairOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	key := &wallet.LegacyPublicKey{
		Pub:       "0xdeadbeef",
		Algorithm: crypto.NewEd25519(),
		Tainted:   false,
		MetaList:  nil,
	}

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().SecureGenerateKeyPair("jeremy", "oh yea?", gomock.Len(0)).Times(1).Return("0xdeadbeef", nil)
	s.handler.EXPECT().GetPublicKey("jeremy", "0xdeadbeef").Times(1).Return(key, nil)

	payload := `{"passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GenerateKeyPair)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceGenKeypairFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GenerateKeyPair)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.GenerateKeyPair)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but no payload
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.GenerateKeyPair)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func testServiceListPublicKeysOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().ListPublicKeys("jeremy").Times(1).
		Return([]wallet.PublicKey{}, nil)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.ListPublicKeys)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceListPublicKeysFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.ListPublicKeys)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.ListPublicKeys)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGetPublicKeyOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	kp := wallet.LegacyKeyPair{
		Pub:       "pub",
		Priv:      "",
		Algorithm: crypto.NewEd25519(),
		Tainted:   false,
		MetaList:  []wallet.Meta{{Key: "a", Value: "b"}},
	}
	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(kp.ToPublicKey(), nil)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceGetPublicKeyFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GetPublicKey)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.GetPublicKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGetPublicKeyFailKeyNotFound(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, wallet.ErrPubKeyDoesNotExist)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func testServiceGetPublicKeyFailMiscError(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("an error"))

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

	resp := w.Result()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func testServiceTaintOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().TaintKey(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).Return(nil)
	payload := `{"passphrase": "some data"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.TaintKey)(w, r, httprouter.Params{{Key: "keyid", Value: "asdasasdasd"}})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceTaintFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Set("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.TaintKey)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"passhp": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"passphrase": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceUpdateMetaOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().UpdateMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).Return(nil)
	payload := `{"passphrase": "some data", "meta": [{"key":"ok", "value":"primary"}]}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.UpdateMeta)(w, r, httprouter.Params{{Key: "keyid", Value: "asdasasdasd"}})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceUpdateMetaFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Set("Authorization", "Bearer")

	w := httptest.NewRecorder()

	service.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"passhp": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"passphrase": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testSigningTransactionSucceeds(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "eyXXzA"
	name := "jeremy"
	payload := `{"pubKey": "0xCAFEDUDE", "orderCancellation": {}}`
	request := newAuthenticatedRequest(payload)
	response := httptest.NewRecorder()

	// setup
	s.auth.EXPECT().
		VerifyToken(token).
		Times(1).
		Return(name, nil)
	s.handler.EXPECT().
		SignTxV2(name, gomock.Any(), gomock.Any()).
		Times(1).
		Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().
		SendTxV2(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionV2Request_TYPE_ASYNC).
		Times(0)
	s.nodeForward.EXPECT().LastBlockHeight(gomock.Any()).
		Times(1).Return(uint64(42), nil)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func testSigningTransactionWithPropagationSucceeds(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "eyXXzA"
	name := "jeremy"
	payload := `{"propagate": true, "pubKey": "0xCAFEDUDE", "orderCancellation": {}}`
	request := newAuthenticatedRequest(payload)
	response := httptest.NewRecorder()

	// setup
	s.auth.EXPECT().
		VerifyToken(token).
		Times(1).
		Return(name, nil)
	s.handler.EXPECT().
		SignTxV2(name, gomock.Any(), gomock.Any()).
		Times(1).
		Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().
		SendTxV2(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionV2Request_TYPE_SYNC).
		Times(1).
		Return(nil)
	s.nodeForward.EXPECT().LastBlockHeight(gomock.Any()).
		Times(1).Return(uint64(42), nil)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func testSigningTransactionWithFailedPropagationFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "eyXXzA"
	name := "jeremy"
	payload := `{"propagate": true, "pubKey": "0xCAFEDUDE", "orderCancellation": {}}`
	request := newAuthenticatedRequest(payload)
	response := httptest.NewRecorder()

	// setup
	s.auth.EXPECT().
		VerifyToken(token).
		Times(1).
		Return(name, nil)
	s.handler.EXPECT().
		SignTxV2(name, gomock.Any(), gomock.Any()).
		Times(1).
		Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().
		SendTxV2(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionV2Request_TYPE_SYNC).
		Times(1).
		Return(errors.New("failure"))
	s.nodeForward.EXPECT().LastBlockHeight(gomock.Any()).
		Times(1).Return(uint64(42), nil)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func testFailedSigningTransactionFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "eyXXzA"
	name := "jeremy"
	payload := `{"pubKey": "0xCAFEDUDE", "orderCancellation": {}}`
	request := newAuthenticatedRequest(payload)
	response := httptest.NewRecorder()

	// setup
	s.auth.EXPECT().
		VerifyToken(token).
		Times(1).
		Return(name, nil)
	s.handler.EXPECT().
		SignTxV2(name, gomock.Any(), gomock.Any()).
		Times(1).
		Return(nil, errors.New("failure"))
	s.nodeForward.EXPECT().LastBlockHeight(gomock.Any()).
		Times(1).Return(uint64(42), nil)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	assert.Equal(t, http.StatusForbidden, result.StatusCode)
}

func testSigningTransactionWithInvalidPayloadFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "eyXXzA"
	payload := `{"badKey": "0xCAFEDUDE"}`
	request := newAuthenticatedRequest(payload)
	response := httptest.NewRecorder()

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func testSigningTransactionWithoutPubKeyFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "0xDEADBEEF"
	payload := `{"orderSubmission": {}}`
	response := httptest.NewRecorder()
	request := newAuthenticatedRequest(payload)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	require.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func testSigningTransactionWithoutCommandFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// given
	token := "0xDEADBEEF"
	payload := `{"pubKey": "0xCAFEDUDE"}`
	response := httptest.NewRecorder()
	request := newAuthenticatedRequest(payload)

	// when
	s.SignTxSyncV2(token, response, request, nil)

	// then
	result := response.Result()
	require.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func testSigningAnythingSucceeds(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().SignAny("jeremy", []byte("spice of dune"), "asdasasdasd").
		Times(1).Return([]byte("some sig"), nil)
	payload := `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.SignAny)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testVerifyingAnythingSucceeds(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().VerifyAny("jeremy", []byte("spice of dune"), []byte("Sietch Tabr"), "asdasasdasd").
		Times(1).Return(true, nil)
	payload := `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd", "signature": "U2lldGNoIFRhYnI="}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.VerifyAny)(w, r, nil)

	httpResponse := w.Result()
	resp := service.SuccessResponse{}
	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
	unmarshalResponse(httpResponse, &resp)
	assert.True(t, resp.Success)
}

func testVerifyingAnythingFails(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().VerifyAny("jeremy", []byte("spice of dune"), []byte("Sietch Tabr"), "asdasasdasd").
		Times(1).Return(false, nil)
	payload := `{"inputData":"c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd", "signature": "U2lldGNoIFRhYnI="}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.VerifyAny)(w, r, nil)

	httpResponse := w.Result()
	resp := service.SuccessResponse{}
	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
	unmarshalResponse(httpResponse, &resp)
	assert.False(t, resp.Success)
}

func newAuthenticatedRequest(payload string) *http.Request {
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")
	return r
}

func unmarshalResponse(r *http.Response, into interface{}) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, into)
	if err != nil {
		panic(err)
	}
}
