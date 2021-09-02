package v1

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/service"
	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/shared/paths"
)

type Store struct {
	pubRsaKeyFilePath  string
	privRsaKeyFilePath string
	configFilePath     string
}

func InitialiseStore(p paths.Paths) (*Store, error) {
	serviceConfigFilePath, err := p.ConfigPathFor(paths.WalletServiceDefaultConfigFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get data path for %s: %w", paths.WalletServiceDefaultConfigFile, err)
	}

	pubRsaKeyFilePath, err := p.DataPathFor(paths.WalletServicePublicRSAKeyDataFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get data path for %s: %w", paths.WalletServicePublicRSAKeyDataFile, err)
	}

	privRsaKeyFilePath, err := p.DataPathFor(paths.WalletServicePrivateRSAKeyDataFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get data path for %s: %w", paths.WalletServicePrivateRSAKeyDataFile, err)
	}

	return &Store{
		pubRsaKeyFilePath:  pubRsaKeyFilePath,
		privRsaKeyFilePath: privRsaKeyFilePath,
		configFilePath:     serviceConfigFilePath,
	}, nil
}

func (s *Store) ConfigExists() (bool, error) {
	return vgfs.FileExists(s.configFilePath)
}

func (s *Store) GetConfigPath() string {
	return s.configFilePath
}

func (s *Store) GetConfig() (*service.Config, error) {
	cfg := service.NewDefaultConfig()

	err := paths.ReadStructuredFile(s.configFilePath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (s *Store) SaveConfig(cfg *service.Config) error {
	return paths.WriteStructuredFile(s.configFilePath, cfg)
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

func (s *Store) SaveRSAKeys(keys *service.RSAKeys) error {
	if err := vgfs.WriteFile(s.privRsaKeyFilePath, keys.Priv); err != nil {
		return fmt.Errorf("unable to save private key: %w", err)
	}

	if err := vgfs.WriteFile(s.pubRsaKeyFilePath, keys.Pub); err != nil {
		return fmt.Errorf("unable to save public key: %w", err)
	}

	return nil
}

func (s *Store) GetRsaKeys() (*service.RSAKeys, error) {
	pub, err := vgfs.ReadFile(s.pubRsaKeyFilePath)
	if err != nil {
		return nil, err
	}

	priv, err := vgfs.ReadFile(s.privRsaKeyFilePath)
	if err != nil {
		return nil, err
	}

	return &service.RSAKeys{
		Pub:  pub,
		Priv: priv,
	}, nil
}
