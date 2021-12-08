package wallets

import (
	"errors"
	"fmt"
	"sync"

	"code.vegaprotocol.io/protos/commands"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	wcommands "code.vegaprotocol.io/vegawallet/commands"
	wcrypto "code.vegaprotocol.io/vegawallet/crypto"
	"code.vegaprotocol.io/vegawallet/wallet"
)

var ErrWalletDoesNotExists = errors.New("wallet does not exist")

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

	w, recoveryPhrase, err := wallet.NewHDWallet(name)
	if err != nil {
		return "", err
	}

	err = h.saveWallet(w, passphrase)
	if err != nil {
		return "", err
	}

	return recoveryPhrase, nil
}

func (h *Handler) ImportWallet(name, passphrase, recoveryPhrase string, version uint32) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.store.WalletExists(name) {
		return wallet.ErrWalletAlreadyExists
	}

	w, err := wallet.ImportHDWallet(name, recoveryPhrase, version)
	if err != nil {
		return err
	}

	return h.saveWallet(w, passphrase)
}

func (h *Handler) LoginWallet(name, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.store.WalletExists(name) {
		return ErrWalletDoesNotExists
	}

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			return err
		}
		return fmt.Errorf("couldn't get wallet %s: %w", name, err)
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
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			return nil, err
		}
		return nil, fmt.Errorf("couldn't get wallet %s: %w", name, err)
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

	data, err := wcommands.ToMarshaledInputData(req, height)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal input data: %w", err)
	}

	pubKey := req.GetPubKey()
	signature, err := w.SignTx(pubKey, data)
	if err != nil {
		return nil, fmt.Errorf("couldn't sign transaction: %w", err)
	}

	protoSignature := &commandspb.Signature{
		Value:   signature.Value,
		Algo:    signature.Algo,
		Version: signature.Version,
	}
	return commands.NewTransaction(pubKey, data, protoSignature), nil
}

func (h *Handler) VerifyAny(inputData, sig []byte, pubKey string) (bool, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return wcrypto.VerifyMessage(&wcrypto.VerifyMessageRequest{
		Message:   inputData,
		Signature: sig,
		PubKey:    pubKey,
	})
}

func (h *Handler) TaintKey(name, pubKey, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	w, err := h.store.GetWallet(name, passphrase)
	if err != nil {
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			return err
		}
		return fmt.Errorf("couldn't get wallet %s: %w", name, err)
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
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			return err
		}
		return fmt.Errorf("couldn't get wallet %s: %w", name, err)
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
		if errors.Is(err, wallet.ErrWrongPassphrase) {
			return err
		}
		return fmt.Errorf("couldn't get wallet %s: %w", name, err)
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

func (h *Handler) getLoggedWallet(name string) (wallet.Wallet, error) {
	if exists := h.store.WalletExists(name); !exists {
		return nil, ErrWalletDoesNotExists
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
		nextID := len(w.ListKeyPairs()) + 1

		meta = append(meta, wallet.Meta{
			Key:   "name",
			Value: fmt.Sprintf("%s key %d", w.Name(), nextID),
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
