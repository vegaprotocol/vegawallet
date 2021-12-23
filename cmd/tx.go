package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

func NewCmdTx(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Provides utilities for interacting with transactions",
		Long:  "Provides utilities for interacting with transactions",
	}

	cmd.AddCommand(NewCmdTxSend(w, rf))
	return cmd
}

func getHostsFromNetwork(rf *RootFlags, networkName string) ([]string, error) {
	netStore, err := netstore.InitialiseStore(paths.New(rf.Home))
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise network store: %w", err)
	}
	net, err := network.GetNetwork(netStore, networkName)
	if err != nil {
		return nil, err
	}

	if err := net.EnsureCanConnectGRPCNode(); err != nil {
		return nil, err
	}

	return net.API.GRPC.Hosts, nil
}
