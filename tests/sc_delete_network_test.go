package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteNetwork(t *testing.T) {
	// given
	home := t.TempDir()
	networkFile := NewFile(t, home, "my-network-1.toml", FakeNetwork("my-network-1"))

	// when
	importNetworkResp, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp).
		WithName("my-network-1").
		LocatedUnder(home)

	// when
	listNetsResp1, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp1)
	AssertListNetwork(t, listNetsResp1).
		WithNetworks("my-network-1")

	// when
	err = NetworkDelete(t, []string{
		"--home", home,
		"--output", "json",
		"--network", "my-network-1",
	})

	// then
	require.NoError(t, err)

	// when
	listNetsResp2, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp2)
	AssertListNetwork(t, listNetsResp2).
		WithoutNetwork()
}
