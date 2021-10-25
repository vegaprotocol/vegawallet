package v1

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/shared/paths"

	"code.vegaprotocol.io/vegawallet/network"
)

const fileExt = ".toml"

type Store struct {
	networksHome string
}

func InitialiseStore(vegaPaths paths.Paths) (*Store, error) {
	networksHome, err := vegaPaths.CreateConfigDirFor(paths.WalletServiceNetworksConfigHome)
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
		networks[i] = s.fileNameToName(entry.Name())
	}
	sort.Strings(networks)
	return networks, nil
}

func (s *Store) GetNetworksPath() string {
	return s.networksHome
}

func (s *Store) GetNetworkPath(name string) string {
	return s.nameToFilePath(name)
}

func (s *Store) NetworkExists(name string) (bool, error) {
	return vgfs.FileExists(s.GetNetworkPath(name))
}

func (s *Store) GetNetwork(name string) (*network.Network, error) {
	net := &network.Network{}
	if err := paths.ReadStructuredFile(s.nameToFilePath(name), &net); err != nil {
		return nil, fmt.Errorf("couldn't read network configuration file: %w", err)
	}
	if name != net.Name {
		return nil, fmt.Errorf("network configuration file name (%s) and network name (%s) don't match", name, net.Name)
	}
	migrateNetwork(net)
	return net, nil
}

func (s *Store) SaveNetwork(net *network.Network) error {
	migrateNetwork(net)
	if err := paths.WriteStructuredFile(s.nameToFilePath(net.Name), net); err != nil {
		return fmt.Errorf("couldn't write network configuration file: %w", err)
	}
	return nil
}

func (s *Store) nameToFilePath(network string) string {
	return filepath.Join(s.networksHome, network+fileExt)
}

func (s *Store) fileNameToName(fileName string) string {
	return fileName[:len(fileName)-len(fileExt)]
}

// migrateNetwork ensures the legacy configuration is migrated to the new
// one.
// TO REMOVE Once the tools use the new API.GRPC.
func migrateNetwork(net *network.Network) {
	if len(net.Nodes.Hosts) > 0 && len(net.API.GRPC.Hosts) == 0 {
		net.API.GRPC = net.Nodes
	}
	// void this legacy property to avoid confusion
	net.Nodes = network.GRPCConfig{}
}
