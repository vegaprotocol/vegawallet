package wallet

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"code.vegaprotocol.io/go-wallet/commands"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
	walletpb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/wallet/v1"

	"github.com/golang/protobuf/proto"
)

var (
	ErrPubKeyIsTainted = errors.New("public key is tainted")
)

// Store abstracts the underlying storage for wallet data.
type Store interface {
	WalletExists(name string) bool
	SaveWallet(w Wallet, passphrase string) error
	GetWallet(name, passphrase string) (Wallet, error)
	GetWalletPath(name string) string
}

type Handler struct {
	store         Store
	loggedWallets wallets

	// just to make sure we do not access same file concurrently or the map
	mu sync.RWMutex
}

func NewHandler(store Store) *Handler {
	return &Handler{
		store:         store,
		loggedWallets: newWallets(),
	}
}

func (h *Handler) WalletExists(name string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.store.WalletExists(name)
}

func (h *Handler) CreateWallet(name, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.store.GetWallet(name, passphrase)
	if err != nil && err != ErrWalletDoesNotExists {
		return err
	} else if err == nil {
		return ErrWalletAlreadyExists
	}

	w := NewWallet(name)

	return h.saveWallet(*w, passphrase)
}

func (h *Handler) LoginWallet(name, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return ErrWalletDoesNotExists
	}

	h.loggedWallets.Add(w)

	return nil
}

func (h *Handler) LogoutWallet(name string) {
	h.loggedWallets.Remove(name)
}

func (h *Handler) GenerateKeyPair(name, passphrase string) (KeyPair, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return KeyPair{}, err
	}

	kp, err := GenKeyPair(crypto.Ed25519)
	if err != nil {
		return KeyPair{}, err
	}

	w.KeyRing.Upsert(*kp)

	err = h.saveWallet(w, passphrase)
	if err != nil {
		return KeyPair{}, err
	}

	return kp.DeepCopy(), nil
}

func (h *Handler) SecureGenerateKeyPair(name, passphrase string) (string, error) {
	kp, err := h.GenerateKeyPair(name, passphrase)
	if err != nil {
		return "", err
	}

	return kp.Pub, nil
}

func (h *Handler) GetPublicKey(name, pubKey string) (*PublicKey, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	kp, err := h.getKeyPair(name, pubKey)
	if err != nil {
		return nil, err
	}

	return kp.ToPublicKey(), nil
}

func (h *Handler) ListPublicKeys(name string) ([]PublicKey, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.loggedWallets.Get(name)
	if err != nil {
		return nil, err
	}

	return w.KeyRing.GetPublicKeys(), nil
}

func (h *Handler) SignAny(name, inputData, pubKey string) ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rawInputData, err := base64.StdEncoding.DecodeString(inputData)
	if err != nil {
		return nil, err
	}

	kp, err := h.getKeyPair(name, pubKey)
	if err != nil {
		return nil, err
	}

	if kp.Tainted {
		return nil, ErrPubKeyIsTainted
	}

	return kp.Algorithm.Sign(kp.privBytes, rawInputData)
}

func (h *Handler) SignTxV2(name string, req *walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	keyPair, err := h.getKeyPair(name, req.GetPubKey())
	if err != nil {
		return nil, err
	}
	if keyPair.Tainted {
		return nil, ErrPubKeyIsTainted
	}

	data := commands.NewInputData(height)
	wrapRequestCommandIntoInputData(data, req)
	marshalledData, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}

	signature, err := keyPair.Sign(marshalledData)
	if err != nil {
		return nil, err
	}

	return commands.NewTransaction(keyPair.pubBytes, marshalledData, signature), nil
}

func (h *Handler) VerifyAny(name, inputData, sig, pubKey string) (bool, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rawSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false, err
	}

	rawInputData, err := base64.StdEncoding.DecodeString(inputData)
	if err != nil {
		return false, err
	}

	kp, err := h.getKeyPair(name, pubKey)
	if err != nil {
		return false, err
	}

	return kp.Algorithm.Verify(kp.pubBytes, rawInputData, rawSig)
}

