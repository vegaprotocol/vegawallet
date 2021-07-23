package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	importArgs struct {
		name           string
		passphrase     string
		passphraseFile string
		mnemonic       string
	}

	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a wallet using the mnemonic",
		Long:  "Import a wallet using the mnemonic.",
		RunE:  runImport,
	}
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&importArgs.name, "name", "n", "", "Name of the wallet to use")
	importCmd.Flags().StringVar(&importArgs.passphrase, "passphrase", "", "Passphrase to access the wallet")
	importCmd.Flags().StringVar(&importArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	importCmd.Flags().StringVarP(&importArgs.mnemonic, "mnemonic", "m", "", `Mnemonic of the wallet "swing ceiling chaos..."`)
}

func runImport(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(importArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	if len(importArgs.mnemonic) == 0 {
		return errors.New("wallet mnemonic is required")
	}

	passphrase, err := getPassphrase(importArgs.passphrase, importArgs.passphraseFile, true)
	if err != nil {
		return err
	}

	err = handler.ImportWallet(importArgs.name, passphrase, importArgs.mnemonic)
	if err != nil {
		return fmt.Errorf("couldn't import wallet: %v", err)
	}

	fmt.Printf("The wallet \"%s\" has been imported.\n", importArgs.name)

	return nil
}
