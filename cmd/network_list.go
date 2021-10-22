package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegawallet/cmd/printer"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
	"github.com/spf13/cobra"
)

var (
	// networkListCmd represents the network list command
	networkListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all registered network",
		Long:  "List all registered network",
		RunE:  runNetworkList,
	}
)

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
			p.InfoText("No network registered").Jump()
		}
		for _, net := range nets {
			p.Text(fmt.Sprintf("- %s", net)).Jump()
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
