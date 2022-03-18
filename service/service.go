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
	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	wcommands "code.vegaprotocol.io/vegawallet/commands"
	"code.vegaprotocol.io/vegawallet/network"
	"code.vegaprotocol.io/vegawallet/version"
	"code.vegaprotocol.io/vegawallet/wallet"

	"github.com/golang/protobuf/jsonpb"
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
// recovery phrase of the created wallet.
type CreateWalletResponse struct {
	RecoveryPhrase string `json:"recoveryPhrase"`
	Token          string `json:"token"`
}

// ImportWalletRequest describes the request for ImportWallet.
type ImportWalletRequest struct {
	Wallet         string `json:"wallet"`
	Passphrase     string `json:"passphrase"`
	RecoveryPhrase string `json:"recoveryPhrase"`
	Version        uint32 `json:"version"`
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

	if len(req.RecoveryPhrase) == 0 {
		errs.AddForProperty("recoveryPhrase", commands.ErrIsRequired)
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

	return req, nil
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

// VerifyAnyResponse describes the response for VerifyAny.
type VerifyAnyResponse struct {
	Valid bool `json:"success"`
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
	ImportWallet(name, passphrase, recoveryPhrase string, version uint32) error
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
	Revoke(token string) (string, error)
}

// NodeForward ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/node_forward_mock.go -package mocks code.vegaprotocol.io/vegawallet/service NodeForward
type NodeForward interface {
	SendTx(context.Context, *commandspb.Transaction, api.SubmitTransactionRequest_Type, int) (string, error)
	HealthCheck(context.Context) error
	LastBlockHeightAndHash(context.Context) (*api.LastBlockHeightResponse, int, error)
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

	s.handle(http.MethodPost, "/api/v1/auth/token", s.Login)
	s.handle(http.MethodDelete, "/api/v1/auth/token", extractToken(s.Revoke))

	s.handle(http.MethodGet, "/api/v1/network", s.GetNetwork)

	s.handle(http.MethodPost, "/api/v1/wallets", s.CreateWallet)
	s.handle(http.MethodPost, "/api/v1/wallets/import", s.ImportWallet)

	s.handle(http.MethodGet, "/api/v1/keys", extractToken(s.ListPublicKeys))
	s.handle(http.MethodPost, "/api/v1/keys", extractToken(s.GenerateKeyPair))
	s.handle(http.MethodGet, "/api/v1/keys/:keyid", extractToken(s.GetPublicKey))
	s.handle(http.MethodPut, "/api/v1/keys/:keyid/taint", extractToken(s.TaintKey))
	s.handle(http.MethodPut, "/api/v1/keys/:keyid/metadata", extractToken(s.UpdateMeta))

	s.handle(http.MethodPost, "/api/v1/command", extractToken(s.SignTx))
	s.handle(http.MethodPost, "/api/v1/command/sync", extractToken(s.SignTxSync))
	s.handle(http.MethodPost, "/api/v1/command/commit", extractToken(s.SignTxCommit))
	s.handle(http.MethodPost, "/api/v1/sign", extractToken(s.SignAny))
	s.handle(http.MethodPost, "/api/v1/verify", s.VerifyAny)

	s.handle(http.MethodGet, "/api/v1/version", s.Version)
	s.handle(http.MethodGet, "/api/v1/status", s.Health)

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

	recoveryPhrase, err := s.handler.CreateWallet(req.Wallet, req.Passphrase)
	if err != nil {
		s.writeBadRequestErr(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, CreateWalletResponse{
		RecoveryPhrase: recoveryPhrase,
		Token:          token,
	})
}

func (s *Service) ImportWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseImportWalletRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	err := s.handler.ImportWallet(req.Wallet, req.Passphrase, req.RecoveryPhrase, req.Version)
	if err != nil {
		s.writeBadRequestErr(w, err)
		return
	}

	token, err := s.auth.NewSession(req.Wallet)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, TokenResponse{Token: token})
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
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, TokenResponse{Token: token})
}

func (s *Service) Revoke(t string, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	name, err := s.auth.Revoke(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	s.handler.LogoutWallet(name)

	s.writeSuccess(w, nil)
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
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			s.writeForbiddenError(w, err)
		} else {
			s.writeInternalError(w, err)
		}
		return
	}

	key, err := s.handler.GetPublicKey(name, pubKey)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, KeyResponse{Key: key})
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
			statusCode = http.StatusInternalServerError
		}
		s.writeError(w, newErrorResponse(err.Error()), statusCode)
		return
	}

	s.writeSuccess(w, KeyResponse{Key: key})
}

