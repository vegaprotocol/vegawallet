package wallet

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.vegaprotocol.io/go-wallet/commands"
	typespb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto"
	"github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/api"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	walletpb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/wallet/v1"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type Service struct {
	*httprouter.Router

	cfg         *Config
	log         *zap.Logger
	s           *http.Server
	handler     WalletHandler
	nodeForward NodeForward
}

// CreateLoginWalletRequest describes the request for CreateWallet, LoginWallet.
type CreateLoginWalletRequest struct {
	Wallet     string `json:"wallet"`
	Passphrase string `json:"passphrase"`
}

// PassphraseRequest describes the request for TaintKey.
type PassphraseRequest struct {
	Passphrase string `json:"passphrase"`
}

// PassphraseMetaRequest describes the request for GenerateKeypair, UpdateMeta.
type PassphraseMetaRequest struct {
	Passphrase string `json:"passphrase"`
	Meta       []Meta `json:"meta"`
}

// SignTxRequest describes the request for SignTx.
type SignTxRequest struct {
	Tx        string `json:"tx"`
	PubKey    string `json:"pubKey"`
	Propagate bool   `json:"propagate"`
}

// SignAnyRequest describes the request for SignAny.
type SignAnyRequest struct {
	InputData string `json:"inputData"`
	PubKey    string `json:"pubKey"`
}

// KeyResponse describes the response to a request that returns a single key.
type KeyResponse struct {
	Key Keypair `json:"key"`
}

// KeysResponse describes the response to a request that returns a list of keys.
type KeysResponse struct {
	Keys []Keypair `json:"keys"`
}

// SignTxResponse describes the response for SignTx.
type SignTxResponse struct {
	SignedTx     SignedBundle `json:"signedTx"`
	HexBundle    string       `json:"hexBundle"`
	Base64Bundle string       `json:"base64Bundle"`
}

// SignAnyResponse describes the response for SignAny.
type SignAnyResponse struct {
	HexSignature    string `json:"hexSignature"`
	Base64Signature string `json:"base64Signature"`
}

// SuccessResponse describes the response to a request that returns a simple true/false answer.
type SuccessResponse struct {
	Success bool `json:"success"`
}

// TokenResponse describes the response to a request that returns a token.
type TokenResponse struct {
	Token string `json:"token"`
}

// WalletHandler ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/wallet_handler_mock.go -package mocks code.vegaprotocol.io/go-wallet/wallet WalletHandler
type WalletHandler interface {
	CreateWallet(name, passphrase string) (string, error)
	LoginWallet(name, passphrase string) (string, error)
	RevokeToken(token string) error
	GenerateKeypair(token, passphrase string) (string, error)
	GetPublicKey(token, pubKey string) (*Keypair, error)
	GetWalletName(token string) (string, error)
	ListPublicKeys(token string) ([]Keypair, error)
	SignTx(token, tx, pubKey string, height uint64) (SignedBundle, error)
	SignTxV2(token string, req walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error)
	SignAny(token, inputData, pubKey string) ([]byte, error)
	TaintKey(token, pubKey, passphrase string) error
	UpdateMeta(token, pubKey, passphrase string, meta []Meta) error
	GetWalletPath(token string) (string, error)
}

// NodeForward ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/node_forward_mock.go -package mocks code.vegaprotocol.io/go-wallet/wallet NodeForward
type NodeForward interface {
	Send(context.Context, *SignedBundle, api.SubmitTransactionRequest_Type) error
	SendTxV2(context.Context, *commandspb.Transaction, api.SubmitTransactionV2Request_Type) error
	HealthCheck(context.Context) error
	LastBlockHeight(context.Context) (uint64, error)
}

func NewService(log *zap.Logger, cfg *Config, rootPath string) (*Service, error) {
	log = log.Named(namedLogger)

	fileStore, err := NewFileStoreV1(rootPath)
	if err != nil {
		return nil, err
	}
	auth, err := NewAuth(log, rootPath, cfg.TokenExpiry.Get())
	if err != nil {
		return nil, err
	}
	nodeForward, err := NewNodeForward(log, cfg.Nodes)
	if err != nil {
		return nil, err
	}
	handler := NewHandler(auth, fileStore)
	return NewServiceWith(log, cfg, handler, nodeForward)
}

