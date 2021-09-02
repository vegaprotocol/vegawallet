package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// File structure for data
//
// XDG_DATA_HOME
// └── vega
// 		├── node/
// 		│	└── wallets/
// 		│		├── vega/
// 		│		│	└── vega-node-wallet
// 		│		└── ethereum/
// 		│			└── eth-node-wallet
// 		├── wallets/
// 		│	├── vega-wallet-1
// 		│	└── vega-wallet-2
// 		└── wallet-service/
// 			└── rsa-keys/
// 				├── private.pem
// 				└── public.pem

var (
	// VegaDataHome is the root folder containing all the data related to Vega.
	VegaDataHome = "vega"

	// NodeDataHome is the folder containing the data dedicated to the node.
	NodeDataHome = filepath.Join(VegaDataHome, "vega")

	// NodeWalletsDataHome is the folder containing the data dedicated to the
	// node wallets.
	NodeWalletsDataHome = filepath.Join(NodeDataHome, "wallets")

	// VegaNodeWalletsDataHome is the folder containing the vega wallet
	// dedicated to the node.
	VegaNodeWalletsDataHome = filepath.Join(NodeWalletsDataHome, "vega")

	// EthereumNodeWalletsDataHome is the folder containing the ethereum wallet
	// dedicated to the node.
	EthereumNodeWalletsDataHome = filepath.Join(NodeWalletsDataHome, "ethereum")

	// WalletsDataHome is the folder containing the user wallets.
	WalletsDataHome = filepath.Join(VegaDataHome, "wallets")

	// WalletServiceDataHome is the folder containing the data dedicated to the
	// wallet service.
	WalletServiceDataHome = filepath.Join(VegaDataHome, "wallet-service")

	// WalletServiceRSAKeysDataHome is the folder containing the RSA keys used by
	// the wallet service.
	WalletServiceRSAKeysDataHome = filepath.Join(WalletServiceDataHome, "rsa-keys")

	// WalletServicePublicRSAKeyDataFile is the file containing the public RSA key
	// used by the wallet service.
	WalletServicePublicRSAKeyDataFile = filepath.Join(WalletServiceRSAKeysDataHome, "public.pem")

	// WalletServicePrivateRSAKeyDataFile is the file containing the private RSA key
	// used by the wallet service.
	WalletServicePrivateRSAKeyDataFile = filepath.Join(WalletServiceRSAKeysDataHome, "private.pem")
)

// DefaultDataPathFor builds the default path for data files and creates
// intermediate directories, if needed.
func DefaultDataPathFor(relFilePath string) (string, error) {
	path, err := xdg.DataFile(relFilePath)
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// CustomDataPathFor builds the path for data files at a given root path and
// creates intermediate directories. It scoped the files under a "data" folder,
// and follow the default structure.
func CustomDataPathFor(rootPath, relFilePath string) (string, error) {
	path := filepath.Join(rootPath, "data", relFilePath)

	if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}
