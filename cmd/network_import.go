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
	importNetworkLong = cli.LongDesc(`
		Import a network configuration from a file or an URL.
	`)

	importNetworkExample = cli.Examples(`
		# import a network configuration from a file
		vegawallet network import --from-file PATH_TO_NETWORK

		# import a network configuration from an URL
		vegawallet network import --from-url URL_TO_NETWORK

		# overwrite existing network configuration
		vegawallet network import --from-url URL_TO_NETWORK --force

		# import a network configuration with a different name
		vegawallet network import --from-url URL_TO_NETWORK --with-name NEW_NAME
	`)
)

type ImportNetworkFromSourceHandler func(*network.ImportNetworkFromSourceRequest) (*network.ImportNetworkFromSourceResponse, error)

func NewCmdImportNetwork(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *network.ImportNetworkFromSourceRequest) (*network.ImportNetworkFromSourceResponse, error) {
		vegaPaths := paths.New(rf.Home)

		s, err := netstore.InitialiseStore(vegaPaths)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise networks store: %w", err)
		}

		return network.ImportNetworkFromSource(s, network.NewReaders(), req)
	}

	return BuildCmdImportNetwork(w, h)
}

func BuildCmdImportNetwork(w io.Writer, handler ImportNetworkFromSourceHandler) *cobra.Command {
	f := &ImportNetworkFlags{}

	cmd := &cobra.Command{
		Use:     "import",
		Short:   "Import a network configuration",
		Long:    importNetworkLong,
		Example: importNetworkExample,
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
				PrintImportNetworkResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&f.FilePath,
		"from-file",
		"",
		"Path to the file containing the network configuration to import",
	)
	cmd.Flags().StringVar(&f.URL,
		"from-url",
		"",
		"URL of the file containing the network configuration to import",
	)
	cmd.Flags().StringVar(&f.Name,
		"with-name",
		"",
		"Change the name of the imported network",
	)
	cmd.Flags().BoolVarP(&f.Force,
		"force", "f",
		false,
		"Overwrite the existing network if it has the same name",
	)

	addOutputFlag(cmd, &f.Output)

	return cmd
}

type ImportNetworkFlags struct {
	FilePath string
	URL      string
	Name     string
	Force    bool
	Output   string
}

func (f *ImportNetworkFlags) Validate() (*network.ImportNetworkFromSourceRequest, error) {
	if err := flags.ValidateOutput(f.Output); err != nil {
		return nil, err
	}

	if len(f.FilePath) == 0 && len(f.URL) == 0 {
		return nil, flags.OneOfFlagsMustBeSpecifiedError("from-file", "from-url")
	}

	if len(f.FilePath) != 0 && len(f.URL) != 0 {
		return nil, flags.FlagsMutuallyExclusiveError("from-file", "from-url")
	}

	return &network.ImportNetworkFromSourceRequest{
		FilePath: f.FilePath,
		URL:      f.URL,
		Name:     f.Name,
		Force:    f.Force,
	}, nil
}

func PrintImportNetworkResponse(w io.Writer, resp *network.ImportNetworkFromSourceResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Importing the network succeeded").NextSection()
	p.Text("Name:").NextLine().WarningText(resp.Name).NextLine()
	p.Text("File path:").NextLine().WarningText(resp.FilePath).NextLine()
}
