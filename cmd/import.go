package cmd

import (
	"errors"
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallet"
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

func runImport(_ *cobra.Command, _ []string) error {
	store, err := newWalletsStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)
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
		return fmt.Errorf("couldn't import wallet: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().Text("Service configuration created at: ").SuccessText(rootArgs.rootPath).Jump()
		p.CheckMark().Text("Service RSA keys created at: ").SuccessText(rootArgs.rootPath).Jump()
		p.CheckMark().SuccessText("Importing the wallet succeeded").NJump(2)

		p.BlueArrow().InfoText("Generate a key pair").Jump()
		p.Text("To generate a key pair on a given wallet, use the following command:").NJump(2)
		p.Code(fmt.Sprintf("%s key generate --name \"%s\"", os.Args[0], importArgs.name)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	}

	return nil
}
