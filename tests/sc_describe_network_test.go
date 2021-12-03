package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDescribeNetwork(t *testing.T) {
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
	describeResp, err := NetworkDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--network", "my-network-1",
	})

	// then
	require.NoError(t, err)
	AssertDescribeNetwork(t, describeResp).
		WithName("my-network-1").
		WithHostAndPort("127.0.0.1", 8000).
		WithTokenExpiry("1h0m0s").
		WithConsole("console.example.com", 1847).
		WithGRPCConfig([]string{"example.com:3007"}, 5).
		WithRESTConfig([]string{"https://example.com/rest"}).
		WithGraphQLConfig([]string{"https://example.com/gql/query"})

	// when
	describeResp, err = NetworkDescribe(t, []string{
		"--home", home,
		"--output", "json",
		"--from-file", networkFile,
		"--network", "i-do-not-exist",
	})

	// then
	require.Error(t, err)
	require.Nil(t, describeResp)
}
