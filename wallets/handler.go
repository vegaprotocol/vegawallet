package wallets

import (
	"encoding/hex"
	"fmt"
	"sync"

	wcrypto "code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/protos/commands"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"

	"github.com/golang/protobuf/proto"
)

// Store abstracts the underlying storage for wallet data.
type Store interface {
	WalletExists(name string) bool
	SaveWallet(w wallet.Wallet, passphrase string) error
	GetWallet(name, passphrase string) (wallet.Wallet, error)
	GetWalletPath(name string) string
	ListWallets() ([]string, error)
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

func (h *Handler) ListWallets() ([]string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.store.ListWallets()
}

func (h *Handler) CreateWallet(name, passphrase string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.store.WalletExists(name) {
		return "", wallet.ErrWalletAlreadyExists
	}

	w, mnemonic, err := wallet.NewHDWallet(name)
	if err != nil {
		return "", err
	}

	err = h.saveWallet(w, passphrase)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (h *Handler) ImportWallet(name, passphrase, mnemonic string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.store.WalletExists(name) {
		return wallet.ErrWalletAlreadyExists
	}

	w, err := wallet.ImportHDWallet(name, mnemonic)
	if err != nil {
		return err
	}

	return h.saveWallet(w, passphrase)
}

func (h *Handler) LoginWallet(name, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return wallet.ErrWalletDoesNotExists
	}

	h.loggedWallets.Add(w)

	return nil
}

func (h *Handler) LogoutWallet(name string) {
	h.loggedWallets.Remove(name)
}

func (h *Handler) GenerateKeyPair(name, passphrase string, meta []wallet.Meta) (wallet.KeyPair, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return nil, err
	}

	meta = addDefaultAlias(meta, w)

	kp, err := w.GenerateKeyPair(meta)
	if err != nil {
		return nil, err
	}

	err = h.saveWallet(w, passphrase)
	if err != nil {
		return nil, err
	}

	return kp, nil
}

func (h *Handler) SecureGenerateKeyPair(name, passphrase string, meta []wallet.Meta) (string, error) {
	kp, err := h.GenerateKeyPair(name, passphrase, meta)
	if err != nil {
		return "", err
	}

	return kp.PublicKey(), nil
}

func (h *Handler) GetPublicKey(name, pubKey string) (wallet.PublicKey, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.getLoggedWallet(name)
	if err != nil {
		return nil, err
	}

	return w.DescribePublicKey(pubKey)
}

func (h *Handler) ListPublicKeys(name string) ([]wallet.PublicKey, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.getLoggedWallet(name)
	if err != nil {
		return nil, err
	}

	return w.ListPublicKeys(), nil
}

func (h *Handler) ListKeyPairs(name string) ([]wallet.KeyPair, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.getLoggedWallet(name)
	if err != nil {
		return nil, err
	}

	return w.ListKeyPairs(), nil
}

func (h *Handler) SignAny(name string, inputData []byte, pubKey string) ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.getLoggedWallet(name)
	if err != nil {
		return nil, err
	}

	return w.SignAny(pubKey, inputData)
}

func (h *Handler) SignTx(name string, req *walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w, err := h.getLoggedWallet(name)
	if err != nil {
		return nil, err
	}

	data := commands.NewInputData(height)
	WrapRequestCommandIntoInputData(data, req)
	marshalledData, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}

	pubKey := req.GetPubKey()
	signature, err := w.SignTx(pubKey, marshalledData)
	if err != nil {
		return nil, err
	}

	protoSignature := &commandspb.Signature{
		Value:   signature.Value,
		Algo:    signature.Algo,
		Version: signature.Version,
	}
	return commands.NewTransaction(pubKey, marshalledData, protoSignature), nil
}

func (h *Handler) VerifyAny(inputData, sig []byte, pubKey string) (bool, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	decodedPubKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, err
	}

	signatureAlgorithm := wcrypto.NewEd25519()
	return signatureAlgorithm.Verify(decodedPubKey, inputData, sig)
}

func (h *Handler) TaintKey(name, pubKey, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return err
	}

	err = w.TaintKey(pubKey)
	if err != nil {
		return err
	}

	return h.saveWallet(w, passphrase)
}

func (h *Handler) UntaintKey(name string, pubKey string, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return err
	}

	err = w.UntaintKey(pubKey)
	if err != nil {
		return err
	}

	return h.saveWallet(w, passphrase)
}

func (h *Handler) UpdateMeta(name, pubKey, passphrase string, meta []wallet.Meta) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		return err
	}

	err = w.UpdateMeta(pubKey, meta)
	if err != nil {
		return err
	}

	return h.saveWallet(w, passphrase)
}

func (h *Handler) GetWalletPath(name string) (string, error) {
	return h.store.GetWalletPath(name), nil
}

func (h *Handler) saveWallet(w wallet.Wallet, passphrase string) error {
	err := h.store.SaveWallet(w, passphrase)
	if err != nil {
		return err
	}

	h.loggedWallets.Add(w)

	return nil
}

func WrapRequestCommandIntoInputData(data *commandspb.InputData, req *walletpb.SubmitTransactionRequest) {
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
	case *walletpb.SubmitTransactionRequest_DelegateSubmission:
		data.Command = &commandspb.InputData_DelegateSubmission{
			DelegateSubmission: req.GetDelegateSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_UndelegateSubmission:
		data.Command = &commandspb.InputData_UndelegateSubmission{
			UndelegateSubmission: req.GetUndelegateSubmission(),
		}
	default:
		panic(fmt.Errorf("command %v is not supported", cmd))
	}
}

func (h *Handler) getLoggedWallet(name string) (wallet.Wallet, error) {
	exists := h.store.WalletExists(name)
	if !exists {
		return nil, wallet.ErrWalletDoesNotExists
	}

	w, loggedIn := h.loggedWallets.Get(name)
	if !loggedIn {
		return nil, wallet.ErrWalletNotLoggedIn
	}
	return w, nil
}

func addDefaultAlias(meta []wallet.Meta, w wallet.Wallet) []wallet.Meta {
	hasName := false
	for _, m := range meta {
		if m.Key == "name" {
			hasName = true
		}
	}
	if !hasName {
		nextId := len(w.ListKeyPairs()) + 1

		meta = append(meta, wallet.Meta{
			Key:   "name",
			Value: fmt.Sprintf("%s key %d", w.Name(), nextId),
		})
	}
	return meta
}

type wallets map[string]wallet.Wallet

func newWallets() wallets {
	return map[string]wallet.Wallet{}
}

func (w wallets) Add(wallet wallet.Wallet) {
	w[wallet.Name()] = wallet
}

func (w wallets) Get(name string) (wallet.Wallet, bool) {
	wal, ok := w[name]
	return wal, ok
}

func (w wallets) Remove(name string) {
	delete(w, name)
}
