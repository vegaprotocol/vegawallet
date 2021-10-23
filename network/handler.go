package network

import (
	"embed"
	"fmt"

	"github.com/zannen/toml"
)

var (
	//go:embed defaults/*.toml
	defaultNetworks embed.FS
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/network Store
type Store interface {
	NetworkExists(string) (bool, error)
	GetNetwork(string) (*Network, error)
	SaveNetwork(*Network) error
}

func InitialiseNetworks(store Store, overwrite bool) error {
	entries, err := defaultNetworks.ReadDir("defaults")
	if err != nil {
		return fmt.Errorf("couldn't read defaults directory: %w", err)
	}

	for _, entry := range entries {
		data, err := defaultNetworks.ReadFile(fmt.Sprintf("defaults/%s", entry.Name()))
		if err != nil {
			return fmt.Errorf("couldn't read file: %w", err)
		}
		net := &Network{}
		if _, err := toml.Decode(string(data), &net); err != nil {
			return fmt.Errorf("couldn't decode embeded data: %w", err)
		}

		if !overwrite {
			exists, err := store.NetworkExists(net.Name)
			if err != nil {
				return fmt.Errorf("couldn't verify network existance: %w", err)
			}
			if exists {
				return fmt.Errorf("network %s already exists", net.Name)
			}
		}

		if err = store.SaveNetwork(net); err != nil {
			return fmt.Errorf("couldn't save network configuration: %w", err)
		}
	}

	return nil
}

func ImportNetwork(store Store, net *Network, overwrite bool) error {
	exists, err := store.NetworkExists(net.Name)
	if err != nil {
		return fmt.Errorf("couldn't verify network existance: %w", err)
	}
	if exists && !overwrite {
		return fmt.Errorf("network \"%s\" already exists", net.Name)
	}

	if err := store.SaveNetwork(net); err != nil {
		return fmt.Errorf("couldn't save the imported network: %w", err)
	}

	return nil
}
