package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const ReleaseAPI = "https://api.github.com/repos/vegaprotocol/vegawallet/releases"

type ReleasesGetter func() ([]string, error)

type releaseResponse struct {
	Name string `json:"name"`
}

func BuildReleasesRequestFromGithub(ctx context.Context) ReleasesGetter {
	return func() ([]string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ReleaseAPI, nil)
		if err != nil {
			return nil, fmt.Errorf("couldn't build API request to get releases: %w", err)
		}
		req.Header.Add("Accept", "application/vnd.github.v3+json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("couldn't successfully deliver API request to get releases: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read releases API response body: %w", err)
		}

		responses := []releaseResponse{}
		if err = json.Unmarshal(body, &responses); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal releases API response body: %w", err)
		}

		releases := make([]string, 0, len(responses))
		for _, response := range responses {
			releases = append(releases, response.Name)
		}

		return releases, nil
	}
}
