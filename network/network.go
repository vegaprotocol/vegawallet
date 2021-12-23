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
	Name        string
	Level       encoding.LogLevel
	TokenExpiry encoding.Duration
	Port        int
	Host        string
	API         APIConfig
	TokenDApp   TokenDAppConfig
	Console     ConsoleConfig
}

type APIConfig struct {
	GRPC    GRPCConfig
	REST    RESTConfig
	GraphQL GraphQLConfig
}

type GRPCConfig struct {
	Hosts   []string
	Retries uint64
}

type RESTConfig struct {
	Hosts []string
}

type GraphQLConfig struct {
	Hosts []string
}

type ConsoleConfig struct {
	URL       string
	LocalPort int
}

type TokenDAppConfig struct {
	URL       string
	LocalPort int
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
