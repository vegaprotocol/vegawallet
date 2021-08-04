package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blang/semver/v4"
)

const (
	ReleasesPage       = "https://github.com/vegaprotocol/go-wallet/releases"
	ReleaseAPI         = "https://api.github.com/repos/vegaprotocol/go-wallet/releases"
	TaggedReleaseURL   = "https://github.com/vegaprotocol/go-wallet/releases/tag"
	defaultVersionHash = "unknown"
	defaultVersion     = "unknown"
)

var (
	// VersionHash specifies the git commit used to build the application.
	// See VERSION_HASH in Makefile for details.
	VersionHash = defaultVersionHash

	// Version specifies the version used to build the application.
	// See VERSION in Makefile for details.
	Version = defaultVersion
)

// Check returns a newer version, or an error or nil for both
// if no error happened, and no updates are needed
func Check(currentVersion string) (*semver.Version, error) {
	req, err := http.NewRequest("GET", ReleaseAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	releases := []struct {
		Name string `json:"name"`
	}{}

	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, err
	}

	last, _ := semver.Make(strings.TrimPrefix(currentVersion, "v"))
	cur := last

	for _, v := range releases {
		newV, err := semver.Make(strings.TrimPrefix(v.Name, "v"))
		if err != nil {
			// unsupported version
			continue
		}
		if newV.GT(last) {
			last = newV
		}
	}

	if cur.EQ(last) {
		// no updates
		return nil, nil
	}

	return &last, nil
}

func GetReleaseURL(v *semver.Version) string {
	return fmt.Sprintf("%v/v%v", TaggedReleaseURL, v)
}
