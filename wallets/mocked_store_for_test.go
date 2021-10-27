package wallets_test

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
)

var errWrongPassphrase = errors.New("wrong passphrase")

type mockedStore struct {
	passphrase string
	wallets    map[string]wallet.Wallet
}

func newMockedStore() *mockedStore {
	return &mockedStore{
		passphrase: "",
		wallets:    map[string]wallet.Wallet{},
	}
}

func (m *mockedStore) WalletExists(name string) bool {
	_, ok := m.wallets[name]
	return ok
}

func (m *mockedStore) ListWallets() ([]string, error) {
	ws := make([]string, 0, len(m.wallets))
	for k := range m.wallets {
		ws = append(ws, k)
	}
	return ws, nil
}

func (m *mockedStore) SaveWallet(w wallet.Wallet, passphrase string) error {
	m.passphrase = passphrase
	m.wallets[w.Name()] = w
	return nil
}

func (m *mockedStore) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	w, ok := m.wallets[name]
	if !ok {
		return nil, wallets.ErrWalletDoesNotExists
	}
	if passphrase != m.passphrase {
		return nil, errWrongPassphrase
	}
	return w, nil
}

func (m *mockedStore) GetWalletPath(name string) string {
	return fmt.Sprintf("some/path/%v", name)
}

func (m *mockedStore) GetKey(name, pubKey string) wallet.PublicKey {
	w, ok := m.wallets[name]
	if !ok {
		panic(fmt.Sprintf("wallet \"%v\" not found", name))
	}
	for _, key := range w.ListPublicKeys() {
		if key.Key() == pubKey {
			return key
		}
	}
	panic(fmt.Sprintf("key \"%v\" not found", pubKey))
}
