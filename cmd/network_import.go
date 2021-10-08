package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/network"
	netstore "code.vegaprotocol.io/go-wallet/network/store/v1"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
	"github.com/spf13/cobra"
)

var (
	networkImportArgs struct {
		filePath string
		name     string
		force    bool
	}

	// networkImportCmd represents the network import command
	networkImportCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a network configuration",
		Long:  "Import a network configuration",
		RunE:  runNetworkImport,
	}
)

func init() {
	networkCmd.AddCommand(networkImportCmd)
	importCmd.Flags().StringVar(&networkImportArgs.filePath, "from-file", "", `Path of the file containing the network configuration to import`)
	importCmd.Flags().StringVar(&networkImportArgs.name, "with-name", "", `Change the name of the imported network`)
	importCmd.Flags().BoolVarP(&networkImportArgs.force, "force", "f", false, "Overwrite the existing network if it has the same name")
}

func runNetworkImport(_ *cobra.Command, _ []string) error {
	vegaPaths := paths.New(rootArgs.home)

	netStore, err := netstore.InitialiseStore(vegaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise networks store: %w", err)
	}

	net := &network.Network{}

	if len(networkImportArgs.filePath) != 0 {
		if err := paths.ReadStructuredFile(networkImportArgs.filePath, net); err != nil {
			return fmt.Errorf("couldn't read file from %s: %w", networkImportArgs.filePath, err)
		}
	} else {
		return errors.New("source must be specified")
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
