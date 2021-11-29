package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteNetwork(t *testing.T) {
	// given
	home := t.TempDir()

	networkFile := NewFile(t, home, "my-network-1.toml", FakeNetwork("my-network-1"))

	// when (import network)
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

	// when (list networks)
	listNetsResp1, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp1)
	require.Equal(t, []string{"my-network-1"}, listNetsResp1.Networks)

	// when (delete the network)
	deleteNetworkResp, err := NetworkDelete(t, []string{
		"--home", home,
		"--output", "json",
		"--network", "my-network-1",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, deleteNetworkResp)
	require.Equal(t, "my-network-1", deleteNetworkResp.Name)

	// when (list networks again)
	listNetsResp2, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp2)
	require.Equal(t, []string{}, listNetsResp2.Networks)
}
