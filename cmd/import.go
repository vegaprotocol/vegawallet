package cmd

import (
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	importWalletArgs struct {
		name       string
		passphrase string
		mnemonic   string
	}

	// importWalletCmd represents the import command
	importWalletCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a wallet using the mnemonic",
		Long:  "Import a wallet using the mnemonic.",
		RunE:  importWallet,
	}
)

func init() {
	rootCmd.AddCommand(importWalletCmd)
	importWalletCmd.Flags().StringVarP(&importWalletArgs.name, "name", "n", "", "Name of the wallet to use")
	importWalletCmd.Flags().StringVarP(&importWalletArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	importWalletCmd.Flags().StringVarP(&importWalletArgs.mnemonic, "mnemonic", "m", "", `Mnemonic of the wallet "swing ceiling chaos..."`)
}

func importWallet(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(importWalletArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	if len(importWalletArgs.mnemonic) == 0 {
		return errors.New("wallet mnemonic is required")
	}

	if len(importWalletArgs.passphrase) == 0 {
		var (
			err          error
			confirmation string
		)
		importWalletArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		if len(importWalletArgs.passphrase) == 0 {
			return fmt.Errorf("passphrase cannot be empty")
		}

		confirmation, err = promptForPassphrase("please confirm passphrase:")
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		if importWalletArgs.passphrase != confirmation {
			return fmt.Errorf("passphrases do not match")
		}
	}

	err = handler.ImportWallet(importWalletArgs.name, importWalletArgs.passphrase, importWalletArgs.mnemonic)
	if err != nil {
		return fmt.Errorf("couldn't import wallet: %v", err)
	}

	fmt.Printf("The wallet \"%s\" has been imported.\n", importWalletArgs.name)

	return nil
}
