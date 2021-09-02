package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
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
	handler, err := newWalletHandler(rootArgs.vegaHome)
	if err != nil {
		return err
	}

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
