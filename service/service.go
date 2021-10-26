package service

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.vegaprotocol.io/protos/commands"
	typespb "code.vegaprotocol.io/protos/vega"
	api "code.vegaprotocol.io/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	wcommands "code.vegaprotocol.io/vegawallet/commands"
	"code.vegaprotocol.io/vegawallet/network"
	"code.vegaprotocol.io/vegawallet/version"
	"code.vegaprotocol.io/vegawallet/wallet"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type Service struct {
	*httprouter.Router

	network     *network.Network
	log         *zap.Logger
	server      *http.Server
	handler     WalletHandler
	auth        Auth
	nodeForward NodeForward
}

// CreateWalletRequest describes the request for CreateWallet.
type CreateWalletRequest struct {
	Wallet     string `json:"wallet"`
	Passphrase string `json:"passphrase"`
}

func ParseCreateWalletRequest(r *http.Request) (*CreateWalletRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &CreateWalletRequest{}
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

// CreateWalletResponse returns the authentication token and the auto-generated
// mnemonic of the created wallet.
type CreateWalletResponse struct {
	Mnemonic string `json:"mnemonic"`
	Token    string `json:"token"`
}

// ImportWalletRequest describes the request for ImportWallet.
type ImportWalletRequest struct {
	Wallet     string `json:"wallet"`
	Passphrase string `json:"passphrase"`
	Mnemonic   string `json:"mnemonic"`
	Version    uint32 `json:"version"`
}

func ParseImportWalletRequest(r *http.Request) (*ImportWalletRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &ImportWalletRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.Wallet) == 0 {
		errs.AddForProperty("wallet", commands.ErrIsRequired)
	}

	if len(req.Passphrase) == 0 {
		errs.AddForProperty("passphrase", commands.ErrIsRequired)
	}

	if len(req.Mnemonic) == 0 {
		errs.AddForProperty("mnemonic", commands.ErrIsRequired)
	}

	if req.Version == 0 {
		req.Version = wallet.LatestVersion
	}

	if !errs.Empty() {
		return nil, errs
	}

	return req, errs
}

// LoginWalletRequest describes the request for CreateWallet, LoginWallet.
type LoginWalletRequest struct {
	Wallet     string `json:"wallet"`
	Passphrase string `json:"passphrase"`
}

func ParseLoginWalletRequest(r *http.Request) (*LoginWalletRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &LoginWalletRequest{}
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

// GenKeyPairRequest describes the request for GenerateKeyPair.
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
	// InputData is the payload to generate a signature from. I should be
	// base 64 encoded.
	InputData string `json:"inputData"`
	// PubKey is used to retrieve the private key to sign the InputDate.
	PubKey string `json:"pubKey"`

	decodedInputData []byte
}

func ParseSignAnyRequest(r *http.Request) (*SignAnyRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &SignAnyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.InputData) == 0 {
		errs.AddForProperty("inputData", commands.ErrIsRequired)
	}
	decodedInputData, err := base64.StdEncoding.DecodeString(req.InputData)
	if err != nil {
		errs.AddForProperty("inputData", ErrShouldBeBase64Encoded)
	} else {
		req.decodedInputData = decodedInputData
	}

	if len(req.PubKey) == 0 {
		errs.AddForProperty("pubKey", commands.ErrIsRequired)
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
	PubKey string `json:"pubKey"`

	decodedInputData []byte
	decodedSignature []byte
}

