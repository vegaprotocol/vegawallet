package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/service"
	"github.com/zannen/toml"
)

const (
	configFile     = "wallet-service-config.toml"
	rsaKeyPath     = "wallet_rsa"
	pubRsaKeyName  = "public.pem"
	privRsaKeyName = "private.pem"
)

var (
	ErrRSAFolderDoesNotExists = errors.New("RSA folder does not exist")
)

type Store struct {
	configPath         string
	keyFolderPath      string
	pubRsaKeyFilePath  string
	privRsaKeyFilePath string
	configFilePath     string
}

func NewStore(configPath string) (*Store, error) {
	keyFolderPath := filepath.Join(configPath, rsaKeyPath)

	return &Store{
		configPath:         configPath,
		keyFolderPath:      keyFolderPath,
		pubRsaKeyFilePath:  filepath.Join(keyFolderPath, pubRsaKeyName),
		privRsaKeyFilePath: filepath.Join(keyFolderPath, privRsaKeyName),
		configFilePath:     filepath.Join(configPath, configFile),
	}, nil
}

// Initialise creates the folders. It does nothing if a folder already
// exists.
func (s *Store) Initialise() error {
	err := createFolder(s.configPath)
	if err != nil {
		return err
	}

	err = createFolder(s.keyFolderPath)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetConfig() (*service.Config, error) {
	buf, err := fs.ReadFile(os.DirFS(s.configPath), configFile)
	if err != nil {
		return nil, err
	}

	cfg := service.NewDefaultConfig()

	if _, err := toml.Decode(string(buf), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *Store) SaveConfig(cfg *service.Config, overwrite bool) error {
	confPathExists, _ := fsutil.FileExists(s.configFilePath)

	if confPathExists {
		if overwrite {
			if err := s.removeConfigFile(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("configuration already exists at path: %v", s.configFilePath)
		}
	}

	// write configuration to toml
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		return err
	}

	if err := writeFile(buf.Bytes(), s.configFilePath); err != nil {
		return fmt.Errorf("unable to save configuration: %v", err)
	}

	return nil
}

func (s *Store) SaveRSAKeys(keys *service.RSAKeys, overwrite bool) error {
	if ok, _ := fsutil.PathExists(s.keyFolderPath); !ok {
		return ErrRSAFolderDoesNotExists
	}

	privKeyExists, _ := fsutil.FileExists(s.privRsaKeyFilePath)
	pubKeyExists, _ := fsutil.FileExists(s.pubRsaKeyFilePath)
	if privKeyExists && pubKeyExists {
		if overwrite {
			if err := s.removeExistingRSAKeys(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("RSA keys already exist at path: %v", s.keyFolderPath)
		}
	}

	if err := writeFile(keys.Priv, s.privRsaKeyFilePath); err != nil {
		return fmt.Errorf("unable to save private key: %v", err)
	}

	if err := writeFile(keys.Pub, s.pubRsaKeyFilePath); err != nil {
		return fmt.Errorf("unable to save public key: %v", err)
	}

	return nil
}

func (s *Store) GetRsaKeys() (*service.RSAKeys, error) {
	pub, err := fs.ReadFile(os.DirFS(s.keyFolderPath), pubRsaKeyName)
	if err != nil {
		return nil, err
	}

	priv, err := fs.ReadFile(os.DirFS(s.keyFolderPath), privRsaKeyName)
	if err != nil {
		return nil, err
	}

	return &service.RSAKeys{
		Pub:  pub,
		Priv: priv,
	}, nil
}

func (s *Store) removeConfigFile() error {
	if err := os.Remove(s.configFilePath); err != nil {
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

func writeFile(content []byte, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
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
