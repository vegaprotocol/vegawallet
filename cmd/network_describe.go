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
	describeNetworkLong = cli.LongDesc(`
	    Describe all known information about the specified network.
	`)

	describeNetworkExample = cli.Examples(`
		# Describe a network
		vegawallet network describe --network NETWORK
	`)
)

type DescribeNetworkHandler func(*network.DescribeNetworkRequest) (*network.DescribeNetworkResponse, error)

func NewCmdDescribeNetwork(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *network.DescribeNetworkRequest) (*network.DescribeNetworkResponse, error) {
		vegaPaths := paths.New(rf.Home)

		netStore, err := netstore.InitialiseStore(vegaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return network.DescribeNetwork(netStore, req)
	}

	return BuildCmdDescribeNetwork(w, h, rf)
}

type DescribeNetworkFlags struct {
	Network string
	Output  string
}

func (f *DescribeNetworkFlags) Validate() (*network.DescribeNetworkRequest, error) {
	req := &network.DescribeNetworkRequest{}

	if err := flags.ValidateOutput(f.Output); err != nil {
		return nil, err
	}

	if len(f.Network) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("network")
	}
	req.Name = f.Network

	return req, nil
}

func BuildCmdDescribeNetwork(w io.Writer, handler DescribeNetworkHandler, rf *RootFlags) *cobra.Command {
	f := &DescribeNetworkFlags{}
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe the specified network",
		Long:    describeNetworkLong,
		Example: describeNetworkExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}
			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch f.Output {
			case flags.InteractiveOutput:
				PrintDescribeNetworkResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network to describe",
	)

	addOutputFlag(cmd, &f.Output)

	autoCompleteNetwork(cmd, rf.Home)

	return cmd
}

func PrintDescribeNetworkResponse(w io.Writer, resp *network.DescribeNetworkResponse) {
	p := printer.NewInteractivePrinter(w)
	p.NextLine().Text("Network").NextLine()
	p.Text("  Name:         ").WarningText(resp.Name).NextLine()
	p.Text("  Address:      ").WarningText(resp.Host).WarningText(":").WarningText(fmt.Sprint(resp.Port)).NextLine()
	p.Text("  Token expiry: ").WarningText(resp.TokenExpiry).NextLine()
	p.Text("  Level:        ").WarningText(resp.Level)
	p.NextSection()

	p.Text("API.GRPC").NextLine()
	p.Text("  Retries: ").WarningText(fmt.Sprint(resp.API.GRPCConfig.Retries)).NextLine()
	p.Text("  Hosts:").NextLine()
	for _, h := range resp.API.GRPCConfig.Hosts {
		p.Text("    - ").WarningText(h).NextLine()
	}
	p.NextLine()

	p.Text("API.REST").NextLine()
	p.Text("  Hosts:").NextLine()
	for _, h := range resp.API.RESTConfig.Hosts {
		p.Text("    - ").WarningText(h).NextLine()
	}

	p.NextLine()
	p.Text("API.GraphQL").NextLine()
	p.Text("  Hosts:").NextLine()
	for _, h := range resp.API.GraphQLConfig.Hosts {
		p.Text("    - ").WarningText(h).NextLine()
	}
	p.NextLine()

	p.Text("Console").NextLine()
	p.Text("  Address: ").WarningText(resp.Console.URL).WarningText(":").WarningText(fmt.Sprint(resp.Console.LocalPort))
	p.NextSection()
}
