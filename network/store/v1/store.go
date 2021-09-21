package v1

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"code.vegaprotocol.io/go-wallet/network"
	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/shared/paths"
)

type Store struct {
	networksHome string
}

func InitialiseStore(vegaPaths paths.Paths) (*Store, error) {
	networksHome, err := vegaPaths.ConfigDirFor(paths.WalletServiceNetworksConfigHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get config path for %s: %w", paths.WalletServiceNetworksConfigHome, err)
	}

	return &Store{
		networksHome: networksHome,
	}, nil
}

func (s *Store) ListNetworks() ([]string, error) {
	networksParentDir, networksDir := filepath.Split(s.networksHome)
	entries, err := fs.ReadDir(os.DirFS(networksParentDir), networksDir)
	if err != nil {
		return nil, fmt.Errorf("couldn't read directory at %s: %w", s.networksHome, err)
	}
	networks := make([]string, len(entries))
	for i, entry := range entries {
		networks[i] = strings.Trim(entry.Name(), ".toml")
	}
	sort.Strings(networks)
	return networks, nil
}

func (s *Store) GetNetworksPath() string {
	return s.networksHome
}

func (s *Store) NetworkExists(name string) (bool, error) {
	return vgfs.FileExists(filepath.Join(s.networksHome, fmt.Sprintf("%s.toml", name)))
}

func (s *Store) GetNetwork(name string) (*network.Network, error) {
	cfg := &network.Network{}
	filePath := filepath.Join(s.networksHome, fmt.Sprintf("%s.toml", name))
	if err := paths.ReadStructuredFile(filePath, &cfg); err != nil {
		return nil, fmt.Errorf("couldn't read network configuration file: %w", err)
	}
	if name != cfg.Name {
		return nil, fmt.Errorf("network configuration file name (%s) and network name (%s) don't match", name, cfg.Name)
	}
	return cfg, nil
}

func (s *Store) SaveNetwork(cfg *network.Network) error {
	filePath := filepath.Join(s.networksHome, fmt.Sprintf("%s.toml", cfg.Name))
	if err := paths.WriteStructuredFile(filePath, cfg); err != nil {
		return fmt.Errorf("couldn't write network configuration file: %w", err)
	}
	return nil
}