func NewServiceWith(log *zap.Logger, cfg *Config, h WalletHandler, n NodeForward) (*Service, error) {
	s := &Service{
		Router:      httprouter.New(),
		log:         log,
		cfg:         cfg,
		handler:     h,
		nodeForward: n,
	}

	s.POST("/api/v1/auth/token", s.Login)
	s.GET("/api/v1/status", s.health)
	s.POST("/api/v1/wallets", s.CreateWallet)

	s.DELETE("/api/v1/auth/token", ExtractToken(s.Revoke))
	s.GET("/api/v1/keys", ExtractToken(s.ListPublicKeys))
	s.POST("/api/v1/keys", ExtractToken(s.GenerateKeypair))
	s.GET("/api/v1/keys/:keyid", ExtractToken(s.GetPublicKey))
	s.PUT("/api/v1/keys/:keyid/taint", ExtractToken(s.TaintKey))
	s.PUT("/api/v1/keys/:keyid/metadata", ExtractToken(s.UpdateMeta))
	s.POST("/api/v1/sign", ExtractToken(s.SignAny))
	s.POST("/api/v1/command", ExtractToken(s.SignTxV2))
	s.POST("/api/v1/command/sync", ExtractToken(s.SignTxSyncV2))
	s.POST("/api/v1/command/commit", ExtractToken(s.SignTxCommitV2))
	s.GET("/api/v1/wallets", ExtractToken(s.DownloadWallet))

	// DEPRECATED Use
	s.POST("/api/v1/messages", ExtractToken(s.SignTx))
	s.POST("/api/v1/messages/async", ExtractToken(s.SignTxAsync))
	s.POST("/api/v1/messages/commit", ExtractToken(s.SignTxCommit))

	return s, nil
}

func (s *Service) Start() error {
	s.s = &http.Server{
		Addr:    fmt.Sprintf("%s:%v", s.cfg.Host, s.cfg.Port),
		Handler: cors.AllowAll().Handler(s), // middlewar with cors
	}

	return s.s.ListenAndServe()
}

func (s *Service) Stop() error {
	return s.s.Shutdown(context.Background())
}

func (s *Service) CreateWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// unmarshal request
	req := CreateLoginWalletRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}

	// validation
	if len(req.Wallet) <= 0 {
		writeError(w, newError("missing wallet field"), http.StatusBadRequest)
		return
	}
	if len(req.Passphrase) <= 0 {
		writeError(w, newError("missing passphrase field"), http.StatusBadRequest)
		return
	}

	token, err := s.handler.CreateWallet(req.Wallet, req.Passphrase)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	writeSuccess(w, TokenResponse{token}, http.StatusOK)
}

func (s *Service) DownloadWallet(token string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	path, err := s.handler.GetWalletPath(token)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}
	http.ServeFile(w, r, path)
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := CreateLoginWalletRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	// validation
	if len(req.Wallet) <= 0 {
		writeError(w, newError("missing wallet field"), http.StatusBadRequest)
		return
	}
	if len(req.Passphrase) <= 0 {
		writeError(w, newError("missing passphrase field"), http.StatusBadRequest)
		return
	}

	token, err := s.handler.LoginWallet(req.Wallet, req.Passphrase)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}
	writeSuccess(w, TokenResponse{token}, http.StatusOK)
}

func (s *Service) Revoke(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := s.handler.RevokeToken(t)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) GenerateKeypair(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := PassphraseMetaRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}
	if len(req.Passphrase) <= 0 {
		writeError(w, newError("missing passphrase field"), http.StatusBadRequest)
		return
	}

	pubKey, err := s.handler.GenerateKeypair(t, req.Passphrase)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	// if any meta specified, lets add them
	if len(req.Meta) > 0 {
		err := s.handler.UpdateMeta(t, pubKey, req.Passphrase, req.Meta)
		if err != nil {
			writeError(w, newError(err.Error()), http.StatusForbidden)
			return
		}
	}

	key, err := s.handler.GetPublicKey(t, pubKey)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}

	writeSuccess(w, KeyResponse{Key: *key}, http.StatusOK)
}

func (s *Service) GetPublicKey(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key, err := s.handler.GetPublicKey(t, ps.ByName("keyid"))
	if err != nil {
		var statusCode int
		if err == ErrPubKeyDoesNotExist {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusForbidden
		}
		writeError(w, newError(err.Error()), statusCode)
		return
	}

	writeSuccess(w, KeyResponse{Key: *key}, http.StatusOK)
}

func (s *Service) ListPublicKeys(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	keys, err := s.handler.ListPublicKeys(t)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	writeSuccess(w, KeysResponse{Keys: keys}, http.StatusOK)
}

