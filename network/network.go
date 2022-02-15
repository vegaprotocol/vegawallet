package network

import (
	"errors"

	"code.vegaprotocol.io/vegawallet/service/encoding"
)

var (
	ErrNetworkDoesNotHaveGRPCHostConfigured              = errors.New("network configuration does not have any gRPC host set")
	ErrNetworkDoesNotHaveHostConfiguredForConsole        = errors.New("network configuration does not have any host set for console")
	ErrNetworkDoesNotHaveLocalPortConfiguredForConsole   = errors.New("network configuration does not have any local port set for console")
	ErrNetworkDoesNotHaveHostConfiguredForTokenDApp      = errors.New("network configuration does not have any host set for token dApp")
	ErrNetworkDoesNotHaveLocalPortConfiguredForTokenDApp = errors.New("network configuration does not have any local port set for token dApp")
)

type Network struct {
	Name        string            `json:"name"`
	Level       encoding.LogLevel `json:"level"`
	TokenExpiry encoding.Duration `json:"tokenExpiry"`
	Port        int               `json:"port"`
	Host        string            `json:"host"`
	API         APIConfig         `json:"api"`
	TokenDApp   TokenDAppConfig   `json:"tokenDApp"`
	Console     ConsoleConfig     `json:"console"`
}

type APIConfig struct {
	GRPC    GRPCConfig    `json:"grpc"`
	REST    RESTConfig    `json:"rest"`
	GraphQL GraphQLConfig `json:"graphQl"`
}

type GRPCConfig struct {
	Hosts   []string `json:"hosts"`
	Retries uint64   `json:"retries"`
}

type RESTConfig struct {
	Hosts []string `json:"hosts"`
}

type GraphQLConfig struct {
	Hosts []string `json:"hosts"`
}

type ConsoleConfig struct {
	URL       string `json:"url"`
	LocalPort int    `json:"localPort"`
}

type TokenDAppConfig struct {
	URL       string `json:"url"`
	LocalPort int    `json:"localPort"`
}

func (n *Network) EnsureCanConnectGRPCNode() error {
	if len(n.API.GRPC.Hosts) > 0 && len(n.API.GRPC.Hosts[0]) > 0 {
		return nil
	}
	return ErrNetworkDoesNotHaveGRPCHostConfigured
}

func (n *Network) EnsureCanConnectConsole() error {
	if len(n.Console.URL) == 0 {
		return ErrNetworkDoesNotHaveHostConfiguredForConsole
	}
	if n.Console.LocalPort == 0 {
		return ErrNetworkDoesNotHaveLocalPortConfiguredForConsole
	}
	return nil
}

func (n *Network) EnsureCanConnectTokenDApp() error {
	if len(n.TokenDApp.URL) == 0 {
		return ErrNetworkDoesNotHaveHostConfiguredForTokenDApp
	}
	if n.TokenDApp.LocalPort == 0 {
		return ErrNetworkDoesNotHaveLocalPortConfiguredForTokenDApp
	}
	return nil
}