func ParseVerifyAnyRequest(r *http.Request) (*VerifyAnyRequest, commands.Errors) {
	errs := commands.NewErrors()

	req := &VerifyAnyRequest{}
	if err := unmarshalBody(r, &req); err != nil {
		return nil, errs.FinalAdd(err)
	}

	if len(req.InputData) == 0 {
		errs.AddForProperty("inputData", commands.ErrIsRequired)
	} else {
		decodedInputData, err := base64.StdEncoding.DecodeString(req.InputData)
		if err != nil {
			errs.AddForProperty("inputData", ErrShouldBeBase64Encoded)
		} else {
			req.decodedInputData = decodedInputData
		}
	}

	if len(req.Signature) == 0 {
		errs.AddForProperty("signature", commands.ErrIsRequired)
	} else {
		decodedSignature, err := base64.StdEncoding.DecodeString(req.Signature)
		if err != nil {
			errs.AddForProperty("signature", ErrShouldBeBase64Encoded)
		} else {
			req.decodedSignature = decodedSignature
		}
	}

	if len(req.PubKey) == 0 {
		errs.AddForProperty("pubKey", commands.ErrIsRequired)
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

	if errs = wcommands.CheckSubmitTransactionRequest(req); !errs.Empty() {
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

// VersionResponse describes the response to a request that returns app version info.
type VersionResponse struct {
	Version     string `json:"version"`
	VersionHash string `json:"versionHash"`
}

// NetworkResponse describes the response to a request that returns app hosts info.
type NetworkResponse struct {
	Network network.Network `json:"network"`
}

// WalletHandler ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/wallet_handler_mock.go -package mocks code.vegaprotocol.io/vegawallet/service WalletHandler
type WalletHandler interface {
	CreateWallet(name, passphrase string) (string, error)
	ImportWallet(name, passphrase, mnemonic string, version uint32) error
	LoginWallet(name, passphrase string) error
	LogoutWallet(name string)
	SecureGenerateKeyPair(name, passphrase string, meta []wallet.Meta) (string, error)
	GetPublicKey(name, pubKey string) (wallet.PublicKey, error)
	ListPublicKeys(name string) ([]wallet.PublicKey, error)
	SignTx(name string, req *walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error)
	SignAny(name string, inputData []byte, pubKey string) ([]byte, error)
	VerifyAny(inputData, sig []byte, pubKey string) (bool, error)
	TaintKey(name, pubKey, passphrase string) error
	UpdateMeta(name, pubKey, passphrase string, meta []wallet.Meta) error
}

// Auth ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/auth_mock.go -package mocks code.vegaprotocol.io/vegawallet/service Auth
type Auth interface {
	NewSession(name string) (string, error)
	VerifyToken(token string) (string, error)
	Revoke(token string) error
}

// NodeForward ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/node_forward_mock.go -package mocks code.vegaprotocol.io/vegawallet/service NodeForward
type NodeForward interface {
	SendTx(context.Context, *commandspb.Transaction, api.SubmitTransactionRequest_Type) error
	HealthCheck(context.Context) error
	LastBlockHeight(context.Context) (uint64, error)
}

func NewService(log *zap.Logger, net *network.Network, h WalletHandler, a Auth, n NodeForward) (*Service, error) {
	s := &Service{
		Router:      httprouter.New(),
		log:         log,
		handler:     h,
		auth:        a,
		nodeForward: n,
		network:     net,
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%v", net.Host, net.Port),
		Handler: cors.AllowAll().Handler(s),
	}

	s.POST("/api/v1/auth/token", s.Login)
	s.DELETE("/api/v1/auth/token", ExtractToken(s.Revoke))

	s.GET("/api/v1/network", s.GetNetwork)

	s.POST("/api/v1/wallets", s.CreateWallet)
	s.POST("/api/v1/wallets/import", s.ImportWallet)

	s.GET("/api/v1/keys", ExtractToken(s.ListPublicKeys))
	s.POST("/api/v1/keys", ExtractToken(s.GenerateKeyPair))
	s.GET("/api/v1/keys/:keyid", ExtractToken(s.GetPublicKey))
	s.PUT("/api/v1/keys/:keyid/taint", ExtractToken(s.TaintKey))
	s.PUT("/api/v1/keys/:keyid/metadata", ExtractToken(s.UpdateMeta))

	s.POST("/api/v1/command", ExtractToken(s.SignTx))
	s.POST("/api/v1/command/sync", ExtractToken(s.SignTxSync))
	s.POST("/api/v1/command/commit", ExtractToken(s.SignTxCommit))
	s.POST("/api/v1/sign", ExtractToken(s.SignAny))
	s.POST("/api/v1/verify", ExtractToken(s.VerifyAny))

	s.GET("/api/v1/version", s.Version)
	s.GET("/api/v1/status", s.Health)

	return s, nil
}

func (s *Service) Start() error {
	return s.server.ListenAndServe()
}

func (s *Service) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Service) CreateWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseCreateWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	mnemonic, err := s.handler.CreateWallet(req.Wallet, req.Passphrase)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, CreateWalletResponse{Mnemonic: mnemonic, Token: token}, http.StatusOK)
}

func (s *Service) ImportWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseImportWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	err := s.handler.ImportWallet(req.Wallet, req.Passphrase, req.Mnemonic, req.Version)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, TokenResponse{Token: token}, http.StatusOK)
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseLoginWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	err := s.handler.LoginWallet(req.Wallet, req.Passphrase)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, TokenResponse{Token: token}, http.StatusOK)
}

func (s *Service) Revoke(t string, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	err = s.auth.Revoke(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.handler.LogoutWallet(name)

	s.writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) GenerateKeyPair(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseGenKeyPairRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	pubKey, err := s.handler.SecureGenerateKeyPair(name, req.Passphrase, req.Meta)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	key, err := s.handler.GetPublicKey(name, pubKey)
	if err != nil {
		s.writeBadRequest(w, commands.NewErrors().FinalAdd(err))
		return
	}

	s.writeSuccess(w, KeyResponse{Key: key}, http.StatusOK)
}

func (s *Service) GetPublicKey(t string, w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	key, err := s.handler.GetPublicKey(name, ps.ByName("keyid"))
	if err != nil {
		var statusCode int
		if errors.Is(err, wallet.ErrPubKeyDoesNotExist) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusForbidden
		}
		s.writeError(w, newErrorResponse(err.Error()), statusCode)
		return
	}

	s.writeSuccess(w, KeyResponse{Key: key}, http.StatusOK)
}

