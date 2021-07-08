package cmd

import (
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	taintArgs struct {
		name       string
		passphrase string
		pubkey     string
	}
	// taintCmd represents the taint command
	taintCmd = &cobra.Command{
		Use:   "taint",
		Short: "Taint a public key",
		Long:  "Taint a public key",
		RunE:  runTaint,
	}
)

func init() {
	rootCmd.AddCommand(taintCmd)
	taintCmd.Flags().StringVarP(&taintArgs.name, "name", "n", "", "Name of the wallet to use")
	taintCmd.Flags().StringVarP(&taintArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	taintCmd.Flags().StringVarP(&taintArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runTaint(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(taintArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	if len(taintArgs.passphrase) == 0 {
		taintArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	err = handler.TaintKey(taintArgs.name, taintArgs.pubkey, taintArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not taint the key: %v", err)
	}

	fmt.Printf("The key has been tainted.\n")
	return nil
}
