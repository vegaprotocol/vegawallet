package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	vgfs "code.vegaprotocol.io/go-wallet/libs/fs"
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
	if err := vgfs.EnsureDir(s.configPath); err != nil {
		return fmt.Errorf("error creating directory %s: %w", s.configPath, err)
	}

	if err := vgfs.EnsureDir(s.keyFolderPath); err != nil {
		return fmt.Errorf("error creating directory %s: %w", s.keyFolderPath, err)
	}

	return nil
}

func (s *Store) ConfigExists() (bool, error) {
	return vgfs.FileExists(s.configFilePath)
}

func (s *Store) GetConfigPath() string {
	return s.configFilePath
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
	confPathExists, _ := vgfs.FileExists(s.configFilePath)

	if confPathExists {
		if overwrite {
			if err := s.removeConfigFile(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("configuration already exists at path: %v", s.configFilePath)
		}
	}

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		return err
	}

	if err := vgfs.WriteFile(s.configFilePath, buf.Bytes()); err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	return nil
}

func (s *Store) RSAKeysExists() (bool, error) {
	privKeyExists, err := vgfs.FileExists(s.privRsaKeyFilePath)
	if err != nil {
		return false, err
	}
	pubKeyExists, err := vgfs.FileExists(s.pubRsaKeyFilePath)
	if err != nil {
		return false, err
	}
	return privKeyExists && pubKeyExists, nil
}

func (s *Store) GetRSAKeysPath() map[string]string {
	return map[string]string{
		"public":  s.pubRsaKeyFilePath,
		"private": s.privRsaKeyFilePath,
	}
}

func (s *Store) SaveRSAKeys(keys *service.RSAKeys, overwrite bool) error {
	if exists, _ := vgfs.PathExists(s.keyFolderPath); !exists {
		return ErrRSAFolderDoesNotExists
	}

	privKeyExists, _ := vgfs.FileExists(s.privRsaKeyFilePath)
	pubKeyExists, _ := vgfs.FileExists(s.pubRsaKeyFilePath)
	if privKeyExists && pubKeyExists {
		if overwrite {
			if err := s.removeExistingRSAKeys(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("RSA keys already exist at path: %vg", s.keyFolderPath)
		}
	}

	if err := vgfs.WriteFile(s.privRsaKeyFilePath, keys.Priv); err != nil {
		return fmt.Errorf("unable to save private key: %w", err)
	}

	if err := vgfs.WriteFile(s.pubRsaKeyFilePath, keys.Pub); err != nil {
		return fmt.Errorf("unable to save public key: %w", err)
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
		return fmt.Errorf("unable to remove configuration: %w", err)
	}

	return nil
}

func (s *Store) removeExistingRSAKeys() error {
	if err := os.RemoveAll(s.keyFolderPath); err != nil {
		return fmt.Errorf("unable to remove RSA keys: %w", err)
	}

	if err := vgfs.EnsureDir(s.keyFolderPath); err != nil {
		return fmt.Errorf("error creating directory %s: %w", s.keyFolderPath, err)
	}

	return nil
}
