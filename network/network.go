package network

import "code.vegaprotocol.io/vegawallet/service/encoding"

type Network struct {
	Name        string
	Level       encoding.LogLevel
	TokenExpiry encoding.Duration
	Port        int
	Host        string
	API         APIConfig
	TokenDApp   TokenDAppConfig
	Console     ConsoleConfig

	// TO REMOVE Once the tools use the new API.GRPC
	Nodes GRPCConfig
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