func (s *Service) ListPublicKeys(t string, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	keys, err := s.handler.ListPublicKeys(name)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, KeysResponse{Keys: keys}, http.StatusOK)
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
		s.writeForbiddenError(w, err)
		return
	}

	err = s.handler.TaintKey(name, keyID, req.Passphrase)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
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
		s.writeForbiddenError(w, err)
		return
	}

	err = s.handler.UpdateMeta(name, keyID, req.Passphrase, req.Meta)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) SignAny(t string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseSignAnyRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	signature, err := s.handler.SignAny(name, req.decodedInputData, req.PubKey)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	res := SignAnyResponse{
		HexSignature:    hex.EncodeToString(signature),
		Base64Signature: base64.StdEncoding.EncodeToString(signature),
	}

	s.writeSuccess(w, res, http.StatusOK)
}

func (s *Service) VerifyAny(_ string, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseVerifyAnyRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	verified, err := s.handler.VerifyAny(req.decodedInputData, req.decodedSignature, req.PubKey)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.writeSuccess(w, SuccessResponse{Success: verified}, http.StatusOK)
}

func (s *Service) SignTxSync(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTx(token, w, r, p, api.SubmitTransactionRequest_TYPE_SYNC)
}

func (s *Service) SignTxCommit(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTx(token, w, r, p, api.SubmitTransactionRequest_TYPE_COMMIT)
}

func (s *Service) SignTx(token string, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.signTx(token, w, r, p, api.SubmitTransactionRequest_TYPE_ASYNC)
}

func (s *Service) signTx(token string, w http.ResponseWriter, r *http.Request, _ httprouter.Params, ty api.SubmitTransactionRequest_Type) {
	defer r.Body.Close()

	req, errs := ParseSubmitTransactionRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	height, err := s.nodeForward.LastBlockHeight(r.Context())
	if err != nil {
		s.writeInternalError(w, ErrCouldNotGetBlockHeight)
		return
	}

	name, err := s.auth.VerifyToken(token)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	tx, err := s.handler.SignTx(name, req, height)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	if req.Propagate {
		if err := s.nodeForward.SendTx(r.Context(), tx, ty); err != nil {
			if st, ok := status.FromError(err); ok {
				var details []string
				for _, v := range st.Details() {
					v, ok := v.(*typespb.ErrorDetail)
					if !ok {
						s.writeError(w, newErrorResponse(fmt.Sprintf("couldn't cast status details to error details: %v", v)), http.StatusInternalServerError)
					}
					details = append(details, v.Message)
				}
				s.writeError(w, newErrorWithDetails(err.Error(), details), http.StatusInternalServerError)
			} else {
				s.writeInternalError(w, err)
			}
			return
		}
	}

	s.writeSuccessProto(w, tx, http.StatusOK)
}

func (s *Service) Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := VersionResponse{
		Version:     version.Version,
		VersionHash: version.VersionHash,
	}

	s.writeSuccess(w, res, http.StatusOK)
}

func (s *Service) GetNetwork(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := NetworkResponse{
		Network: *s.network,
	}
	s.writeSuccess(w, res, http.StatusOK)
}

func (s *Service) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := s.nodeForward.HealthCheck(r.Context()); err != nil {
		s.writeSuccess(w, SuccessResponse{Success: false}, http.StatusFailedDependency)
		return
	}
	s.writeSuccess(w, SuccessResponse{Success: true}, http.StatusOK)
}

func (s *Service) writeBadRequest(w http.ResponseWriter, errs commands.Errors) {
	s.writeErrors(w, http.StatusBadRequest, errs)
}

func (s *Service) writeErrors(w http.ResponseWriter, statusCode int, errs commands.Errors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	buf, _ := json.Marshal(ErrorsResponse{Errors: errs})
	if _, err := w.Write(buf); err != nil {
		s.log.Error("couldn't marshal error", zap.Error(errs))
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

func (s *Service) writeForbiddenError(w http.ResponseWriter, e error) {
	s.writeError(w, newErrorResponse(e.Error()), http.StatusForbidden)
}

func (s *Service) writeInternalError(w http.ResponseWriter, e error) {
	s.writeError(w, newErrorResponse(e.Error()), http.StatusInternalServerError)
}

func (s *Service) writeError(w http.ResponseWriter, e error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	buf, err := json.Marshal(e)
	if err != nil {
		s.log.Error("couldn't marshal error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		s.log.Error("couldn't write error to HTTP response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) writeSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	buf, err := json.Marshal(data)
	if err != nil {
		s.log.Error("couldn't marshal error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		s.log.Error("couldn't write error to HTTP response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Service) writeSuccessProto(w http.ResponseWriter, data proto.Message, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	marshaller := jsonpb.Marshaler{}
	if err := marshaller.Marshal(w, data); err != nil {
		s.log.Error("couldn't marshal proto message", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
