package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/spf13/cobra"
)

var (
	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all registered wallets",
		Long:  "List all registered wallets",
		RunE:  runList,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	wallets, err := handler.ListWallets()
	if err != nil {
		return err
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		for _, w := range wallets {
			p.Text(fmt.Sprintf("- %s", w)).Jump()
		}
	} else if rootArgs.output == "json" {
		return vgjson.Print(struct {
			Wallets []string
		}{
			Wallets: wallets,
		})
	}

	return nil
}
