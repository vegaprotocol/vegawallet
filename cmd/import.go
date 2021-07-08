package cmd

import (
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	walletImportArgs struct {
		name       string
		passphrase string
		mnemonic   string
	}

	walletImportCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a wallet using the mnemonic",
		Long:  "Import a wallet using the mnemonic.",
		RunE:  importWallet,
	}
)

func init() {
	rootCmd.AddCommand(walletImportCmd)
	walletImportCmd.Flags().StringVarP(&walletImportArgs.name, "name", "n", "", "Name of the wallet to use")
	walletImportCmd.Flags().StringVarP(&walletImportArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	walletImportCmd.Flags().StringVarP(&walletImportArgs.mnemonic, "mnemonic", "m", "", `Mnemonic of the wallet "swing ceiling chaos..."`)
}

func importWallet(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(walletImportArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	if len(walletImportArgs.mnemonic) == 0 {
		return errors.New("wallet mnemonic is required")
	}

	if len(walletImportArgs.passphrase) == 0 {
		var (
			err          error
			confirmation string
		)
		walletImportArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		if len(walletImportArgs.passphrase) == 0 {
			return fmt.Errorf("passphrase cannot be empty")
		}

		confirmation, err = promptForPassphrase("please confirm passphrase:")
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		if walletImportArgs.passphrase != confirmation {
			return fmt.Errorf("passphrases do not match")
		}
	}

	err = handler.ImportWallet(walletImportArgs.name, walletImportArgs.passphrase, walletImportArgs.mnemonic)
	if err != nil {
		return fmt.Errorf("couldn't import wallet: %v", err)
	}

	fmt.Printf("The wallet \"%s\" has been imported.\n", walletImportArgs.name)

	return nil
}
