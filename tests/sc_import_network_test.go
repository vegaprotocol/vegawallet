package tests_test

import (
	"fmt"
	"sort"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"github.com/stretchr/testify/require"
)

func TestImportNetwork(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	networkFile1 := NewFile(t, home, "my-network-1.toml", fakeNetwork("my-network-1"))

	// when
	importNetworkResp1, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile1,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp1).
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
	require.Equal(t, []string{"my-network-1"}, listNetsResp1.Networks)

	// given
	networkFile2 := NewFile(t, home, "my-network-2.toml", fakeNetwork("my-network-2"))

	// when
	importNetworkResp2, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile2,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp2).
		WithName("my-network-2").
		LocatedUnder(home)

	// when
	listNetsResp2, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp2)
	require.Equal(t, []string{"my-network-1", "my-network-2"}, listNetsResp2.Networks)
}

func TestForceImportNetwork(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	networkFile := NewFile(t, home, "my-network.toml", fakeNetwork("my-network"))

	// when
	importNetworkResp1, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp1).
		WithName("my-network").
		LocatedUnder(home)

	// when
	importNetworkResp2, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
	})

	// then
	require.Error(t, err)
	require.Nil(t, importNetworkResp2)

	// when
	importNetworkResp3, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
		"--force",
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp3).
		WithName("my-network").
		LocatedUnder(home)

	// when
	listNetsResp, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp)
	require.Equal(t, []string{"my-network"}, listNetsResp.Networks)
}

func TestImportNetworkWithNewName(t *testing.T) {
	// given
	home, cleanUpFn := NewTempDir(t)
	defer cleanUpFn(t)

	networkFile := NewFile(t, home, "my-network.toml", fakeNetwork("my-network"))

	// when
	importNetworkResp1, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp1).
		WithName("my-network").
		LocatedUnder(home)

	// when
	listNetsResp1, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp1)
	require.Equal(t, []string{"my-network"}, listNetsResp1.Networks)

	// given
	networkName := vgrand.RandomStr(5)

	// when
	importNetworkResp2, err := NetworkImport(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
		"--with-name", networkName,
	})

	// then
	require.NoError(t, err)
	AssertImportNetwork(t, importNetworkResp2).
		WithName(networkName).
		LocatedUnder(home)

	// when
	listNetsResp2, err := NetworkList(t, []string{
		"--home", home,
		"--output", "json",
	})

	// then
	require.NoError(t, err)
	require.NotNil(t, listNetsResp2)
	expectedNets := []string{"my-network", networkName}
	sort.Strings(expectedNets)
	require.Equal(t, expectedNets, listNetsResp2.Networks)
}

func fakeNetwork(name string) string {
	return fmt.Sprintf(`
Name = "%s"
Level = "info"
TokenExpiry = "1h0m0s"
Port = 8000
Host = "127.0.0.1"

[API.GRPC]
Retries = 5
Hosts = [
    "example.com:3007",
]

[API.REST]
Hosts = [
    "https://example.com/rest"
]

[API.GraphQL]
Hosts = [
    "https://example.com/gql/query"
]

[Console]
URL = "console.example.com"
LocalPort = 1847
`, name)
}
