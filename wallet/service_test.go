package wallet_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/go-wallet/mocks"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// this tests in general ensure request / response contracts are not broken for the service

type testService struct {
	*wallet.Service

	ctrl        *gomock.Controller
	handler     *mocks.MockWalletHandler
	nodeForward *mocks.MockNodeForward
}

func getTestService(t *testing.T) *testService {
	ctrl := gomock.NewController(t)
	handler := mocks.NewMockWalletHandler(ctrl)
	nodeForward := mocks.NewMockNodeForward(ctrl)
	// no needs of the conf or path as we do not run an actual service
	s, _ := wallet.NewServiceWith(logging.NewTestLogger(), nil, "", handler, nodeForward)
	return &testService{
		Service:     s,
		ctrl:        ctrl,
		handler:     handler,
		nodeForward: nodeForward,
	}
}

func TestService(t *testing.T) {
	t.Run("create wallet ok", testServiceCreateWalletOK)
	t.Run("create wallet fail invalid request", testServiceCreateWalletFailInvalidRequest)
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
	t.Run("sign ok", testServiceSignOK)
	t.Run("sign fail invalid request", testServiceSignFailInvalidRequest)
	t.Run("taint ok", testServiceTaintOK)
	t.Run("taint fail invalid request", testServiceTaintFailInvalidRequest)
	t.Run("update meta", testServiceUpdateMetaOK)
	t.Run("update meta invalid request", testServiceUpdateMetaFailInvalidRequest)
}

func testServiceCreateWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().CreateWallet(gomock.Any(), gomock.Any()).Times(1).Return("this is a token", nil)

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

func testServiceLoginWalletOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(1).Return("this is a token", nil)

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

	s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(1).Return("this is a token", nil)

	payload := `{"wallet": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.Login(w, r, nil)

	resp := w.Result()
	var token struct {
		Data string
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
	s.handler.EXPECT().WalletPath(token.Data).Times(1).Return(tmpFile.Name(), nil)

	// now get the file:
	r = httptest.NewRequest(http.MethodGet, "scheme://host/path", bytes.NewBufferString(""))
	w = httptest.NewRecorder()

	s.DownloadWallet(token.Data, w, r, nil)
	resp = w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func testServiceLoginWalletFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	payload := `{"wall": "jeremy", "passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	s.Login(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"wallet": "jeremy", "passrase": "oh yea?"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()

	s.Login(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceRevokeTokenOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().RevokeToken(gomock.Any()).Times(1).Return(nil)

	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.Revoke)(w, r, nil)

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

	wallet.ExtractToken(s.Revoke)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.Revoke)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGenKeypairOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().GenerateKeypair(gomock.Any(), gomock.Any()).Times(1).Return("", nil)
	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).Return(&wallet.Keypair{}, nil)

	payload := `{"passphrase": "oh yea?"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.GenerateKeypair)(w, r, nil)

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

	wallet.ExtractToken(s.GenerateKeypair)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.GenerateKeypair)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but no payload
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.GenerateKeypair)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func testServiceListPublicKeysOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().ListPublicKeys(gomock.Any()).Times(1).
		Return([]wallet.Keypair{}, nil)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.ListPublicKeys)(w, r, nil)

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

	wallet.ExtractToken(s.ListPublicKeys)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.ListPublicKeys)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGetPublicKeyOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	kp := wallet.Keypair{
		Pub:       "pub",
		Priv:      "",
		Algorithm: crypto.NewEd25519(),
		Tainted:   false,
		Meta:      []wallet.Meta{{Key: "a", Value: "b"}},
	}
	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(&kp, nil)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

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

	wallet.ExtractToken(s.GetPublicKey)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.GetPublicKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func testServiceGetPublicKeyFailKeyNotFound(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, wallet.ErrPubKeyDoesNotExist)

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func testServiceGetPublicKeyFailMiscError(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Times(1).
		Return(nil, errors.New("an error"))

	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.GetPublicKey)(w, r, httprouter.Params{{Key: "keyid", Value: "apubkey"}})

	resp := w.Result()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func testServiceSignOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().SignTx(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).Return(wallet.SignedBundle{}, nil)
	payload := `{"tx": "some data", "pubKey": "asdasasdasd"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceSignFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// InvalidMethod
	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	w := httptest.NewRecorder()

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// invalid token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w = httptest.NewRecorder()

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"t": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"tx": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func testServiceTaintOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().TaintKey(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).Return(nil)
	payload := `{"passphrase": "some data"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.TaintKey)(w, r, httprouter.Params{{Key: "keyid", Value: "asdasasdasd"}})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceTaintFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.TaintKey)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"passhp": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"passphrase": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.TaintKey)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func testServiceUpdateMetaOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.handler.EXPECT().UpdateMeta(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).Return(nil)
	payload := `{"passphrase": "some data", "meta": [{"key":"ok", "value":"primary"}]}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Add("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.UpdateMeta)(w, r, httprouter.Params{{Key: "keyid", Value: "asdasasdasd"}})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceUpdateMetaFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// invalid token
	r := httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Add("Authorization", "Bearer")

	w := httptest.NewRecorder()

	wallet.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	wallet.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"passhp": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"passphrase": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Add("Authorization", "Bearer eyXXzA")

	wallet.ExtractToken(s.UpdateMeta)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}
