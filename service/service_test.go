package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	api "code.vegaprotocol.io/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/crypto"
	"code.vegaprotocol.io/vegawallet/network"
	"code.vegaprotocol.io/vegawallet/service"
	"code.vegaprotocol.io/vegawallet/service/mocks"
	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testRecoveryPhrase = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"

	testRequestTimeout = 10 * time.Second
)

type testService struct {
	*service.Service

	ctrl                 *gomock.Controller
	handler              *mocks.MockWalletHandler
	nodeForward          *mocks.MockNodeForward
	auth                 *mocks.MockAuth
	ConsentConfirmations chan service.ConsentConfirmation
}

func getTestService(t *testing.T, consentPolicy string) *testService {
	t.Helper()

	ctrl := gomock.NewController(t)
	handler := mocks.NewMockWalletHandler(ctrl)
	auth := mocks.NewMockAuth(ctrl)
	nodeForward := mocks.NewMockNodeForward(ctrl)

	pendingConsents := make(chan service.ConsentRequest, 1)
	consentConfirmations := make(chan service.ConsentConfirmation, 1)

	var policy service.Policy
	switch consentPolicy {
	case "automatic":
		policy = service.NewAutomaticConsentPolicy(pendingConsents, consentConfirmations)
	case "manual":
		policy = service.NewExplicitConsentPolicy(pendingConsents, consentConfirmations)
	default:
		t.Fatalf("unknown consent policy: %s", consentPolicy)
	}
	// no needs of the conf or path as we do not run an actual service
	s, err := service.NewService(zap.NewNop(), &network.Network{}, handler, auth, nodeForward, policy)
	if err != nil {
		t.Fatalf("couldn't create service: %v", err)
	}

	return &testService{
		Service:              s,
		ctrl:                 ctrl,
		handler:              handler,
		auth:                 auth,
		nodeForward:          nodeForward,
		ConsentConfirmations: consentConfirmations,
	}
}

func TestService(t *testing.T) {
	t.Run("create wallet ok", testServiceCreateWalletOK)
	t.Run("create wallet fail invalid request", testServiceCreateWalletFailInvalidRequest)
	t.Run("Importing a wallet succeeds", testServiceImportWalletOK)
	t.Run("Importing a wallet with and invalid request fails", testServiceImportWalletFailInvalidRequest)
	t.Run("login wallet ok", testServiceLoginWalletOK)
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
	t.Run("update metadata", testServiceUpdateMetaOK)
	t.Run("update metadata invalid request", testServiceUpdateMetaFailInvalidRequest)
	t.Run("Signing transaction automatically succeeds", testSigningTransactionSucceeds)
	t.Run("Signing transaction manually succeeds", testAcceptSigningTransactionManuallySucceeds)
	t.Run("Decline signing transaction manually succeeds", testDeclineSigningTransactionManuallySucceeds)
	t.Run("Signing transaction with propagation succeeds", testSigningTransactionWithPropagationSucceeds)
	t.Run("Signing transaction with failed propagation fails", testSigningTransactionWithFailedPropagationFails)
	t.Run("Failed signing of transaction fails", testFailedTransactionSigningFails)
	t.Run("Signing transaction with invalid request fails", testSigningTransactionWithInvalidRequestFails)
	t.Run("Signing anything succeeds", testSigningAnythingSucceeds)
	t.Run("Signing anything with invalid request fails", testSigningAnyDataWithInvalidRequestFails)
	t.Run("Verifying anything succeeds", testVerifyingAnythingSucceeds)
	t.Run("Failed verification fails", testVerifyingAnythingFails)
	t.Run("Verifying anything with invalid request fails", testVerifyingAnyDataWithInvalidRequestFails)
}

