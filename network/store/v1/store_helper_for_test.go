package v1_test

import (
	"os"
	"path/filepath"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/shared/paths"
)

type vegaHome struct {
	customPaths *paths.CustomPaths
}

func newVegaHome() vegaHome {
	rootPath := filepath.Join("/tmp", "vega_tests", vgrand.RandomStr(10))

	return vegaHome{
		customPaths: &paths.CustomPaths{CustomHome: rootPath},
	}
}

func (h *vegaHome) Paths() paths.Paths {
	return h.customPaths
}

func (h *vegaHome) Remove() {
	err := os.RemoveAll(h.customPaths.CustomHome)
	if err != nil {
		panic(err)
	}
}

func (h *vegaHome) NetworksHome() string {
	networksHome, err := h.customPaths.ConfigDirFor(paths.WalletServiceNetworksConfigHome)
	if err != nil {
		panic(err)
	}

	return networksHome
}

func (h *vegaHome) NetworkPath(name string) string {
	return filepath.Join(h.NetworksHome(), name + ".toml")
}
