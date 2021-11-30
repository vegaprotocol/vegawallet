package network

import (
	"fmt"

	"code.vegaprotocol.io/shared/paths"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/network Store
type Store interface {
	NetworkExists(string) (bool, error)
	GetNetwork(string) (*Network, error)
	SaveNetwork(*Network) error
	ListNetworks() ([]string, error)
	GetNetworkPath(string) string
	DeleteNetwork(string) error
}

//go:generate go run github.com/golang/mock/mockgen -destination mocks/reader_mock.go -package mocks code.vegaprotocol.io/vegawallet/network Reader
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

func DescribeNetwork(store Store, req *DescribeNetworkRequest) (*DescribeNetworkResponse, error) {
	resp := &DescribeNetworkResponse{}
	net, err := GetNetwork(store, req.Name)
	if err != nil {
		return nil, err
	}

	resp.Name = net.Name
	resp.TokenExpiry = net.TokenExpiry.String()
	resp.Level = net.Level.String()
	resp.Host = net.Host
	resp.Port = net.Port
	resp.API.GRPCConfig.Hosts = net.API.GRPC.Hosts
	resp.API.GRPCConfig.Retries = net.API.GRPC.Retries
	resp.API.RESTConfig.Hosts = net.API.REST.Hosts
	resp.API.GraphQLConfig.Hosts = net.API.GraphQL.Hosts
	resp.Console.LocalPort = net.Console.LocalPort
	resp.Console.URL = net.Console.URL

	return resp, nil
}

type DescribeNetworkRequest struct {
	Name string
}

type DescribeNetworkResponse struct {
	Name        string `json:"name"`
	Level       string `json:"logLevel"`
	TokenExpiry string `json:"tokenExpiry"`
	Port        int    `json:"port"`
	Host        string `json:"host"`
	API         struct {
		GRPCConfig struct {
			Hosts   []string `json:"hosts"`
			Retries uint64   `json:"retries"`
		} `json:"grpcConfig"`
		RESTConfig struct {
			Hosts []string `json:"hosts"`
		} `json:"restConfig"`
		GraphQLConfig struct {
			Hosts []string `json:"hosts"`
		} `json:"graphQLConfig"`
	} `json:"api"`
	Console struct {
		URL       string `json:"url"`
		LocalPort int    `json:"localPort"`
	}
}

type DeleteNetworkRequest struct {
	Name string
}

func DeleteNetwork(store Store, req *DeleteNetworkRequest) error {
	exists, err := store.NetworkExists(req.Name)
	if err != nil {
		return fmt.Errorf("couldn't verify network existence: %w", err)
	}
	if !exists {
		return NewNetworkDoesNotExistError(req.Name)
	}
	if err = store.DeleteNetwork(req.Name); err != nil {
		return fmt.Errorf("couldn't delete network: %w", err)
	}

	return nil
}
