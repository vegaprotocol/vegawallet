package cmd

import (
	vgjson "code.vegaprotocol.io/go-wallet/libs/json"

	"github.com/spf13/cobra"
)

var (
	listArgs struct{}

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

func runList(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	wallets, err := handler.ListWallets()
	if err != nil {
		return err
	}

	return vgjson.PrettyPrint(wallets)
}
