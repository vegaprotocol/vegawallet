package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocateNetworks(t *testing.T) {
	// given
	home := RandomPath()

	// when
	locateNetworkResp, err := NetworkLocate(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	AssertLocateNetwork(t, locateNetworkResp).
		LocatedUnder(home)
}
