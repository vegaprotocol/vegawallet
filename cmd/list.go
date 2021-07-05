package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	listArgs struct {
		name       string
		passphrase string
	}

	// listCmd represents the list command
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List keys of a wallet",
		Long:  "List all the keys for a given wallet",
		RunE:  runList,
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listArgs.name, "name", "n", "", "Name of the wallet to use")
	listCmd.Flags().StringVarP(&listArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
}

func runList(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(listArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(listArgs.passphrase) == 0 {
		var err error
		listArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	err = handler.LoginWallet(listArgs.name, listArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	keys, err := handler.ListKeyPairs(listArgs.name)
	if err != nil {
		return fmt.Errorf("could not list the public keys: %v", err)
	}

	buf, err := json.MarshalIndent(keys, " ", " ")
	if err != nil {
		return fmt.Errorf("unable to marshal message: %v", err)
	}

	// print the new keys for user info
	fmt.Printf("List of all your keys:\n")
	fmt.Printf("%v\n", string(buf))

	return nil
}
