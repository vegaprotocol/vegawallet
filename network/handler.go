package network

import (
	"embed"
	"fmt"

	"code.vegaprotocol.io/shared/paths"
	"github.com/zannen/toml"
)

//go:embed defaults/*.toml
var defaultNetworks embed.FS

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/network Store
type Store interface {
	NetworkExists(string) (bool, error)
	GetNetwork(string) (*Network, error)
	SaveNetwork(*Network) error
	ListNetworks() ([]string, error)
	GetNetworkPath(string) string
}

type Reader func(uri string, net interface{}) error

type Readers struct {
	ReadFromFile Reader
	ReadFromURL  Reader
}

func NewReaders() Readers {
	return Readers{
		ReadFromFile: paths.ReadStructuredFile,
		ReadFromURL:  paths.FetchStructuredFile,
	}
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
			return fmt.Errorf("couldn't decode embedded data: %w", err)
		}

		if !overwrite {
			exists, err := store.NetworkExists(net.Name)
			if err != nil {
				return fmt.Errorf("couldn't verify network existence: %w", err)
			}
			if exists {
				return NewNetworkAlreadyExistsError(net.Name)
			}
		}

		if err = store.SaveNetwork(net); err != nil {
			return fmt.Errorf("couldn't save network configuration: %w", err)
		}
	}

	return nil
}

func GetNetwork(store Store, name string) (*Network, error) {
	exists, err := store.NetworkExists(name)
	if err != nil {
		return nil, fmt.Errorf("couldn't verify network existence: %w", err)
	}
	if !exists {
		return nil, NewNetworkDoesNotExistError(name)
	}
	n, err := store.GetNetwork(name)
	if err != nil {
		return nil, fmt.Errorf("couldn't get network %s: %w", name, err)
	}

	return n, nil
}

func ImportNetwork(store Store, net *Network, overwrite bool) error {
	exists, err := store.NetworkExists(net.Name)
	if err != nil {
		return fmt.Errorf("couldn't verify network existence: %w", err)
	}
	if exists && !overwrite {
		return NewNetworkAlreadyExistsError(net.Name)
	}

	if err := store.SaveNetwork(net); err != nil {
		return fmt.Errorf("couldn't save the imported network: %w", err)
	}

	return nil
}

type ImportNetworkFromSourceRequest struct {
	FilePath string
	URL      string
	Name     string
	Force    bool
}

type ImportNetworkFromSourceResponse struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

func ImportNetworkFromSource(s Store, rs Readers, req *ImportNetworkFromSourceRequest) (*ImportNetworkFromSourceResponse, error) {
	net := &Network{}

	if len(req.FilePath) != 0 {
		if err := rs.ReadFromFile(req.FilePath, net); err != nil {
			return nil, fmt.Errorf("couldn't read network configuration at %s: %w", req.FilePath, err)
		}
	} else if len(req.URL) != 0 {
		if err := rs.ReadFromURL(req.URL, net); err != nil {
			return nil, fmt.Errorf("couldn't fetch network configuration from %s: %w", req.URL, err)
		}
	}

	if len(req.Name) != 0 {
		net.Name = req.Name
	}

	if err := ImportNetwork(s, net, req.Force); err != nil {
		return nil, fmt.Errorf("couldn't import network: %w", err)
	}

	return &ImportNetworkFromSourceResponse{
		Name:     net.Name,
		FilePath: s.GetNetworkPath(net.Name),
	}, nil
}

func ListNetworks(store Store) (*ListNetworksResponse, error) {
	nets, err := store.ListNetworks()
	if err != nil {
		return nil, fmt.Errorf("couldn't list networks: %w", err)
	}

	return &ListNetworksResponse{
		Networks: nets,
	}, nil
}

type ListNetworksResponse struct {
	Networks []string `json:"networks"`
}
