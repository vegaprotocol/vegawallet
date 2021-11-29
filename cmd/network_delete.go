package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

var (
	deleteNetworkLong = cli.LongDesc(`
	    Delete the specified network
	`)

	deleteNetworkExample = cli.Examples(`
		# Delete a network
		vegawallet network delete --network NETWORK
	`)
)

func NewCmdDeleteNetwork(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *network.DeleteNetworkRequest) (*network.DeleteNetworkResponse, error) {
		vegaPaths := paths.New(rf.Home)

		netStore, err := netstore.InitialiseStore(vegaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return network.DeleteNetwork(netStore, req)
	}

	return BuildCmdDeleteNetwork(w, h, rf)
}

type DeleteNetworkFlags struct {
	Network string
}

func (f *DeleteNetworkFlags) Validate() (*network.DeleteNetworkRequest, error) {
	req := &network.DeleteNetworkRequest{}

	if len(f.Network) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("network")
	}
	req.Name = f.Network

	return req, nil
}

func BuildCmdDeleteNetwork(w io.Writer, handler DeleteNetworkHandler, rf *RootFlags) *cobra.Command {
	f := &DeleteNetworkFlags{}
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete the specified network",
		Long:    deleteNetworkLong,
		Example: deleteNetworkExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}
			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintDeleteNetworkResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network to delete",
	)

	return cmd
}

type DeleteNetworkHandler func(*network.DeleteNetworkRequest) (*network.DeleteNetworkResponse, error)

func PrintDeleteNetworkResponse(w io.Writer, resp *network.DeleteNetworkResponse) {
	p := printer.NewInteractivePrinter(w)
	p.NextLine().Text("Network ")
	p.WarningText(resp.Name)
	p.Text(" deleted")
}
