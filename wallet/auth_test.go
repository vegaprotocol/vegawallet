package wallet_test

import (
	"os"
	"testing"
	"time"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testAuth struct {
	wallet.Auth
	rootPath string
}

func getTestAuth(t *testing.T) *testAuth {
	rootPath := rootDir()
	tokenExpiry := 10 * time.Hour
	fsutil.EnsureDir(rootPath)
	log := zap.NewNop()

	// gen keys
	err := wallet.GenRsaKeyFiles(log, rootPath, false)
	if err != nil {
		t.Fatal(err)
	}
	a, err := wallet.NewAuth(log, rootPath, tokenExpiry)
	if err != nil {
		t.Fatal(err)
	}

	return &testAuth{
		Auth:     a,
		rootPath: rootPath,
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
	wallet := "jeremy"

	// get a new session
	tok, err := auth.NewSession(wallet)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	wallet2, err := auth.VerifyToken(tok)
	assert.NoError(t, err)
	assert.Equal(t, wallet, wallet2)

	assert.NoError(t, os.RemoveAll(auth.rootPath))
}

func testVerifyInvalidToken(t *testing.T) {
	auth := getTestAuth(t)
	tok := "hehehe that's not a toekn"

	w, err := auth.VerifyToken(tok)
	assert.EqualError(t, err, "token is malformed: token contains an invalid number of segments")
	assert.Empty(t, w)

	assert.NoError(t, os.RemoveAll(auth.rootPath))
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
	assert.EqualError(t, err, wallet.ErrSessionNotFound.Error())
	assert.Empty(t, w)

	assert.NoError(t, os.RemoveAll(auth.rootPath))
}

func testRevokeInvalidToken(t *testing.T) {
	auth := getTestAuth(t)
	tok := "hehehe that's not a toekn"

	err := auth.Revoke(tok)
	assert.EqualError(t, err, "token is malformed: token contains an invalid number of segments")

	assert.NoError(t, os.RemoveAll(auth.rootPath))
}
