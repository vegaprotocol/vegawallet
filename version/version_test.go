package version_test

import (
	"fmt"
	"testing"

	"code.vegaprotocol.io/go-wallet/version"
	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckVersion(t *testing.T) {
	t.Run("Checking current version succeeds", testCheckingCurrentVersionSucceeds)
	t.Run("Verifying unreleased version succeeds", testVerifyingUnreleasedVersionSucceeds)
}

func testCheckingCurrentVersionSucceeds(t *testing.T) {
	tcs := []struct {
		name           string
		currentVersion string
		releases       []string
		expectedResult *semver.Version
	}{
		{
			name:           "stable release against older stable releases should not update",
			currentVersion: "v0.4.0",
			releases:       []string{"v0.2.0", "v0.3.0"},
			expectedResult: nil, // means no update
		}, {
			name:           "stable release against newer stable releases should update to latest stable release",
			currentVersion: "v0.1.0",
			releases:       []string{"v0.2.0", "v0.3.0"},
			expectedResult: toSemVer("0.3.0"),
		}, {
			name:           "stable release against same stable release should not update",
			currentVersion: "v0.3.0",
			releases:       []string{"v0.2.0", "v0.3.0"},
			expectedResult: nil, // means no update
		}, {
			name:           "stable against newer pre-release should update to latest stable release",
			currentVersion: "v0.1.0",
			releases:       []string{"v0.2.0", "v0.3.0-alpha"},
			expectedResult: toSemVer("0.2.0"),
		}, {
			name:           "stable release against older pre-release should not update",
			currentVersion: "v0.3.0",
			releases:       []string{"v0.2.0", "v0.3.0-alpha"},
			expectedResult: nil, // means no update
		}, {
			name:           "pre-release against newer stable release should update to latest stable release",
			currentVersion: "v0.3.0-alpha",
			releases:       []string{"v0.2.0", "v0.3.0"},
			expectedResult: toSemVer("0.3.0"),
		}, {
			name:           "pre-release against older stable release should not update",
			currentVersion: "v0.4.0-alpha",
			releases:       []string{"v0.2.0", "v0.3.0"},
			expectedResult: nil, // means no update
		}, {
			name:           "pre-release against newer pre-release should update to latest pre-release",
			currentVersion: "v0.4.0-alpha",
			releases:       []string{"v0.2.0", "v0.3.0", "v0.4.0-alpha", "v0.4.0-beta"},
			expectedResult: toSemVer("0.4.0-beta"),
		}, {
			name:           "pre-release against same pre-release should not update",
			currentVersion: "v0.4.0-alpha",
			releases:       []string{"v0.2.0", "v0.3.0", "v0.4.0-alpha"},
			expectedResult: nil, // means no update
		}, {
			name:           "pre-release against newer pre-release separated by stable releases should update to latest stable release",
			currentVersion: "v0.2.0-alpha",
			releases:       []string{"v0.2.0-alpha", "v0.2.0", "v0.3.0-alpha", "v0.3.0", "v0.3.1", "v0.4.0-alpha"},
			expectedResult: toSemVer("0.3.1"),
		}, {
			name:           "pre-release against newer development pre-release should not update",
			currentVersion: "v0.4.0-alpha",
			releases:       []string{"v0.2.0", "v0.3.0", "v0.4.0-alpha+dev"},
			expectedResult: nil, // means no update
		}, {
			name:           "development pre-release against non-development pre-release should update to non-development pre-release",
			currentVersion: "v0.4.0-alpha+dev",
			releases:       []string{"v0.2.0", "v0.3.0", "v0.4.0-alpha"},
			expectedResult: toSemVer("0.4.0-alpha"),
		}, {
			name:           "development pre-release against newer pre-release should update to latest pre-release",
			currentVersion: "v0.4.0-alpha+dev",
			releases:       []string{"v0.2.0", "v0.3.0", "v0.4.0-alpha", "v0.4.0-beta"},
			expectedResult: toSemVer("0.4.0-beta"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// given
			releasesGetter := func() ([]string, error) {
				return tc.releases, nil
			}

			// when
			newerVersion, err := version.Check(releasesGetter, tc.currentVersion)

			// then
			require.NoError(tt, err)
			assert.Equal(tt, tc.expectedResult, newerVersion)
		})
	}
}

func testVerifyingUnreleasedVersionSucceeds(t *testing.T) {
	tcs := []struct {
		name         string
		version      string
		isUnreleased bool
	}{
		{
			name:         "stable version is not unreleased",
			version:      "0.1.0",
			isUnreleased: false,
		}, {
			name:         "pre-release version is not unreleased",
			version:      "0.1.0-pre",
			isUnreleased: false,
		}, {
			name:         "development build on stable version is unreleased",
			version:      "0.1.0+dev",
			isUnreleased: true,
		}, {
			name:         "development pre-release version is unreleased",
			version:      "0.1.0-alpha+dev",
			isUnreleased: true,
		}, {
			name:         "annotated build on pre-release version is released",
			version:      "0.1.0-alpha+12345",
			isUnreleased: false,
		}, {
			name:         "non-semver version is unreleased",
			version:      "test",
			isUnreleased: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			version.Version = tc.version

			// when
			isUnreleased := version.IsUnreleased()

			// then
			assert.Equal(tt, tc.isUnreleased, isUnreleased)
		})
	}
}

func toSemVer(s string) *semver.Version {
	expectVersion, err := semver.New(s)
	if err != nil {
		panic(fmt.Errorf("couldn't parse the semver: %w", err))
	}
	return expectVersion
}
