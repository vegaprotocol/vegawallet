package cmd

import (
	"fmt"
	"os"
	"strings"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	importArgs struct {
		wallet         string
		passphraseFile string
		mnemonicFile   string
		version        uint32
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
	importCmd.Flags().StringVarP(&importArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	importCmd.Flags().StringVarP(&importArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	importCmd.Flags().StringVarP(&importArgs.mnemonicFile, "mnemonic-file", "m", "", `Path of the file containing the mnemonic of the wallet "swing ceiling chaos..."`)
	importCmd.Flags().Uint32Var(&importArgs.version, "version", wallet.LatestVersion, fmt.Sprintf("Version of the wallet to import: %v", wallet.SupportedVersions))
	_ = importCmd.MarkFlagRequired("wallet")
	_ = importCmd.MarkFlagRequired("mnemonic-file")
}

func runImport(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	passphrase, err := getPassphrase(importArgs.passphraseFile, true)
	if err != nil {
		return err
	}

	rawMnemonic, err := vgfs.ReadFile(importArgs.mnemonicFile)
	if err != nil {
		return fmt.Errorf("couldn't read mnemonic file: %w", err)
	}
	mnemonic := strings.Trim(string(rawMnemonic), "\n")

	err = handler.ImportWallet(importArgs.wallet, passphrase, mnemonic, importArgs.version)
	if err != nil {
		return fmt.Errorf("couldn't import wallet: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Importing the wallet succeeded").NextSection()

		p.BlueArrow().InfoText("Generate a key pair").NextLine()
		p.Text("To generate a key pair on a given wallet, use the following command:").NextSection()
		p.Code(fmt.Sprintf("%s key generate --wallet \"%s\"", os.Args[0], importArgs.wallet)).NextSection()
		p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
	}

	return nil
}
