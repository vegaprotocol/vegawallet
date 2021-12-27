package network_test

import (
	"testing"

	"code.vegaprotocol.io/vegawallet/network"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("Ensure network can connect to a gRPC node fails", testEnsureNetworkCanConnectGRPCNodeFails)
	t.Run("Ensure network can connect to a console fails", testEnsureNetworkCanConnectConsoleFails)
	t.Run("Ensure network can connect to a token dApp fails", testEnsureNetworkCanConnectTokenDAppFails)
}

func testEnsureNetworkCanConnectGRPCNodeFails(t *testing.T) {
	// given
	net := &network.Network{
		API: network.APIConfig{GRPC: network.GRPCConfig{
			Hosts:   nil,
			Retries: 0,
		}},
	}

	// when
	err := net.EnsureCanConnectGRPCNode()

	// then
	require.ErrorIs(t, err, network.ErrNetworkDoesNotHaveGRPCHostConfigured)
}

func testEnsureNetworkCanConnectConsoleFails(t *testing.T) {
	tcs := []struct {
		name    string
		err     error
		network *network.Network
	}{
		{
			name: "without host",
			err:  network.ErrNetworkDoesNotHaveHostConfiguredForConsole,
			network: &network.Network{
				Console: network.ConsoleConfig{
					URL:       "",
					LocalPort: 1234,
				},
			},
		}, {
			name: "without local port",
			err:  network.ErrNetworkDoesNotHaveLocalPortConfiguredForConsole,
			network: &network.Network{
				Console: network.ConsoleConfig{
					URL:       "https://example.com",
					LocalPort: 0,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			err := tc.network.EnsureCanConnectConsole()

			// then
			require.ErrorIs(tt, err, tc.err)
		})
	}
}

func testEnsureNetworkCanConnectTokenDAppFails(t *testing.T) {
	tcs := []struct {
		name    string
		err     error
		network *network.Network
	}{
		{
			name: "without host",
			err:  network.ErrNetworkDoesNotHaveHostConfiguredForTokenDApp,
			network: &network.Network{
				TokenDApp: network.TokenDAppConfig{
					URL:       "",
					LocalPort: 1234,
				},
			},
		}, {
			name: "without host",
			err:  network.ErrNetworkDoesNotHaveLocalPortConfiguredForTokenDApp,
			network: &network.Network{
				TokenDApp: network.TokenDAppConfig{
					URL:       "https://example.com",
					LocalPort: 0,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			err := tc.network.EnsureCanConnectTokenDApp()

			// then
			require.ErrorIs(tt, err, tc.err)
		})
	}
}
