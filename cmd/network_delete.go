package cmd

import (
	"fmt"
	"io"

	vgterm "code.vegaprotocol.io/shared/libs/term"
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
		# Delete the specified network
		vegawallet network delete --network NETWORK

		# Delete the specified network without asking for confirmation
		vegawallet delete --wallet WALLET --force
	`)
)

type DeleteNetworkHandler func(*network.DeleteNetworkRequest) error

func NewCmdDeleteNetwork(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *network.DeleteNetworkRequest) error {
		vegaPaths := paths.New(rf.Home)

		netStore, err := netstore.InitialiseStore(vegaPaths)
		if err != nil {
			return fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return network.DeleteNetwork(netStore, req)
	}

	return BuildCmdDeleteNetwork(w, h, rf)
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

			if !f.Force && vgterm.HasTTY() {
				confirm, err := flags.DoYouConfirm()
				if err != nil {
					return err
				}
				if !confirm {
					return nil
				}
			}

			if err = handler(req); err != nil {
				return err
			}

			if rf.Output == flags.InteractiveOutput {
				PrintDeleteNetworkResponse(w, f.Network)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network to delete",
	)
	cmd.Flags().BoolVarP(&f.Force,
		"force", "f",
		false,
		"Do not ask for confirmation",
	)

	autoCompleteNetwork(cmd, rf.Home)

	return cmd
}

type DeleteNetworkFlags struct {
	Network string
	Force   bool
}

func (f *DeleteNetworkFlags) Validate() (*network.DeleteNetworkRequest, error) {
	req := &network.DeleteNetworkRequest{}

	if len(f.Network) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("network")
	}
	req.Name = f.Network

	return req, nil
}

func PrintDeleteNetworkResponse(w io.Writer, networkName string) {
	p := printer.NewInteractivePrinter(w)
	p.CheckMark().SuccessText("Network ").SuccessBold(networkName).SuccessText(" deleted").NextLine()
}
