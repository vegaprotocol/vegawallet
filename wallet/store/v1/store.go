package v1

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/fsutil"
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
	err := createFolder(s.walletsPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) WalletExists(name string) bool {
	walletPath := s.walletPath(name)

	ok, _ := fsutil.PathExists(walletPath)
	return ok
}

func (s *Store) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	walletPath := s.walletPath(name)

	if ok, _ := fsutil.PathExists(walletPath); !ok {
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

	f, err := os.Create(s.walletPath(w.Name()))
	if err != nil {
		return err
	}

	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	_, err = f.Write(encBuf)
	if err != nil {
		return err
	}

	return f.Close()
}

func (s *Store) GetWalletPath(name string) string {
	return s.walletPath(name)
}

func (s *Store) walletPath(name string) string {
	return filepath.Join(s.walletsPath, name)
}

func createFolder(folder string) error {
	ok, err := fsutil.PathExists(folder)
	if !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return fmt.Errorf("invalid directory path %s: %v", folder, err)
		}

		if err := fsutil.EnsureDir(folder); err != nil {
			return fmt.Errorf("error creating directory %s: %v", folder, err)
		}
	}
	return nil
}
