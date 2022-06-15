package version

import (
	vgversion "code.vegaprotocol.io/shared/libs/version"

	"github.com/blang/semver/v4"
)

const (
	ReleasesAPI        = "https://api.github.com/repos/vegaprotocol/vegawallet/releases"
	ReleasesURL        = "https://github.com/vegaprotocol/vegawallet/releases"
	defaultVersionHash = "unknown"
	defaultVersion     = "v0.16.1"
)

var (
	// VersionHash specifies the git commit used to build the application.
	// See VERSION_HASH in Makefile for details.
	VersionHash = defaultVersionHash

	// Version specifies the version used to build the application.
	// See VERSION in Makefile for details.
	Version = defaultVersion
)

func IsUnreleased() bool {
	return vgversion.IsUnreleased(Version)
}

type GetVersionResponse struct {
	Version string `json:"version"`
	GitHash string `json:"gitHash"`
}

func GetVersionInfo() *GetVersionResponse {
	return &GetVersionResponse{
		Version: Version,
		GitHash: VersionHash,
	}
}

func Check(releasesGetterFn vgversion.ReleasesGetter, currentRelease string) (*semver.Version, error) {
	return vgversion.Check(releasesGetterFn, currentRelease)
}
