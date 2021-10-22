package service

import (
	"errors"
	"fmt"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/service Store
type Store interface {
	RSAKeysExists() (bool, error)
	SaveRSAKeys(*RSAKeys) error
}

func InitialiseService(store Store, overwrite bool) error {
	keys, err := GenerateRSAKeys()
	if err != nil {
		return fmt.Errorf("couldn't generate RSA keys: %w", err)
	}

	if !overwrite {
		rsaKeysExists, err := store.RSAKeysExists()
		if err != nil {
			return fmt.Errorf("couldn't verify RSA keys existance: %w", err)
		}
		if rsaKeysExists {
			return errors.New("RSA keys already exist")
		}
	}

	if err := store.SaveRSAKeys(keys); err != nil {
		return fmt.Errorf("couldn't save RSA keys: %w", err)
	}

	return nil
}
