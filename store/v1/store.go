package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/config"
	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/zannen/toml"
)

const (
	configFile       = "wallet-service-config.toml"
	rsaKeyPath       = "wallet_rsa"
	pubRsaKeyName    = "public.pem"
	privRsaKeyName   = "private.pem"
	walletBaseFolder = "wallets"
)

var (
	ErrRSAFolderDoesNotExists = errors.New("RSA folder does not exist")
)

type Store struct {
	rootPath           string
	keyFolderPath      string
	pubRsaKeyFileName  string
	privRsaKeyFileName string
	configFileName     string
	walletBasePath     string
}

func NewStore(rootPath string) (*Store, error) {
	keyFolderPath := filepath.Join(rootPath, rsaKeyPath)

	return &Store{
		rootPath:           rootPath,
		keyFolderPath:      keyFolderPath,
		pubRsaKeyFileName:  filepath.Join(keyFolderPath, pubRsaKeyName),
		privRsaKeyFileName: filepath.Join(keyFolderPath, privRsaKeyName),
		configFileName:     filepath.Join(rootPath, configFile),
		walletBasePath:     filepath.Join(rootPath, walletBaseFolder),
	}, nil
}

// Initialise creates the folders. It does nothing if a folder already
// exists.
func (s *Store) Initialise() error {
	err := createFolder(s.rootPath)
	if err != nil {
		return err
	}

	err = createFolder(s.walletBasePath)
	if err != nil {
		return err
	}

	err = createFolder(s.keyFolderPath)
	if err != nil {
		return err
	}

	return nil
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

func (s *Store) WalletExists(name string) bool {
	walletPath := s.walletPath(name)

	ok, _ := fsutil.PathExists(walletPath)
	return ok
}

func (s *Store) GetWallet(name, passphrase string) (wallet.Wallet, error) {
	walletPath := s.walletPath(name)

	if ok, _ := fsutil.PathExists(walletPath); !ok {
		return wallet.Wallet{}, wallet.ErrWalletDoesNotExists
	}

	buf, err := ioutil.ReadFile(walletPath)
	if err != nil {
		return wallet.Wallet{}, err
	}

	decBuf, err := crypto.Decrypt(buf, passphrase)
	if err != nil {
		return wallet.Wallet{}, err
	}

	w := &wallet.Wallet{}
	err = json.Unmarshal(decBuf, w)
	return *w, err
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

	f, err := os.Create(s.walletPath(w.Name))
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

func (s *Store) GetConfig() (*config.Config, error) {
	buf, err := ioutil.ReadFile(s.configFileName)
	if err != nil {
		return nil, err
	}

	cfg := config.NewDefaultConfig()

	if _, err := toml.Decode(string(buf), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *Store) SaveConfig(cfg *config.Config, overwrite bool) error {
	confPathExists, _ := fsutil.FileExists(s.configFileName)

	if confPathExists {
		if overwrite {
			if err := s.removeConfigFile(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("configuration already exists at path: %v", s.configFileName)
		}
	}

	// write configuration to toml
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		return err
	}

	// create the configuration file
	f, err := os.Create(s.configFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(buf.String()); err != nil {
		return err
	}

	return nil
}

func (s *Store) SaveRSAKeys(keys *wallet.RSAKeys, overwrite bool) error {
	if ok, _ := fsutil.PathExists(s.keyFolderPath); !ok {
		return ErrRSAFolderDoesNotExists
	}

	privKeyExists, _ := fsutil.FileExists(s.privRsaKeyFileName)
	pubKeyExists, _ := fsutil.FileExists(s.pubRsaKeyFileName)
	if privKeyExists && pubKeyExists {
		if overwrite {
			if err := s.removeExistingRSAKeys(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("RSA keys already exist at path: %v", s.keyFolderPath)
		}
	}

	if err := writeFile(keys.Priv, s.privRsaKeyFileName); err != nil {
		return fmt.Errorf("unable to write private key: %v", err)
	}

	if err := writeFile(keys.Pub, s.pubRsaKeyFileName); err != nil {
		return fmt.Errorf("unable to write private key: %v", err)
	}

	return nil
}

func (s *Store) GetRsaKeys() (*wallet.RSAKeys, error) {
	pub, err := ioutil.ReadFile(s.pubRsaKeyFileName)
	if err != nil {
		return nil, err
	}

	priv, err := ioutil.ReadFile(s.privRsaKeyFileName)
	if err != nil {
		return nil, err
	}

	return &wallet.RSAKeys{
		Pub:  pub,
		Priv: priv,
	}, nil
}

func (s *Store) removeConfigFile() error {
	if err := os.Remove(s.configFileName); err != nil {
		return fmt.Errorf("unable to remove configuration: %v", err)
	}

	return nil
}

func (s *Store) removeExistingRSAKeys() error {
	if err := os.RemoveAll(s.keyFolderPath); err != nil {
		return fmt.Errorf("unable to remove RSA keys: %v", err)
	}

	return createFolder(s.keyFolderPath)
}

func (s *Store) walletPath(name string) string {
	return filepath.Join(s.rootPath, walletBaseFolder, name)
}

func writeFile(content []byte, fileName string) error {
	pemFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer pemFile.Close()

	_, err = pemFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}
