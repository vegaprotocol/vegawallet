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

func (h *vegaHome) ConfigFilePath() string {
	serviceConfigFilePath, err := h.customPaths.ConfigPathFor(paths.WalletServiceDefaultConfigFile)
	if err != nil {
		panic(err)
	}

	return serviceConfigFilePath
}

func (h *vegaHome) ConfigHome() string {
	configHome, err := h.customPaths.ConfigPathFor(paths.WalletServiceConfigHome)
	if err != nil {
		panic(err)
	}

	return configHome
}

func (h *vegaHome) RSAKeysHome() string {
	rsaKeyHome, err := h.customPaths.DataPathFor(paths.WalletServiceRSAKeysDataHome)
	if err != nil {
		panic(err)
	}

	return rsaKeyHome
}

func (h *vegaHome) PublicRSAKeyFilePath() string {
	pubRsaKeyFilePath, err := h.customPaths.DataPathFor(paths.WalletServicePublicRSAKeyDataFile)
	if err != nil {
		panic(err)
	}

	return pubRsaKeyFilePath
}

func (h *vegaHome) PrivateRSAKeyFilePath() string {
	privRsaKeyFilePath, err := h.customPaths.DataPathFor(paths.WalletServicePrivateRSAKeyDataFile)
	if err != nil {
		panic(err)
	}

	return privRsaKeyFilePath
}
