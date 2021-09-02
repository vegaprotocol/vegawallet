package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// File structure for configuration
//
// XDG_CONFIG_HOME
// └── vega
// 		├── console/
// 		│	├── config.toml
// 		│	└── proxy.toml
// 		├── data-node/
// 		│	└── config.toml
// 		├── node/
// 		│	├── config.toml
// 		│	└── wallets.toml
// 		├── wallet-cli/
// 		│	└── config.toml
// 		├── wallet-desktop/
// 		│	└── config.toml
// 		└── wallet-service/
// 			└── config.toml

var (
	// VegaConfigHome is the root folder containing all the configuration files.
	// Configuration files should be scoped inside sub-folder base on their
	// context.
	VegaConfigHome = "vega"

	// ConsoleConfigHome is the folder containing the configuration files
	// dedicated to the console.
	ConsoleConfigHome = filepath.Join(VegaConfigHome, "console")

	// ConsoleDefaultConfigFile is the default configuration file for the
	// console.
	ConsoleDefaultConfigFile = filepath.Join(ConsoleConfigHome, "config.toml")

	// ConsoleProxyConfigFile is the configuration file for the
	// console proxy.
	ConsoleProxyConfigFile = filepath.Join(ConsoleConfigHome, "proxy.toml")

	// DataNodeConfigHome is the folder containing the configuration files
	// dedicated to the node.
	DataNodeConfigHome = filepath.Join(VegaConfigHome, "data-node")

	// DataNodeDefaultConfigFile is the default configuration file for the
	// data-node.
	DataNodeDefaultConfigFile = filepath.Join(DataNodeConfigHome, "config.toml")

	// NodeConfigHome is the folder containing the configuration files dedicated
	// to the node.
	NodeConfigHome = filepath.Join(VegaConfigHome, "node")

	// NodeDefaultConfigFile is the default configuration file for the node.
	NodeDefaultConfigFile = filepath.Join(NodeConfigHome, "config.toml")

	// NodeWalletsConfigFile is the configuration file for the node wallets.
	NodeWalletsConfigFile = filepath.Join(NodeConfigHome, "wallets.encrypted")

	// WalletCLIConfigHome is the folder containing the configuration files
	// dedicated to the wallet CLI.
	WalletCLIConfigHome = filepath.Join(VegaConfigHome, "wallet-cli")

	// WalletCLIDefaultConfigFile is the default configuration file for the
	// wallet CLI.
	WalletCLIDefaultConfigFile = filepath.Join(WalletCLIConfigHome, "config.toml")

	// WalletDesktopConfigHome is the folder containing the configuration files
	// dedicated to the wallet desktop application.
	WalletDesktopConfigHome = filepath.Join(VegaConfigHome, "wallet-desktop")

	// WalletDesktopDefaultConfigFile is the default configuration file for the
	// wallet desktop application.
	WalletDesktopDefaultConfigFile = filepath.Join(WalletDesktopConfigHome, "config.toml")

	// WalletServiceConfigHome is the folder containing the configuration files
	// dedicated to the wallet desktop application.
	WalletServiceConfigHome = filepath.Join(VegaConfigHome, "wallet-service")

	// WalletServiceDefaultConfigFile is the default configuration file for the
	// wallet desktop application.
	WalletServiceDefaultConfigFile = filepath.Join(WalletServiceConfigHome, "config.toml")
)

// DefaultConfigPathFor builds the default path for configuration files and
// creates intermediate directories, if needed.
func DefaultConfigPathFor(relFilePath string) (string, error) {
	path, err := xdg.ConfigFile(relFilePath)
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// CustomConfigPathFor builds the path for configuration files at a given root
// path and creates intermediate directories. It scoped the files under a
// "config" folder, and follow the default structure.
func CustomConfigPathFor(rootPath, relFilePath string) (string, error) {
	path := filepath.Join(rootPath, "config", relFilePath)

	if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}
