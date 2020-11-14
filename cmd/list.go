package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	listArgs struct {
		walletOwner string
		passphrase  string
	}

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List keypairs of a wallet",
		Long:  "List all the keypairs for a given wallet",
		RunE:  runList,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	listCmd.Flags().StringVarP(&listArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")

}

func runList(cmd *cobra.Command, args []string) error {
	if len(listArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(listArgs.passphrase) <= 0 {
		var err error
		listArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootArgs.rootPath, listArgs.walletOwner, listArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	buf, err := json.MarshalIndent(wal, " ", " ")
	if err != nil {
		return fmt.Errorf("unable to marshal message: %v", err)
	}

	// print the new keys for user info
	fmt.Printf("List of all your keypairs:\n")
	fmt.Printf("%v\n", string(buf))

	return nil
}
