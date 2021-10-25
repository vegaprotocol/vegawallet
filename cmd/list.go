package cmd

import (
	"fmt"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered wallets",
	Long:  "List all registered wallets",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	ws, err := handler.ListWallets()
	if err != nil {
		return err
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		if len(ws) == 0 {
			p.InfoText("No wallet registered").Jump()
		}
		for _, w := range ws {
			p.Text(fmt.Sprintf("- %s", w)).Jump()
		}
	} else if rootArgs.output == "json" {
		return vgjson.Print(struct {
			Wallets []string `json:"wallets"`
		}{
			Wallets: ws,
		})
	}

	return nil
}
