package service

import (
	"time"

	"code.vegaprotocol.io/go-wallet/service/encoding"

	"go.uber.org/zap"
)

const (
	// 7 days, needs to be in seconds for the token
	tokenExpiry = time.Hour * 24 * 7
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/go-wallet/service Store
type Store interface {
	SaveConfig(*Config, bool) error
	SaveRSAKeys(*RSAKeys, bool) error
}

type Config struct {
	Level       encoding.LogLevel
	TokenExpiry encoding.Duration
	Port        int
	Host        string
	Nodes       NodesConfig
	Console     ConsoleConfig
}

type ConsoleConfig struct {
	URL       string
	LocalPort int
}

type NodesConfig struct {
	Hosts   []string
	Retries uint64
}

func GenerateConfig(log *zap.Logger, store Store, overwrite bool) error {
	config := NewDefaultConfig()
	err := store.SaveConfig(&config, overwrite)
	if err != nil {
		return err
	}

	log.Info("wallet service configuration generated successfully")

	keys, err := GenerateRSAKeys()
	if err != nil {
		return err
	}
	if err := store.SaveRSAKeys(keys, overwrite); err != nil {
		return err
	}

	log.Info("wallet RSA keys generated successfully")

	return nil
}

// NewDefaultConfig creates an instance of the package specific configuration,
// given a pointer to a logger instance to be used for logging within the
// package.
func NewDefaultConfig() Config {
	return Config{
		Level:       encoding.LogLevel{Level: zap.InfoLevel},
		TokenExpiry: encoding.Duration{Duration: tokenExpiry},
		Nodes: NodesConfig{
			Hosts: []string{
				"n01.testnet.vega.xyz:3002",
				"n02.testnet.vega.xyz:3002",
				"n03.testnet.vega.xyz:3002",
				"n04.testnet.vega.xyz:3002",
				"n05.testnet.vega.xyz:3002",
				"n06.testnet.vega.xyz:3002",
				"n07.testnet.vega.xyz:3002",
				"n08.testnet.vega.xyz:3002",
				"n09.testnet.vega.xyz:3002",
				"n10.testnet.vega.xyz:3002",
				"n11.testnet.vega.xyz:3002",
				"n12.testnet.vega.xyz:3002",
				"n13.testnet.vega.xyz:3002",
				"n14.testnet.vega.xyz:3002",
				"n15.testnet.vega.xyz:3002",
				"n16.testnet.vega.xyz:3002",
				"n17.testnet.vega.xyz:3002",
				"n18.testnet.vega.xyz:3002",
				"n19.testnet.vega.xyz:3002",
				"n20.testnet.vega.xyz:3002",
			},
			Retries: 5,
		},
		Host: "127.0.0.1",
		Port: 1789,
		Console: ConsoleConfig{
			URL:       "console.fairground.wtf",
			LocalPort: 1847,
		},
	}
}
