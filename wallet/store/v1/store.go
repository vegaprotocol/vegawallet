package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"code.vegaprotocol.io/go-wallet/wallet"
	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	vgfs "code.vegaprotocol.io/shared/libs/fs"
)

var (
	ErrWrongPassphrase = errors.New("wrong passphrase")
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
		return fmt.Errorf("error creating directory %s: %w", s.walletsPath, err)
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
		return nil, fmt.Errorf("couldn't read directory at %s: %w", s.walletsPath, err)
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

	if exists, err := vgfs.FileExists(walletPath); !exists {
		return nil, fmt.Errorf("couldn't verify file presence at %s: %w", walletPath, err)
	}

	buf, err := fs.ReadFile(os.DirFS(s.walletsPath), name)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file at %s: %w", s.walletsPath, err)
	}

	decBuf, err := vgcrypto.Decrypt(buf, passphrase)
	if err != nil {
		if err.Error() == "cipher: message authentication failed" {
			return nil, ErrWrongPassphrase
		}
		return nil, err
	}

	versionedWallet := &struct {
		Version uint32 `json:"version"`
	}{}

	err = json.Unmarshal(decBuf, versionedWallet)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal wallet verion: %w", err)
	}

	var w wallet.Wallet
	switch versionedWallet.Version {
	case 1:
		w = &wallet.HDWallet{}
		break
	default:
		return nil, fmt.Errorf("wallet with version %d isn't supported", versionedWallet.Version)
	}

	err = json.Unmarshal(decBuf, w)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal wallet: %w", err)
	}

	return w, nil
}

func (s *Store) SaveWallet(w wallet.Wallet, passphrase string) error {
	buf, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("couldn't marshal wallet: %w", err)
	}

	encBuf, err := vgcrypto.Encrypt(buf, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't encrypt wallet: %w", err)
	}

	walletPath := s.walletPath(w.Name())
	err = vgfs.WriteFile(walletPath, encBuf)
	if err != nil {
		return fmt.Errorf("couldn't write wallet file at %s: %w", walletPath, err)
	}
	return nil
}

func (s *Store) GetWalletPath(name string) string {
	return s.walletPath(name)
}

func (s *Store) walletPath(name string) string {
	return filepath.Join(s.walletsPath, name)
}
