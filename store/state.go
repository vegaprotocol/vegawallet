package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// File structure for state
//
// XDG_STATE_HOME
// └── vega
// 		├── node/
// 		│	├── logs/
// 		│	├── checkpoints/
// 		│	└── snapshots/
// 		├── wallet-cli/
// 		│	└── logs/
// 		├── wallet-desktop/
// 		│	└── logs/
// 		└── wallet-service/
// 			└── logs/

var (
	// VegaStateHome is the root folder containing all the state related to Vega.
	VegaStateHome = "vega"

	// DataNodeStateHome is the folder containing the state dedicated to the
	// data-node.
	DataNodeStateHome = filepath.Join(VegaStateHome, "data-node")

	// DataNodeLogsHome is the folder containing the logs of the data-node.
	DataNodeLogsHome = filepath.Join(DataNodeStateHome, "logs")

	// NodeStateHome is the folder containing the state of the node.
	NodeStateHome = filepath.Join(VegaStateHome, "node")

	// NodeLogsHome is the folder containing the logs of the node.
	NodeLogsHome = filepath.Join(NodeStateHome, "logs")

	// CheckpointStateHome is the folder containing the checkpoint files
	// of to the node.
	CheckpointStateHome = filepath.Join(NodeStateHome, "checkpoints")

	// SnapshotStateHome is the folder containing the snapshot files
	// of to the node.
	SnapshotStateHome = filepath.Join(NodeStateHome, "snapshots")

	// WalletCLIStateHome is the folder containing the state of the wallet CLI.
	WalletCLIStateHome = filepath.Join(VegaStateHome, "wallet-cli")

	// WalletCLILogsHome is the folder containing the logs of the wallet CLI.
	WalletCLILogsHome = filepath.Join(WalletCLIStateHome, "logs")

	// WalletDesktopStateHome is the folder containing the state of the wallet
	// desktop.
	WalletDesktopStateHome = filepath.Join(VegaStateHome, "wallet-desktop")

	// WalletDesktopLogsHome is the folder containing the logs of the wallet
	// desktop.
	WalletDesktopLogsHome = filepath.Join(WalletDesktopStateHome, "logs")

	// WalletServiceStateHome is the folder containing the state of the node.
	WalletServiceStateHome = filepath.Join(VegaStateHome, "wallet-service")

	// WalletServiceLogsHome is the folder containing the logs of the node.
	WalletServiceLogsHome = filepath.Join(WalletServiceStateHome, "logs")
)

// DefaultStatePathFor builds the default path for state files and creates
// intermediate directories, if needed.
func DefaultStatePathFor(relFilePath string) (string, error) {
	path, err := xdg.StateFile(relFilePath)
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// CustomStatePathFor builds the path for cache files at a given root path and
// creates intermediate directories. It scoped the files under a "cache" folder,
// and follow the default structure.
func CustomStatePathFor(rootPath, relFilePath string) (string, error) {
	path := filepath.Join(rootPath, "state", relFilePath)

	if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}
