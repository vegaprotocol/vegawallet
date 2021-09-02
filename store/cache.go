package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// File structure for cache
//
// XDG_CACHE_HOME
// └── vega
// 		└── data-node/

var (
	// VegaCacheHome is the root folder containing all the cache related to Vega.
	VegaCacheHome = "vega"

	// DataNodeCacheHome is the folder containing the data dedicated to the
	// data-node.
	DataNodeCacheHome = filepath.Join(VegaDataHome, "data-node")
)

// DefaultCachePathFor builds the default path for cache files and creates
// intermediate directories, if needed.
func DefaultCachePathFor(relFilePath string) (string, error) {
	path, err := xdg.CacheFile(relFilePath)
	if err != nil {
		return "", fmt.Errorf("couldn't get the default path for %s: %w", relFilePath, err)
	}
	return path, nil
}

// CustomCachePathFor builds the path for cache files at a given root path and
// creates intermediate directories. It scoped the files under a "cache" folder,
// and follow the default structure.
func CustomCachePathFor(rootPath, relFilePath string) (string, error) {
	path := filepath.Join(rootPath, "cache", relFilePath)

	if err := os.MkdirAll(path, os.ModeDir|0700); err != nil {
		return "", fmt.Errorf("couldn't create directories for %s: %w", path, err)
	}
	return path, nil
}
