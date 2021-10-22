package v1

import (
	"fmt"

	"code.vegaprotocol.io/vegawallet/service"
	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/shared/paths"
)

type Store struct {
	pubRsaKeyFilePath  string
	privRsaKeyFilePath string
}

func InitialiseStore(p paths.Paths) (*Store, error) {
	pubRsaKeyFilePath, err := p.CreateDataPathFor(paths.WalletServicePublicRSAKeyDataFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get data path for %s: %w", paths.WalletServicePublicRSAKeyDataFile, err)
	}

	privRsaKeyFilePath, err := p.CreateDataPathFor(paths.WalletServicePrivateRSAKeyDataFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get data path for %s: %w", paths.WalletServicePrivateRSAKeyDataFile, err)
	}

	return &Store{
		pubRsaKeyFilePath:  pubRsaKeyFilePath,
		privRsaKeyFilePath: privRsaKeyFilePath,
	}, nil
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

func (s *Store) GetRSAKeysPath() (string, string) {
	return s.pubRsaKeyFilePath, s.privRsaKeyFilePath
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