func (s *Service) TaintKey(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := PassphraseRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}
	keyID := ps.ByName("keyid")
	if len(keyID) <= 0 {
		writeError(w, newError("missing keyID"), http.StatusBadRequest)
		return
	}
	if len(req.Passphrase) <= 0 {
		writeError(w, newError("missing passphrase field"), http.StatusBadRequest)
		return
	}

	err := s.handler.TaintKey(t, keyID, req.Passphrase)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) UpdateMeta(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := PassphraseMetaRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}
	keyID := ps.ByName("keyid")
	if len(keyID) <= 0 {
		writeError(w, newError("missing keyID"), http.StatusBadRequest)
		return
	}
	if len(req.Passphrase) <= 0 {
		writeError(w, newError("missing passphrase field"), http.StatusBadRequest)
		return
	}

	err := s.handler.UpdateMeta(t, keyID, req.Passphrase, req.Meta)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) SignAny(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := SignAnyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		writeError(w, newError(err.Error()), http.StatusBadRequest)
		return
	}
	if len(req.InputData) <= 0 {
		writeError(w, newError("missing inputData field"), http.StatusBadRequest)
		return
	}
	if len(req.PubKey) <= 0 {
		writeError(w, newError("missing pubKey field"), http.StatusBadRequest)
		return
	}

	signature, err := s.handler.SignAny(t, req.InputData, req.PubKey)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	res := SignAnyResponse{
		HexSignature:    hex.EncodeToString(signature),
		Base64Signature: base64.StdEncoding.EncodeToString(signature),
	}

	writeSuccess(w, res, http.StatusOK)
}

func (s *Service) SignTxSyncV2(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTxV2(token, w, r, p, api.SubmitTransactionV2Request_TYPE_SYNC)
}

func (s *Service) SignTxCommitV2(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTxV2(token, w, r, p, api.SubmitTransactionV2Request_TYPE_COMMIT)
}

func (s *Service) SignTxV2(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTxV2(token, w, r, p, api.SubmitTransactionV2Request_TYPE_ASYNC)
}

func (s *Service) signTxV2(token string, w http.ResponseWriter, r *http.Request, _ httprouter.Params, ty api.SubmitTransactionV2Request_Type) {
	defer r.Body.Close()

	errs := commands.NewErrors()
	req := walletpb.SubmitTransactionRequest{}

	if err := jsonpb.Unmarshal(r.Body, &req); err != nil {
		errs.Add(errors.New("couldn't parse the request"))
		s.writeBadRequest(w, errs)
		return
	}

	errs.Merge(CheckSubmitTransactionRequest(req))
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	height, err := s.nodeForward.LastBlockHeight(r.Context())
	if err != nil {
		writeError(w, newError("could not get last block height"), http.StatusInternalServerError)
		return
	}

	tx, err := s.handler.SignTxV2(token, req, height)
	if err != nil {
		writeError(w, newError(err.Error()), http.StatusForbidden)
		return
	}

	if req.Propagate {
		if err := s.nodeForward.SendTxV2(r.Context(), tx, ty); err != nil {
			if s, ok := status.FromError(err); ok {
				details := []string{}
				for _, v := range s.Details() {
					v := v.(*typespb.ErrorDetail)
					details = append(details, v.Message)
				}
				writeError(w, newErrorWithDetails(err.Error(), details), http.StatusInternalServerError)
			} else {
				writeError(w, newError(err.Error()), http.StatusInternalServerError)
			}
			return
		}
	}

	writeSuccessProto(w, tx, http.StatusOK)
}

type ErrorResponse struct {
	Errors commands.Errors `json:"errors"`
}

func (s *Service) writeBadRequest(w http.ResponseWriter, errs commands.Errors) {
	s.writeErrors(w, http.StatusBadRequest, errs)
}

func (s *Service) writeErrors(w http.ResponseWriter, statusCode int, errs commands.Errors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	buf, _ := json.Marshal(ErrorResponse{Errors: errs})
	if _, err := w.Write(buf); err != nil {
		s.log.Error(fmt.Sprintf("couldn't marshal errors as JSON because of: %s", err.Error()),
			zap.Error(errs),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := s.nodeForward.HealthCheck(r.Context()); err != nil {
		writeSuccess(w, SuccessResponse{Success: false}, http.StatusFailedDependency)
		return
	}
	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func unmarshalBody(r *http.Request, into interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrInvalidRequest
	}
	return json.Unmarshal(body, into)
}

func writeError(w http.ResponseWriter, e error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	buf, _ := json.Marshal(e)
	w.Write(buf)
}

func writeSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	buf, _ := json.Marshal(data)
	w.Write(buf)
}

func writeSuccessProto(w http.ResponseWriter, data proto.Message, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	marshaler := jsonpb.Marshaler{}
	marshaler.Marshal(w, data)
}

var (
	ErrInvalidRequest        = newError("invalid request")
	ErrInvalidOrMissingToken = newError("invalid or missing token")
)

type HTTPError struct {
	ErrorStr string   `json:"error"`
	Details  []string `json:"details"`
}

func (e HTTPError) Error() string {
	return e.ErrorStr
}

func newError(e string) HTTPError {
	return HTTPError{
		ErrorStr: e,
	}
}

func newErrorWithDetails(e string, details []string) HTTPError {
	return HTTPError{
		ErrorStr: e,
		Details:  details,
	}
}
