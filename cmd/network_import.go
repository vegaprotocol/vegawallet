package cmd

import (
	"errors"
	"fmt"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

var (
	ErrOnlySingleSourceMustBeSpecified = errors.New("only a single source must be specified")
	ErrSourceMustBeSpecified           = errors.New("a source must be specified")

	networkImportArgs struct {
		filePath string
		url      string
		name     string
		force    bool
	}

	// networkImportCmd represents the network import command.
	networkImportCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a network configuration",
		Long:  "Import a network configuration",
		RunE:  runNetworkImport,
	}
)

func init() {
	networkCmd.AddCommand(networkImportCmd)
	networkImportCmd.Flags().StringVar(&networkImportArgs.filePath, "from-file", "", `Path of the file containing the network configuration to import`)
	networkImportCmd.Flags().StringVar(&networkImportArgs.url, "from-url", "", `URL of the file containing the network configuration to import`)
	networkImportCmd.Flags().StringVar(&networkImportArgs.name, "with-name", "", `Change the name of the imported network`)
	networkImportCmd.Flags().BoolVarP(&networkImportArgs.force, "force", "f", false, "Overwrite the existing network if it has the same name")
}

func runNetworkImport(_ *cobra.Command, _ []string) error {
	vegaPaths := paths.New(rootArgs.home)

	netStore, err := netstore.InitialiseStore(vegaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise networks store: %w", err)
	}

	net, err := GetNetworkFromSource()
	if err != nil {
		return err
	}

	if len(networkImportArgs.name) != 0 {
		net.Name = networkImportArgs.name
	}

	if err := network.ImportNetwork(netStore, net, networkImportArgs.force); err != nil {
		return fmt.Errorf("couldn't import network: %w", err)
	}

	filePath := netStore.GetNetworkPath(net.Name)
	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Importing the network succeeded").NJump(2)
		p.Text("Name:").Jump().WarningText(net.Name).Jump()
		p.Text("File path:").Jump().WarningText(filePath).Jump()
	} else if rootArgs.output == "json" {
		return vgjson.Print(struct {
			FilePath string `json:"filePath"`
		}{
			FilePath: filePath,
		})
	}

	return nil
}

func GetNetworkFromSource() (*network.Network, error) {
	net := &network.Network{}

	if len(networkImportArgs.filePath) != 0 && len(networkImportArgs.url) != 0 {
		return nil, ErrOnlySingleSourceMustBeSpecified
	}

	if len(networkImportArgs.filePath) != 0 {
		if err := paths.ReadStructuredFile(networkImportArgs.filePath, net); err != nil {
			return nil, fmt.Errorf("couldn't read file from %s: %w", networkImportArgs.filePath, err)
		}
	} else if len(networkImportArgs.url) != 0 {
		if err := paths.FetchStructuredFile(networkImportArgs.url, net); err != nil {
			return nil, fmt.Errorf("couldn't fetch file from %s: %w", networkImportArgs.url, err)
		}
	} else {
		return nil, ErrSourceMustBeSpecified
	}

	return net, nil
}
