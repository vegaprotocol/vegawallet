package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"github.com/spf13/cobra"
)

var (
	importArgs struct {
		name           string
		passphraseFile string
		mnemonicFile   string
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
	importCmd.Flags().StringVarP(&importArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	importCmd.Flags().StringVarP(&importArgs.mnemonicFile, "mnemonic-file", "m", "", `Path of the file containing the mnemonic of the wallet "swing ceiling chaos..."`)
}

func runImport(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(importArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	if len(importArgs.mnemonicFile) == 0 {
		return errors.New("path to wallet mnemonic is required")
	}

	passphrase, err := getPassphrase(importArgs.passphraseFile, true)
	if err != nil {
		return err
	}

	rawMnemonic, err := vgfs.ReadFile(importArgs.mnemonicFile)
	if err != nil {
		return fmt.Errorf("couldn't read mnemonic file: %w", err)
	}
	mnemonic := strings.Trim(string(rawMnemonic), "\n")

	err = handler.ImportWallet(importArgs.name, passphrase, mnemonic)
	if err != nil {
		return fmt.Errorf("couldn't import wallet: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Importing the wallet succeeded").NJump(2)

		p.BlueArrow().InfoText("Generate a key pair").Jump()
		p.Text("To generate a key pair on a given wallet, use the following command:").NJump(2)
		p.Code(fmt.Sprintf("%s key generate --name \"%s\"", os.Args[0], importArgs.name)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	}

	return nil
}
