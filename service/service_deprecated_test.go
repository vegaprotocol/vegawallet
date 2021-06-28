package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeprecatedService(t *testing.T) {
	t.Run("sign ok", testServiceSignOK)
	t.Run("sign fail invalid request", testServiceSignFailInvalidRequest)
}

func testServiceSignOK(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	s.auth.EXPECT().VerifyToken("eyXXzA").Times(1).Return("jeremy", nil)
	s.handler.EXPECT().SignTx("jeremy", gomock.Any(), "asdasasdasd", uint64(42)).
		Times(1).Return(wallet.SignedBundle{}, nil)
	s.nodeForward.EXPECT().LastBlockHeight(gomock.Any()).
		Times(1).Return(uint64(42), nil)
	payload := `{"tx": "some data", "pubKey": "asdasasdasd"}`
	r := httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	r.Header.Set("Authorization", "Bearer eyXXzA")

	w := httptest.NewRecorder()

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func testServiceSignFailInvalidRequest(t *testing.T) {
	s := getTestService(t)
	defer s.ctrl.Finish()

	// InvalidMethod
	r := httptest.NewRequest("GET", "scheme://host/path", nil)
	w := httptest.NewRecorder()

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// invalid token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	r.Header.Set("Authorization", "Bearer")

	w = httptest.NewRecorder()

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// no token
	r = httptest.NewRequest("POST", "scheme://host/path", nil)
	w = httptest.NewRecorder()

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// token but invalid payload
	payload := `{"t": "some data", "pubKey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	payload = `{"tx": "some data", "puey": "asdasasdasd"}`
	r = httptest.NewRequest("POST", "scheme://host/path", bytes.NewBufferString(payload))
	w = httptest.NewRecorder()
	r.Header.Set("Authorization", "Bearer eyXXzA")

	service.ExtractToken(s.SignTx)(w, r, nil)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}
