package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/fsutil"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keypairs of a wallet",
	Long: "List all the keypairs for a given wallet",
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&walletOwner, "name", "n", "", "Name of the wallet to use")
	listCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase to access the wallet")

}

func  runList(cmd *cobra.Command, args []string) error {
	if len(walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return errors.New("passphrase is required")
	}

	if ok, err := fsutil.PathExists(rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, walletOwner, passphrase)
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
