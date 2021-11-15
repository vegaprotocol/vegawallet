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
	listNetworkLong = cli.LongDesc(`
		List all registered networks.
	`)

	listNetworkExample = cli.Examples(`
		# List networks
		vegawallet network list"
	`)
)

type ListNetworksHandler func() (*network.ListNetworksResponse, error)

func NewCmdListNetworks(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func() (*network.ListNetworksResponse, error) {
		vegaPaths := paths.New(rf.Home)

		netStore, err := netstore.InitialiseStore(vegaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return network.ListNetworks(netStore)
	}

	return BuildCmdListNetworks(w, h, rf)
}

func BuildCmdListNetworks(w io.Writer, handler ListNetworksHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all registered networks",
		Long:    listNetworkLong,
		Example: listNetworkExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp, err := handler()
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintListNetworksResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintListNetworksResponse(w io.Writer, resp *network.ListNetworksResponse) {
	p := printer.NewInteractivePrinter(w)

	if len(resp.Networks) == 0 {
		p.InfoText("No network registered").NextLine()
		return
	}

	for _, net := range resp.Networks {
		p.Text(fmt.Sprintf("- %s", net)).NextLine()
	}
}