func testServiceCreateWalletOK(t *testing.T) {
	s := getTestService(t, "automatic")
	t.Cleanup(func() {
		s.ctrl.Finish()
	})

	// given
	walletName := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"wallet": "%s", "passphrase": "%s"}`, walletName, passphrase)

	// setup
	s.handler.EXPECT().CreateWallet(walletName, passphrase).Times(1).Return(testRecoveryPhrase, nil)
	s.auth.EXPECT().NewSession(walletName).Times(1).Return("this is a token", nil)

	// when
	statusCode, _ := serveHTTP(t, s, createWalletRequest(t, payload))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceCreateWalletFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		payload string
	}{
		{
			name:    "misspelled wallet property",
			payload: `{"wall": "jeremy", "passphrase": "oh yea?"}`,
		}, {
			name:    "misspelled passphrase property",
			payload: `{"wallet": "jeremy", "passrase": "oh yea?"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, createWalletRequest(tt, tc.payload))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceImportWalletOK(t *testing.T) {
	tcs := []struct {
		name    string
		version uint32
	}{
		{
			name:    "version 1",
			version: 1,
		}, {
			name:    "version 2",
			version: 2,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// given
			walletName := vgrand.RandomStr(5)
			passphrase := vgrand.RandomStr(5)
			payload := fmt.Sprintf(`{"wallet": "%s", "passphrase": "%s", "recoveryPhrase": "%s", "version": %d}`, walletName, passphrase, testRecoveryPhrase, tc.version)

			// setup
			s.handler.EXPECT().ImportWallet(walletName, passphrase, testRecoveryPhrase, tc.version).Times(1).Return(nil)
			s.auth.EXPECT().NewSession(walletName).Times(1).Return("this is a token", nil)

			// when
			statusCode, _ := serveHTTP(tt, s, importWalletRequest(tt, payload))

			// then
			assert.Equal(tt, http.StatusOK, statusCode)
		})
	}
}

func testServiceImportWalletFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		payload string
	}{
		{
			name:    "misspelled wallet property",
			payload: fmt.Sprintf(`{"wall": "jeremy", "passphrase": "oh yea?", "recoveryPhrase": \"%s\"}`, testRecoveryPhrase),
		}, {
			name:    "misspelled passphrase property",
			payload: fmt.Sprintf(`{"wallet": "jeremy", "password": "oh yea?", "recoveryPhrase": \"%s\"}`, testRecoveryPhrase),
		}, {
			name:    "misspelled recovery phrase property",
			payload: fmt.Sprintf(`{"wallet": "jeremy", "passphrase": "oh yea?", "little_words": \"%s\"}`, testRecoveryPhrase),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, importWalletRequest(tt, tc.payload))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceLoginWalletOK(t *testing.T) {
	s := getTestService(t, "automatic")
	t.Cleanup(func() {
		s.ctrl.Finish()
	})

	// given
	walletName := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"wallet": "%s", "passphrase": "%s"}`, walletName, passphrase)

	// setup
	s.handler.EXPECT().LoginWallet(walletName, passphrase).Times(1).Return(nil)
	s.auth.EXPECT().NewSession(walletName).Times(1).Return("this is a token", nil)

	// when
	statusCode, _ := serveHTTP(t, s, loginRequest(t, payload))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceLoginWalletFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		payload string
	}{
		{
			name:    "misspelled wallet property",
			payload: `{"wall": "jeremy", "passphrase": "oh yea?"}`,
		}, {
			name:    "misspelled passphrase property",
			payload: `{"wallet": "jeremy", "passrase": "oh yea?"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			t.Cleanup(func() {
				s.ctrl.Finish()
			})

			// setup
			s.handler.EXPECT().LoginWallet(gomock.Any(), gomock.Any()).Times(0)
			s.auth.EXPECT().NewSession(gomock.Any()).Times(0)

			// when
			statusCode, _ := serveHTTP(tt, s, loginRequest(tt, tc.payload))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceRevokeTokenOK(t *testing.T) {
	s := getTestService(t, "automatic")
	t.Cleanup(func() {
		s.ctrl.Finish()
	})

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)

	// setup
	s.auth.EXPECT().Revoke(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().LogoutWallet(walletName).Times(1)

	// when
	statusCode, _ := serveHTTP(t, s, logoutRequest(t, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceRevokeTokenFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(t, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, logoutRequest(t, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceGenKeypairOK(t *testing.T) {
	s := getTestService(t, "automatic")
	t.Cleanup(func() {
		s.ctrl.Finish()
	})

	// given
	ed25519 := crypto.NewEd25519()
	key := &wallet.HDPublicKey{
		PublicKey: vgrand.RandomStr(5),
		Algorithm: wallet.Algorithm{
			Name:    ed25519.Name(),
			Version: ed25519.Version(),
		},
		Tainted:  false,
		MetaList: nil,
	}
	walletName := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"passphrase": "%s"}`, passphrase)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SecureGenerateKeyPair(walletName, passphrase, gomock.Len(0)).Times(1).Return(key.PublicKey, nil)
	s.handler.EXPECT().GetPublicKey(walletName, key.PublicKey).Times(1).Return(key, nil)

	// when
	statusCode, _ := serveHTTP(t, s, generateKeyRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceGenKeypairFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
		payload string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
			payload: `{"passphrase": "oh yea?"}`,
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
			payload: `{"passphrase": "oh yea?"}`,
		}, {
			name:    "invalid request",
			headers: authHeaders(t, vgrand.RandomStr(5)),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, generateKeyRequest(t, tc.payload, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceListPublicKeysOK(t *testing.T) {
	s := getTestService(t, "automatic")
	t.Cleanup(func() {
		s.ctrl.Finish()
	})

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().ListPublicKeys(walletName).Times(1).Return([]wallet.PublicKey{}, nil)

	// when
	statusCode, _ := serveHTTP(t, s, listKeysRequest(t, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceListPublicKeysFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, listKeysRequest(t, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceGetPublicKeyOK(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	hdPubKey := &wallet.HDPublicKey{
		Idx:       1,
		PublicKey: vgrand.RandomStr(5),
		Algorithm: wallet.Algorithm{
			Name:    "some/algo",
			Version: 1,
		},
		Tainted:  false,
		MetaList: []wallet.Meta{{Key: "a", Value: "b"}},
	}
	headers := authHeaders(t, token)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().GetPublicKey(walletName, hdPubKey.PublicKey).Times(1).Return(hdPubKey, nil)

	// when
	statusCode, _ := serveHTTP(t, s, getKeyRequest(t, hdPubKey.PublicKey, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceGetPublicKeyFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(t, s, getKeyRequest(t, vgrand.RandomStr(5), tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceGetPublicKeyFailKeyNotFound(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	pubKey := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().GetPublicKey(walletName, pubKey).Times(1).Return(nil, wallet.ErrPubKeyDoesNotExist)

	// when
	statusCode, _ := serveHTTP(t, s, getKeyRequest(t, pubKey, headers))

	// then
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func testServiceGetPublicKeyFailMiscError(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	pubKey := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().GetPublicKey(walletName, pubKey).Times(1).Return(nil, assert.AnError)

	// when
	statusCode, _ := serveHTTP(t, s, getKeyRequest(t, pubKey, headers))

	// then
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func testServiceTaintOK(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	pubKey := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"passphrase": "%s"}`, passphrase)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().TaintKey(walletName, pubKey, passphrase).Times(1).Return(nil)

	// when
	statusCode, _ := serveHTTP(t, s, taintKeyRequest(t, pubKey, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceTaintFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
		payload string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
			payload: `{"passphrase": "some data"}`,
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
			payload: `{"passphrase": "some data"}`,
		}, {
			name:    "misspelled passphrase property",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"passhp": "some data"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, taintKeyRequest(tt, vgrand.RandomStr(5), tc.payload, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testServiceUpdateMetaOK(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// when
	walletName := vgrand.RandomStr(5)
	pubKey := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	passphrase := vgrand.RandomStr(5)
	metaRole := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"passphrase": "%s", "meta": [{"key":"role", "value":"%s"}]}`, passphrase, metaRole)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().UpdateMeta(walletName, pubKey, passphrase, []wallet.Meta{{
		Key:   "role",
		Value: metaRole,
	}}).Times(1).Return(nil)

	// when
	statusCode, _ := serveHTTP(t, s, annotateKeyRequest(t, pubKey, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testServiceUpdateMetaFailInvalidRequest(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
		payload string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
			payload: `{"passphrase": "some data", "meta": [{"key": "role", "value": "signing"}]}`,
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
			payload: `{"passphrase": "some data", "meta": [{"key": "role", "value": "signing"}]}`,
		}, {
			name:    "misspelled passphrase property",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"pssphrse": "some data", "meta": [{"key": "role", "value": "signing"}]}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, annotateKeyRequest(tt, vgrand.RandomStr(5), tc.payload, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testSigningTransactionSucceeds(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"pubKey": "%s", "orderCancellation": {}}`, vgrand.RandomStr(5))

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignTx(walletName, gomock.Any(), gomock.Any()).Times(1).Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(0)
	s.nodeForward.EXPECT().LastBlockHeightAndHash(gomock.Any()).Times(1).Return(&api.LastBlockHeightResponse{
		Height:              42,
		Hash:                "0292041e2f0cf741894503fb3ead4cb817bca2375e543aa70f7c4d938157b5a6",
		SpamPowDifficulty:   2,
		SpamPowHashFunction: "sha3_24_rounds",
	}, 0, nil)

	// when
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testAcceptSigningTransactionManuallySucceeds(t *testing.T) {
	s := getTestService(t, "manual")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	pubKey := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"pubKey": "%s", "orderCancellation": {}}`, pubKey)
	txStr := fmt.Sprintf(`pub_key:"%s" order_cancellation:{}`, pubKey)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignTx(walletName, gomock.Any(), gomock.Any()).Times(1).Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(0)
	s.nodeForward.EXPECT().LastBlockHeightAndHash(gomock.Any()).Times(1).Return(&api.LastBlockHeightResponse{
		Height:              42,
		Hash:                "0292041e2f0cf741894503fb3ead4cb817bca2375e543aa70f7c4d938157b5a6",
		SpamPowDifficulty:   2,
		SpamPowHashFunction: "sha3_24_rounds",
	}, 0, nil)

	// when
	s.ConsentConfirmations <- service.ConsentConfirmation{TxStr: txStr, Decision: true}
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testDeclineSigningTransactionManuallySucceeds(t *testing.T) {
	s := getTestService(t, "manual")
	defer s.ctrl.Finish()

	// given
	token := vgrand.RandomStr(5)
	walletName := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	pubKey := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"pubKey": "%s", "orderCancellation": {}}`, pubKey)
	txStr := fmt.Sprintf(`pub_key:"%s" order_cancellation:{}`, pubKey)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(0)

	// when
	s.ConsentConfirmations <- service.ConsentConfirmation{TxStr: txStr, Decision: false}
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusUnauthorized, statusCode)
}

func testSigningTransactionWithPropagationSucceeds(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"propagate": true, "pubKey": "%s", "orderCancellation": {}}`, vgrand.RandomStr(5))

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignTx(walletName, gomock.Any(), gomock.Any()).Times(1).Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), gomock.Any(), api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(1)
	s.nodeForward.EXPECT().LastBlockHeightAndHash(gomock.Any()).Times(1).Return(&api.LastBlockHeightResponse{
		Height:              42,
		Hash:                "0292041e2f0cf741894503fb3ead4cb817bca2375e543aa70f7c4d938157b5a6",
		SpamPowDifficulty:   2,
		SpamPowHashFunction: "sha3_24_rounds",
	}, 0, nil)

	// when
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testSigningTransactionWithFailedPropagationFails(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"propagate": true, "pubKey": "%s", "orderCancellation": {}}`, vgrand.RandomStr(5))

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignTx(walletName, gomock.Any(), gomock.Any()).Times(1).Return(&commandspb.Transaction{}, nil)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), gomock.Any(), api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(1).Return("", assert.AnError)
	s.nodeForward.EXPECT().LastBlockHeightAndHash(gomock.Any()).Times(1).Return(&api.LastBlockHeightResponse{
		Height:              42,
		Hash:                "0292041e2f0cf741894503fb3ead4cb817bca2375e543aa70f7c4d938157b5a6",
		SpamPowDifficulty:   2,
		SpamPowHashFunction: "sha3_24_rounds",
	}, 0, nil)

	// when
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func testFailedTransactionSigningFails(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"propagate": true, "pubKey": "%s", "orderCancellation": {}}`, vgrand.RandomStr(5))

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignTx(walletName, gomock.Any(), gomock.Any()).Times(1).Return(nil, assert.AnError)
	s.nodeForward.EXPECT().SendTx(gomock.Any(), &commandspb.Transaction{}, api.SubmitTransactionRequest_TYPE_ASYNC, gomock.Any()).Times(0)
	s.nodeForward.EXPECT().LastBlockHeightAndHash(gomock.Any()).Times(1).Return(&api.LastBlockHeightResponse{
		Height:              42,
		Hash:                "0292041e2f0cf741894503fb3ead4cb817bca2375e543aa70f7c4d938157b5a6",
		SpamPowDifficulty:   2,
		SpamPowHashFunction: "sha3_24_rounds",
	}, 0, nil)

	// when
	statusCode, _ := serveHTTP(t, s, signTxRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func testSigningTransactionWithInvalidRequestFails(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
		payload string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
			payload: `{"propagate": true, "pubKey": "0xCAFEDUDE", "orderCancellation": {}}`,
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
			payload: `{"propagate": true, "pubKey": "0xCAFEDUDE", "orderCancellation": {}}`,
		}, {
			name:    "misspelled pubKey property",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"propagate": true, "puey": "0xCAFEDUDE", "orderCancellation": {}}`,
		}, {
			name:    "without command",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"propagate": true, "pubKey": "0xCAFEDUDE", "robMoney": {}}`,
		}, {
			name:    "with unknown command",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"propagate": true, "pubKey": "0xCAFEDUDE"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			walletName := vgrand.RandomStr(5)
			token := vgrand.RandomStr(5)

			s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)

			// when
			statusCode, _ := serveHTTP(tt, s, signTxRequest(tt, tc.payload, tc.headers))
			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testSigningAnythingSucceeds(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	walletName := vgrand.RandomStr(5)
	pubKey := vgrand.RandomStr(5)
	token := vgrand.RandomStr(5)
	headers := authHeaders(t, token)
	payload := fmt.Sprintf(`{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "%s"}`, pubKey)

	// setup
	s.auth.EXPECT().VerifyToken(token).Times(1).Return(walletName, nil)
	s.handler.EXPECT().SignAny(walletName, []byte("spice of dune"), pubKey).Times(1).Return([]byte("some sig"), nil)

	// when
	statusCode, _ := serveHTTP(t, s, signAnyRequest(t, payload, headers))

	// then
	assert.Equal(t, http.StatusOK, statusCode)
}

func testSigningAnyDataWithInvalidRequestFails(t *testing.T) {
	tcs := []struct {
		name    string
		headers map[string]string
		payload string
	}{
		{
			name:    "no header",
			headers: map[string]string{},
			payload: `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd"}`,
		}, {
			name:    "no token",
			headers: authHeaders(t, ""),
			payload: `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd"}`,
		}, {
			name:    "misspelled pubKey property",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "puey": "asdasasdasd"}`,
		}, {
			name:    "misspelled inputData property",
			headers: authHeaders(t, vgrand.RandomStr(5)),
			payload: `{"data": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, signAnyRequest(tt, tc.payload, tc.headers))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func testVerifyingAnythingSucceeds(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	pubKey := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "%s", "signature": "U2lldGNoIFRhYnI="}`, pubKey)

	// setup
	s.handler.EXPECT().VerifyAny([]byte("spice of dune"), []byte("Sietch Tabr"), pubKey).Times(1).Return(true, nil)

	// when
	statusCode, body := serveHTTP(t, s, verifyAnyRequest(t, payload))

	// then
	assert.Equal(t, http.StatusOK, statusCode)

	resp := &service.VerifyAnyResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("couldn't unmarshal responde: %v", err)
	}
	assert.True(t, resp.Valid)
}

func testVerifyingAnythingFails(t *testing.T) {
	s := getTestService(t, "automatic")
	defer s.ctrl.Finish()

	// given
	pubKey := vgrand.RandomStr(5)
	payload := fmt.Sprintf(`{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "%s", "signature": "U2lldGNoIFRhYnI="}`, pubKey)

	// setup
	s.handler.EXPECT().VerifyAny([]byte("spice of dune"), []byte("Sietch Tabr"), pubKey).Times(1).Return(false, nil)

	// when
	statusCode, body := serveHTTP(t, s, verifyAnyRequest(t, payload))

	// then
	assert.Equal(t, http.StatusOK, statusCode)

	resp := &service.VerifyAnyResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		t.Fatalf("couldn't unmarshal responde: %v", err)
	}
	assert.False(t, resp.Valid)
}

func testVerifyingAnyDataWithInvalidRequestFails(t *testing.T) {
	tcs := []struct {
		name    string
		payload string
	}{
		{
			name:    "misspelled pubKey property",
			payload: `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "puey": "asdasasdasd", "signature": "U2lldGNoIFRhYnI="}`,
		}, {
			name:    "misspelled inputData property",
			payload: `{"data": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd", "signature": "U2lldGNoIFRhYnI="}`,
		}, {
			name:    "misspelled signature property",
			payload: `{"inputData": "c3BpY2Ugb2YgZHVuZQ==", "pubKey": "asdasasdasd", "sign": "U2lldGNoIFRhYnI="}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			s := getTestService(tt, "automatic")
			tt.Cleanup(func() {
				s.ctrl.Finish()
			})

			// when
			statusCode, _ := serveHTTP(tt, s, verifyAnyRequest(tt, tc.payload))

			// then
			assert.Equal(tt, http.StatusBadRequest, statusCode)
		})
	}
}

func loginRequest(t *testing.T, payload string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/auth/token", payload, nil)
}

func logoutRequest(t *testing.T, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodDelete, "/api/v1/auth/token", "", headers)
}

func createWalletRequest(t *testing.T, payload string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/wallets", payload, nil)
}

func importWalletRequest(t *testing.T, payload string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/wallets/import", payload, nil)
}

func generateKeyRequest(t *testing.T, payload string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/keys", payload, headers)
}

func listKeysRequest(t *testing.T, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodGet, "/api/v1/keys", "", headers)
}

func getKeyRequest(t *testing.T, keyID string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodGet, fmt.Sprintf("/api/v1/keys/%s", keyID), "", headers)
}

func taintKeyRequest(t *testing.T, id, payload string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/keys/%s/taint", id), payload, headers)
}

func annotateKeyRequest(t *testing.T, id, payload string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPut, fmt.Sprintf("/api/v1/keys/%s/metadata", id), payload, headers)
}

func signTxRequest(t *testing.T, payload string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/command", payload, headers)
}

func signAnyRequest(t *testing.T, payload string, headers map[string]string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/sign", payload, headers)
}

func verifyAnyRequest(t *testing.T, payload string) *http.Request {
	t.Helper()
	return buildRequest(t, http.MethodPost, "/api/v1/verify", payload, nil)
}

func authHeaders(t *testing.T, token string) map[string]string {
	t.Helper()
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
}

func buildRequest(t *testing.T, method, path, payload string, headers map[string]string) *http.Request {
	t.Helper()

	ctx, cancelFn := context.WithTimeout(context.Background(), testRequestTimeout)
	t.Cleanup(func() {
		cancelFn()
	})

	req, _ := http.NewRequestWithContext(ctx, method, path, strings.NewReader(payload))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func serveHTTP(t *testing.T, s *testService, req *http.Request) (int, []byte) {
	t.Helper()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	resp := w.Result() // nolint:bodyclose
	defer func() {
		if err := w.Result().Body.Close(); err != nil {
			t.Fatalf("couldn't close response body: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("couldn't read body: %v", err)
	}

	return resp.StatusCode, body
}
