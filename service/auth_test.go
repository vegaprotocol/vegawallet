package service_test

import (
	"testing"
	"time"

	"code.vegaprotocol.io/vegawallet/service"
	"code.vegaprotocol.io/vegawallet/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testAuth struct {
	service.Auth
	ctrl *gomock.Controller
}

func getTestAuth(t *testing.T) *testAuth {
	rsaKeys, err := service.GenerateRSAKeys()
	if err != nil {
		t.Fatal(err)
	}

	ctrl := gomock.NewController(t)
	store := mocks.NewMockRSAStore(ctrl)
	store.EXPECT().GetRsaKeys().Return(rsaKeys, nil)

	tokenExpiry := 10 * time.Hour
	a, err := service.NewAuth(zap.NewNop(), store, tokenExpiry)
	if err != nil {
		t.Fatal(err)
	}

	return &testAuth{
		Auth: a,
		ctrl: ctrl,
	}
}

func TestAuth(t *testing.T) {
	t.Run("verify a valid token", testVerifyValidToken)
	t.Run("verify an invalid token fail", testVerifyInvalidToken)
	t.Run("revoke a valid token", testRevokeValidToken)
	t.Run("revoke an invalid token fail", testRevokeInvalidToken)
}

func testVerifyValidToken(t *testing.T) {
	auth := getTestAuth(t)
	w := "jeremy"

	// get a new session
	tok, err := auth.NewSession(w)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	wallet2, err := auth.VerifyToken(tok)
	assert.NoError(t, err)
	assert.Equal(t, w, wallet2)
}

func testVerifyInvalidToken(t *testing.T) {
	auth := getTestAuth(t)
	tok := "that's not a token"

	w, err := auth.VerifyToken(tok)
	assert.EqualError(t, err, "token is malformed: token contains an invalid number of segments")
	assert.Empty(t, w)
}

func testRevokeValidToken(t *testing.T) {
	auth := getTestAuth(t)
	walletname := "jeremy"

	// get a new session
	tok, err := auth.NewSession(walletname)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	wallet2, err := auth.VerifyToken(tok)
	assert.NoError(t, err)
	assert.Equal(t, walletname, wallet2)

	// now we made sure the token exists, let's revoke and re-verify it
	err = auth.Revoke(tok)
	assert.NoError(t, err)

	w, err := auth.VerifyToken(tok)
	assert.EqualError(t, err, service.ErrSessionNotFound.Error())
	assert.Empty(t, w)
}

func testRevokeInvalidToken(t *testing.T) {
	auth := getTestAuth(t)
	tok := "hehehe that's not a toekn"

	err := auth.Revoke(tok)
	assert.EqualError(t, err, "token is malformed: token contains an invalid number of segments")
}
