package cmd

import (
	"errors"
	"fmt"
	"io"
	"time"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

const (
	DefaultForwarderRetryCount = 5
	ForwarderRequestTimeout    = 5 * time.Second
)

var ErrNetworkDoesNotHaveGRPCHostConfigured = errors.New("network does not have gRPC hosts configured")

func NewCmdSend(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send data to the Vega network",
		Long:  "Send data to the Vega network",
	}

	// create subcommands
	cmd.AddCommand(NewCmdSendCommand(w, rf))
	cmd.AddCommand(NewCmdSendTx(w, rf))
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

	if len(net.API.GRPC.Hosts) == 0 {
		return nil, ErrNetworkDoesNotHaveGRPCHostConfigured
	}

	return net.API.GRPC.Hosts, nil
}