func (h *Handler) TaintKey(name, pubKey, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return err
	}

	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	if err := keyPair.Taint(); err != nil {
		return err
	}

	w.KeyRing.Upsert(keyPair)

	return h.saveWallet(w, passphrase)
}

func (h *Handler) UpdateMeta(name, pubKey, passphrase string, meta []Meta) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return err
	}

	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	keyPair.Meta = meta

	w.KeyRing.Upsert(keyPair)

	return h.saveWallet(w, passphrase)
}

func (h *Handler) GetWalletPath(name string) (string, error) {
	return h.store.GetWalletPath(name), nil
}

func (h *Handler) getKeyPair(name, pubKey string) (*KeyPair, error) {
	wallet, err := h.loggedWallets.Get(name)
	if err != nil {
		return nil, err
	}

	keyPair, err := wallet.KeyRing.FindPair(pubKey)

	return &keyPair, err
}

func (h *Handler) saveWallet(w Wallet, passphrase string) error {
	err := h.store.SaveWallet(w, passphrase)
	if err != nil {
		return err
	}

	h.loggedWallets.Add(w)

	return nil
}

func wrapRequestCommandIntoInputData(data *commandspb.InputData, req *walletpb.SubmitTransactionRequest) {
	switch cmd := req.Command.(type) {
	case *walletpb.SubmitTransactionRequest_OrderSubmission:
		data.Command = &commandspb.InputData_OrderSubmission{
			OrderSubmission: req.GetOrderSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_OrderCancellation:
		data.Command = &commandspb.InputData_OrderCancellation{
			OrderCancellation: req.GetOrderCancellation(),
		}
	case *walletpb.SubmitTransactionRequest_OrderAmendment:
		data.Command = &commandspb.InputData_OrderAmendment{
			OrderAmendment: req.GetOrderAmendment(),
		}
	case *walletpb.SubmitTransactionRequest_VoteSubmission:
		data.Command = &commandspb.InputData_VoteSubmission{
			VoteSubmission: req.GetVoteSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_WithdrawSubmission:
		data.Command = &commandspb.InputData_WithdrawSubmission{
			WithdrawSubmission: req.GetWithdrawSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission:
		data.Command = &commandspb.InputData_LiquidityProvisionSubmission{
			LiquidityProvisionSubmission: req.GetLiquidityProvisionSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_ProposalSubmission:
		data.Command = &commandspb.InputData_ProposalSubmission{
			ProposalSubmission: req.GetProposalSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_NodeRegistration:
		data.Command = &commandspb.InputData_NodeRegistration{
			NodeRegistration: req.GetNodeRegistration(),
		}
	case *walletpb.SubmitTransactionRequest_NodeVote:
		data.Command = &commandspb.InputData_NodeVote{
			NodeVote: req.GetNodeVote(),
		}
	case *walletpb.SubmitTransactionRequest_NodeSignature:
		data.Command = &commandspb.InputData_NodeSignature{
			NodeSignature: req.GetNodeSignature(),
		}
	case *walletpb.SubmitTransactionRequest_ChainEvent:
		data.Command = &commandspb.InputData_ChainEvent{
			ChainEvent: req.GetChainEvent(),
		}
	case *walletpb.SubmitTransactionRequest_OracleDataSubmission:
		data.Command = &commandspb.InputData_OracleDataSubmission{
			OracleDataSubmission: req.GetOracleDataSubmission(),
		}
	default:
		panic(fmt.Errorf("command %v is not supported", cmd))
	}
}

type wallets map[string]Wallet

func newWallets() wallets {
	return map[string]Wallet{}
}

func (w wallets) Add(wallet Wallet) {
	w[wallet.Name] = wallet
}

func (w wallets) Get(name string) (Wallet, error) {
	wallet, ok := w[name]
	if !ok {
		return Wallet{}, ErrWalletDoesNotExists
	}
	return wallet, nil
}

func (w wallets) Remove(name string) {
	delete(w, name)
}
