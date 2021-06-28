package service

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.vegaprotocol.io/go-wallet/commands"
	"code.vegaprotocol.io/go-wallet/config"
	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
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

	cfg         *config.Config
	log         *zap.Logger
	s           *http.Server
	handler     WalletHandler
	auth        Auth
	nodeForward NodeForward
}

// CreateLoginWalletRequest describes the request for CreateWallet, LoginWallet.
type CreateLoginWalletRequest struct {
	Wallet     string `json:"wallet"`
	Passphrase string `json:"passphrase"`
}

func ParseCreateLoginWalletRequest(r *http.Request) (*CreateLoginWalletRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &CreateLoginWalletRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.Wallet) == 0 {
		errs.AddForProperty("wallet", commands.ErrIsRequired)
	}

	if len(req.Passphrase) == 0 {
		errs.AddForProperty("passphrase", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// TaintKeyRequest describes the request for TaintKey.
type TaintKeyRequest struct {
	Passphrase string `json:"passphrase"`
}

func ParseTaintKeyRequest(r *http.Request, keyID string) (*TaintKeyRequest, commands.Errors) {
	errs := commands.NewErrors()

	if len(keyID) == 0 {
		errs.AddForProperty("keyid", commands.ErrIsRequired)
	}

	req := &TaintKeyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.Passphrase) == 0 {
		errs.AddForProperty("passphrase", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// GenKeyPairRequest describes the request for GenerateKeypair
type GenKeyPairRequest struct {
	Passphrase string        `json:"passphrase"`
	Meta       []wallet.Meta `json:"meta"`
}

func ParseGenKeyPairRequest(r *http.Request) (*GenKeyPairRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &GenKeyPairRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.Passphrase) == 0 {
		errs.AddForProperty("passphrase", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// UpdateMetaRequest describes the request for UpdateMeta.
type UpdateMetaRequest struct {
	Passphrase string        `json:"passphrase"`
	Meta       []wallet.Meta `json:"meta"`
}

func ParseUpdateMetaRequest(r *http.Request, keyID string) (*UpdateMetaRequest, commands.Errors) {
	errs := commands.NewErrors()

	if len(keyID) == 0 {
		errs.AddForProperty("keyid", commands.ErrIsRequired)
	}

	req := &UpdateMetaRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.Passphrase) == 0 {
		errs.AddForProperty("passphrase", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// SignAnyRequest describes the request for SignAny.
type SignAnyRequest struct {
	InputData string `json:"inputData"`
	PubKey    string `json:"pubKey"`
}

func ParseSignAnyRequest(r *http.Request) (*SignAnyRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &SignAnyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.InputData) == 0 {
		errs.AddForProperty("input_data", commands.ErrIsRequired)
	}

	if len(req.PubKey) == 0 {
		errs.AddForProperty("pub_key", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// VerifyAnyRequest describes the request for VerifyAny.
type VerifyAnyRequest struct {
	// InputData is the payload to be verified. It should be base64 encoded.
	InputData string `json:"inputData"`
	// Signature is the signature to check against the InputData. It should be
	// base64 encoded.
	Signature string `json:"signature"`
	// PubKey is the public key used along the signature to check the InputData.
	PubKey    string `json:"pubKey"`
}

func ParseVerifyAnyRequest(r *http.Request) (*VerifyAnyRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &VerifyAnyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.InputData) == 0 {
		errs.AddForProperty("input_data", commands.ErrIsRequired)
	}

	if len(req.Signature) == 0 {
		errs.AddForProperty("signature", commands.ErrIsRequired)
	}

	if len(req.PubKey) == 0 {
		errs.AddForProperty("pub_key", commands.ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

func ParseSubmitTransactionRequest(r *http.Request) (*walletpb.SubmitTransactionRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &walletpb.SubmitTransactionRequest{}
	if err := jsonpb.Unmarshal(r.Body, req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if errs = wallet.CheckSubmitTransactionRequest(req); !errs.Empty() {
		return nil, errs
	}

	return req, nil
}

// KeyResponse describes the response to a request that returns a single key.
type KeyResponse struct {
	Key wallet.PublicKey `json:"key"`
}

// KeysResponse describes the response to a request that returns a list of keys.
type KeysResponse struct {
	Keys []wallet.PublicKey `json:"keys"`
}

// SignTxResponse describes the response for SignTx.
type SignTxResponse struct {
	SignedTx     wallet.SignedBundle `json:"signedTx"`
	HexBundle    string              `json:"hexBundle"`
	Base64Bundle string              `json:"base64Bundle"`
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
//go:generate go run github.com/golang/mock/mockgen -destination mocks/wallet_handler_mock.go -package mocks code.vegaprotocol.io/go-wallet/service WalletHandler
type WalletHandler interface {
	CreateWallet(name, passphrase string) error
	LoginWallet(name, passphrase string) error
	SecureGenerateKeyPair(name, passphrase string) (string, error)
	GetPublicKey(name, pubKey string) (*wallet.PublicKey, error)
	ListPublicKeys(name string) ([]wallet.PublicKey, error)
	SignTx(name, tx, pubKey string, height uint64) (wallet.SignedBundle, error)
	SignTxV2(name string, req *walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error)
	SignAny(name, inputData, pubKey string) ([]byte, error)
	VerifyAny(name, inputData, sig, pubKey string) (bool, error)
	TaintKey(name, pubKey, passphrase string) error
	UpdateMeta(name, pubKey, passphrase string, meta []wallet.Meta) error
	GetWalletPath(name string) (string, error)
}

// Auth ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/auth_mock.go -package mocks code.vegaprotocol.io/go-wallet/service Auth
type Auth interface {
	NewSession(name string) (string, error)
	VerifyToken(token string) (string, error)
	Revoke(token string) error
}

// NodeForward ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/node_forward_mock.go -package mocks code.vegaprotocol.io/go-wallet/service NodeForward
type NodeForward interface {
	Send(context.Context, *wallet.SignedBundle, api.SubmitTransactionRequest_Type) error
	SendTxV2(context.Context, *commandspb.Transaction, api.SubmitTransactionV2Request_Type) error
	HealthCheck(context.Context) error
	LastBlockHeight(context.Context) (uint64, error)
}

func NewService(log *zap.Logger, cfg *config.Config, rootPath string) (*Service, error) {
	log = log.Named("wallet")

	store, err := storev1.NewStore(rootPath)
	if err != nil {
		return nil, err
	}
	auth, err := NewAuth(log, store, cfg.TokenExpiry.Get())
	if err != nil {
		return nil, err
	}
	nodeForward, err := newNodeForward(log, cfg.Nodes)
	if err != nil {
		return nil, err
	}
	handler := wallet.NewHandler(store)
	return NewServiceWith(log, cfg, handler, auth, nodeForward)
}

func NewServiceWith(log *zap.Logger, cfg *config.Config, h WalletHandler, a Auth, n NodeForward) (*Service, error) {
	s := &Service{
		Router:      httprouter.New(),
		log:         log,
		cfg:         cfg,
		handler:     h,
		auth:        a,
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
	s.POST("/api/v1/verify", ExtractToken(s.VerifyAny))
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
		Handler: cors.AllowAll().Handler(s),
	}

	return s.s.ListenAndServe()
}

func (s *Service) Stop() error {
	return s.s.Shutdown(context.Background())
}

func (s *Service) CreateWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseCreateLoginWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	err := s.handler.CreateWallet(req.Wallet, req.Passphrase)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, TokenResponse{Token: token}, http.StatusOK)
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseCreateLoginWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	err := s.handler.LoginWallet(req.Wallet, req.Passphrase)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, TokenResponse{Token: token}, http.StatusOK)
}

func (s *Service) DownloadWallet(token string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, err := s.auth.VerifyToken(token)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	path, err := s.handler.GetWalletPath(name)
	if err != nil {
		s.writeBadRequest(w, commands.NewErrors().FinalAdd(err))
		return
	}

	http.ServeFile(w, r, path)
}

func (s *Service) Revoke(t string, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := s.auth.Revoke(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) GenerateKeypair(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseGenKeyPairRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	pubKey, err := s.handler.SecureGenerateKeyPair(name, req.Passphrase)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	// if any meta specified, lets add them
	if len(req.Meta) > 0 {
		err := s.handler.UpdateMeta(name, pubKey, req.Passphrase, req.Meta)
		if err != nil {
			writeForbiddenError(w, err)
			return
		}
	}

	key, err := s.handler.GetPublicKey(name, pubKey)
	if err != nil {
		s.writeBadRequest(w, commands.NewErrors().FinalAdd(err))
		return
	}

	writeSuccess(w, KeyResponse{Key: *key}, http.StatusOK)
}

func (s *Service) GetPublicKey(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	key, err := s.handler.GetPublicKey(name, ps.ByName("keyid"))
	if err != nil {
		var statusCode int
		if err == wallet.ErrPubKeyDoesNotExist {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusForbidden
		}
		writeError(w, newErrorResponse(err.Error()), statusCode)
		return
	}

	writeSuccess(w, KeyResponse{Key: *key}, http.StatusOK)
}

func (s *Service) ListPublicKeys(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	keys, err := s.handler.ListPublicKeys(name)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, KeysResponse{Keys: keys}, http.StatusOK)
}

func (s *Service) TaintKey(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keyID := ps.ByName("keyid")
	req, errs := ParseTaintKeyRequest(r, keyID)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	err = s.handler.TaintKey(name, keyID, req.Passphrase)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) UpdateMeta(t string, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	keyID := ps.ByName("keyid")
	req, errs := ParseUpdateMetaRequest(r, keyID)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	err = s.handler.UpdateMeta(name, keyID, req.Passphrase, req.Meta)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) SignAny(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseSignAnyRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	signature, err := s.handler.SignAny(name, req.InputData, req.PubKey)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	res := SignAnyResponse{
		HexSignature:    hex.EncodeToString(signature),
		Base64Signature: base64.StdEncoding.EncodeToString(signature),
	}

	writeSuccess(w, res, http.StatusOK)
}

func (s *Service) VerifyAny(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseVerifyAnyRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	verified, err := s.handler.VerifyAny(name, req.InputData, req.Signature, req.PubKey)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	writeSuccess(w, SuccessResponse{Success: verified}, http.StatusOK)
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

	req, errs := ParseSubmitTransactionRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	height, err := s.nodeForward.LastBlockHeight(r.Context())
	if err != nil {
		writeInternalError(w, ErrCouldNotGetBlockHeight)
		return
	}

	name, err := s.auth.VerifyToken(token)
	if err != nil {
		writeForbiddenError(w, err)
		return
	}

	tx, err := s.handler.SignTxV2(name, req, height)
	if err != nil {
		writeForbiddenError(w, err)
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
				writeInternalError(w, err)
			}
			return
		}
	}

	writeSuccessProto(w, tx, http.StatusOK)
}

func (s *Service) health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := s.nodeForward.HealthCheck(r.Context()); err != nil {
		writeSuccess(w, SuccessResponse{Success: false}, http.StatusFailedDependency)
		return
	}
	writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) writeBadRequest(w http.ResponseWriter, errs commands.Errors) {
	s.writeErrors(w, http.StatusBadRequest, errs)
}

func (s *Service) writeErrors(w http.ResponseWriter, statusCode int, errs commands.Errors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	buf, _ := json.Marshal(ErrorsResponse{Errors: errs})
	if _, err := w.Write(buf); err != nil {
		s.log.Error(fmt.Sprintf("couldn't marshal errors as JSON because of: %s", err.Error()),
			zap.Error(errs),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func unmarshalBody(r *http.Request, into interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrCouldNotReadRequest
	}
	return json.Unmarshal(body, into)
}

func writeForbiddenError(w http.ResponseWriter, e error) {
	writeError(w, newErrorResponse(e.Error()), http.StatusForbidden)
}

func writeInternalError(w http.ResponseWriter, e error) {
	writeError(w, newErrorResponse(e.Error()), http.StatusInternalServerError)
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