func (s *Service) ListPublicKeys(t string, w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	name, err := s.auth.VerifyToken(t)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	keys, err := s.handler.ListPublicKeys(name)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, KeysResponse{Keys: keys})
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

	if err = s.handler.TaintKey(name, keyID, req.Passphrase); err != nil {
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			s.writeForbiddenError(w, err)
		} else {
			s.writeInternalError(w, err)
		}
		return
	}

	s.writeSuccess(w, nil)
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

	if err = s.handler.UpdateMeta(name, keyID, req.Passphrase, req.Meta); err != nil {
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			s.writeForbiddenError(w, err)
		} else {
			s.writeInternalError(w, err)
		}
		return
	}

	s.writeSuccess(w, nil)
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
		s.writeInternalError(w, err)
		return
	}

	res := SignAnyResponse{
		HexSignature:    hex.EncodeToString(signature),
		Base64Signature: base64.StdEncoding.EncodeToString(signature),
	}

	s.writeSuccess(w, res)
}

func (s *Service) VerifyAny(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, errs := ParseVerifyAnyRequest(r)
	if !errs.Empty() {
		s.writeBadRequest(w, errs)
		return
	}

	verified, err := s.handler.VerifyAny(req.decodedInputData, req.decodedSignature, req.PubKey)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	s.writeSuccess(w, VerifyAnyResponse{Valid: verified})
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

	blockData, cltIdx, err := s.nodeForward.LastBlockHeightAndHash(r.Context())
	if err != nil {
		s.writeInternalError(w, ErrCouldNotGetBlockHeight)
		return
	}

	name, err := s.auth.VerifyToken(token)
	if err != nil {
		s.writeForbiddenError(w, err)
		return
	}

	tx, err := s.handler.SignTx(name, req, blockData.Height)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}

	if !req.Propagate {
		s.writeSuccess(w, nil)
		return
	}

	// generate proof of work for the transaction
	tid := vgcrypto.RandomHash()
	powNonce, _, err := vgcrypto.PoW(blockData.Hash, tid, uint(blockData.SpamPowDifficulty), vgcrypto.Sha3)
	if err != nil {
		s.writeInternalError(w, err)
		return
	}
	tx.Pow = &commandspb.ProofOfWork{
		Tid:   tid,
		Nonce: powNonce,
	}

	txHash, err := s.nodeForward.SendTx(r.Context(), tx, ty, cltIdx)
	if err != nil {
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

	s.writeSuccess(w, struct {
		TxHash string                  `json:"txHash"`
		Tx     *commandspb.Transaction `json:"tx"`
	}{
		TxHash: txHash,
		Tx:     tx,
	})
}

func (s *Service) Version(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := VersionResponse{
		Version:     version.Version,
		VersionHash: version.VersionHash,
	}

	s.writeSuccess(w, res)
}

func (s *Service) GetNetwork(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	res := NetworkResponse{
		Network: *s.network,
	}
	s.writeSuccess(w, res)
}

func (s *Service) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := s.nodeForward.HealthCheck(r.Context()); err != nil {
		s.writeError(w, newErrorResponse(err.Error()), http.StatusFailedDependency)
		return
	}
	s.writeSuccess(w, nil)
}

func (s *Service) writeBadRequestErr(w http.ResponseWriter, err error) {
	errs := commands.NewErrors()
	s.writeErrors(w, http.StatusBadRequest, errs.FinalAdd(err))
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
	s.log.Info(fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)))
}

func unmarshalBody(r *http.Request, into interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrCouldNotReadRequest
	}
	if len(body) == 0 {
		return nil
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
		return
	}
	s.log.Info(fmt.Sprintf("%d %s", status, http.StatusText(status)))
}

func (s *Service) writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if data == nil {
		s.log.Info(fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK)))
		return
	}

	buf, err := json.Marshal(data)
	if err != nil {
		s.log.Error("couldn't marshal error", zap.Error(err))
		s.writeInternalError(w, fmt.Errorf("couldn't marshal error: %w", err))
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		s.log.Error("couldn't write error to HTTP response", zap.Error(err))
		s.writeInternalError(w, fmt.Errorf("couldn't write error to HTTP response: %w", err))
		return
	}
	s.log.Info(fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK)))
}

func (s *Service) handle(method string, path string, handle httprouter.Handle) {
	loggedEndpoint := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		s.log.Info(fmt.Sprintf("--> %s %s", method, path))
		handle(w, r, p)
		s.log.Info(fmt.Sprintf("<-- %s %s", method, path))
	}
	s.Handle(method, path, loggedEndpoint)
}
