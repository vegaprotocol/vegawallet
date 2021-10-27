package version

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
)

const (
	TaggedReleaseURL   = "https://github.com/vegaprotocol/vegawallet/releases/tag"
	defaultVersionHash = "unknown"
	defaultVersion     = "v0.9.2"
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
	v, err := extractVersionFromRelease(Version)
	if err != nil {
		// unsupported version, considered unreleased
		return true
	}

	return v.isDevelopmentVersion
}

// Check returns a newer version, or an error or nil for both
// if no error happened, and no updates are needed.
func Check(releasesGetterFn ReleasesGetter, currentRelease string) (*semver.Version, error) {
	releases, err := releasesGetterFn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get releases: %w", err)
	}

	currentVersion, err := extractVersionFromRelease(currentRelease)
	if err != nil {
		return nil, fmt.Errorf("couldn't extract version from release: %w", err)
	}
	latestVersion := currentVersion

	var updateAvailable bool
	for _, release := range releases {
		comparedVersion, err := extractVersionFromRelease(release)
		if err != nil {
			// unsupported version
			continue
		}

		if shouldUpdate(latestVersion, comparedVersion) {
			updateAvailable = true
			latestVersion = comparedVersion
		}
	}

	if !updateAvailable {
		return nil, nil
	}

	return latestVersion.version, nil
}

func GetReleaseURL(v *semver.Version) string {
	return fmt.Sprintf("%v/v%v", TaggedReleaseURL, v)
}

func shouldUpdate(latestVersion *cachedVersion, comparedVersion *cachedVersion) bool {
	if latestVersion.isStable && !comparedVersion.isStable {
		return false
	}

	if latestVersion.isDevelopmentVersion && nonDevelopmentVersionAvailable(latestVersion, comparedVersion) {
		return true
	}

	return comparedVersion.version.GT(*latestVersion.version)
}

// nonDevelopmentVersionAvailable verifies if the compared version is the
// non-development equivalent of the latest version.
// For example, 0.9.0-pre1 is the non-development version of 0.9.0-pre1+dev.
// In semantic versioning, we don't compare the `build` annotation, so verifying
// equality is safe.
func nonDevelopmentVersionAvailable(latestVersion *cachedVersion, comparedVersion *cachedVersion) bool {
	return comparedVersion.version.EQ(*latestVersion.version)
}

func extractVersionFromRelease(release string) (*cachedVersion, error) {
	version, err := semver.New(strings.TrimPrefix(release, "v"))
	return asCachedVersion(version), err
}

type cachedVersion struct {
	// version is the original version
	version *semver.Version
	// isDevelopmentVersion tells if the version has a `dev` build annotation.
	isDevelopmentVersion bool
	// isStable tells if the version has any pre-release annotations.
	isStable bool
}

func asCachedVersion(v *semver.Version) *cachedVersion {
	lv := &cachedVersion{
		version: v,
	}

	for _, build := range v.Build {
		if build == "dev" {
			lv.isDevelopmentVersion = true
		}
	}

	lv.isStable = !lv.isDevelopmentVersion && len(v.Pre) == 0

	return lv
}
