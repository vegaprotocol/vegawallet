package cmd

import (
	"fmt"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

// networkListCmd represents the network list command.
var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered network",
	Long:  "List all registered network",
	RunE:  runNetworkList,
}

func init() {
	networkCmd.AddCommand(networkListCmd)
}

func runNetworkList(_ *cobra.Command, _ []string) error {
	vegaPaths := paths.New(rootArgs.home)

	netStore, err := netstore.InitialiseStore(vegaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise networks store: %w", err)
	}

	nets, err := netStore.ListNetworks()
	if err != nil {
		return fmt.Errorf("couldn't list networks: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		if len(nets) == 0 {
			p.InfoText("No network registered").NextLine()
		}
		for _, net := range nets {
			p.Text(fmt.Sprintf("- %s", net)).NextLine()
		}
	} else if rootArgs.output == "json" {
		return vgjson.Print(struct {
			Networks []string `json:"networks"`
		}{
			Networks: nets,
		})
	}

	return nil
}
