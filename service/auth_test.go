package service_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	rootDirPath = "/tmp/vegatests/wallet/"
)

type testAuth struct {
	service.Auth
	rootPath string
}

func getTestAuth(t *testing.T) *testAuth {
	rootPath := rootDir()
	tokenExpiry := 10 * time.Hour
	err := fsutil.EnsureDir(rootPath)
	if err != nil {
		panic(err)
	}
	log := zap.NewNop()

	// gen keys
	err = wallet.GenRsaKeyFiles(log, rootPath, false)
	if err != nil {
		t.Fatal(err)
	}
	a, err := service.NewAuth(log, rootPath, tokenExpiry)
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
	w := "jeremy"

	// get a new session
	tok, err := auth.NewSession(w)
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)

	wallet2, err := auth.VerifyToken(tok)
	assert.NoError(t, err)
	assert.Equal(t, w, wallet2)

	assert.NoError(t, os.RemoveAll(auth.rootPath))
}

func testVerifyInvalidToken(t *testing.T) {
	auth := getTestAuth(t)
	tok := "that's not a token"

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
	assert.EqualError(t, err, service.ErrSessionNotFound.Error())
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

func rootDir() string {
	path := filepath.Join(rootDirPath, crypto.RandomStr(10))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return path
}
