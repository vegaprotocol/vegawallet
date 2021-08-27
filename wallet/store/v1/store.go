package v1

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"code.vegaprotocol.io/go-wallet/crypto"
	vgfs "code.vegaprotocol.io/go-wallet/libs/fs"
	"code.vegaprotocol.io/go-wallet/wallet"
)

type Store struct {
	walletsPath string
}

func NewStore(walletsPath string) (*Store, error) {
	return &Store{
		walletsPath: walletsPath,
	}, nil
}

// Initialise creates the folders. It does nothing if a folder already
// exists.
func (s *Store) Initialise() error {
	if err := vgfs.EnsureDir(s.walletsPath); err != nil {
		return fmt.Errorf("error creating directory %s: %v", s.walletsPath, err)
	}

	return nil
}

func (s *Store) WalletExists(name string) bool {
	walletPath := s.walletPath(name)

	exists, _ := vgfs.PathExists(walletPath)
	return exists
}

func (s *Store) ListWallets() ([]string, error) {
	walletsParentDir, walletsDir := filepath.Split(s.walletsPath)
	entries, err := fs.ReadDir(os.DirFS(walletsParentDir), walletsDir)
	if err != nil {
		return nil, err
	}
	wallets := make([]string, len(entries))
	for i, entry := range entries {
		wallets[i] = entry.Name()
	}
	sort.Strings(wallets)
	return wallets, nil
}

func (s *Store) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	walletPath := s.walletPath(name)

	if exists, _ := vgfs.FileExists(walletPath); !exists {
		return nil, wallet.ErrWalletDoesNotExists
	}

	buf, err := fs.ReadFile(os.DirFS(s.walletsPath), name)
	if err != nil {
		return nil, err
	}

	decBuf, err := crypto.Decrypt(buf, passphrase)
	if err != nil {
		return nil, err
	}

	versionedWallet := &struct {
		Version uint32 `json:"version"`
	}{}

	err = json.Unmarshal(decBuf, versionedWallet)
	if err != nil {
		return nil, err
	}

	var w wallet.Wallet
	switch versionedWallet.Version {
	case 0:
		w = &wallet.LegacyWallet{}
		break
	case 1:
		w = &wallet.HDWallet{}
		break
	default:
		return nil, fmt.Errorf("wallet with version %d isn't supported", versionedWallet.Version)
	}

	err = json.Unmarshal(decBuf, w)

	return w, nil
}

func (s *Store) SaveWallet(w wallet.Wallet, passphrase string) error {
	buf, err := json.Marshal(w)
	if err != nil {
		return err
	}

	encBuf, err := crypto.Encrypt(buf, passphrase)
	if err != nil {
		return err
	}

	return vgfs.WriteFile(encBuf, s.walletPath(w.Name()))
}

func (s *Store) GetWalletPath(name string) string {
	return s.walletPath(name)
}

func (s *Store) walletPath(name string) string {
	return filepath.Join(s.walletsPath, name)
}
