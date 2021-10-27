package cmd

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	keyTaintArgs struct {
		wallet         string
		passphraseFile string
		pubKey         string
	}

	keyTaintCmd = &cobra.Command{
		Use:   "taint",
		Short: "Taint a public key",
		Long:  "Taint a public key",
		RunE:  runKeyTaint,
	}
)

func init() {
	keyCmd.AddCommand(keyTaintCmd)
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.pubKey, "pubkey", "k", "", "Public key to be used (hex)")
	_ = keyTaintCmd.MarkFlagRequired("wallet")
	_ = keyTaintCmd.MarkFlagRequired("pubkey")
}

func runKeyTaint(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	passphrase, err := getPassphrase(keyTaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.TaintKey(keyTaintArgs.wallet, keyTaintArgs.pubKey, passphrase)
	if err != nil {
		return fmt.Errorf("could not taint the key: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Tainting succeeded").NextSection()

		p.RedArrow().DangerText("Important").NextLine()
		p.Text("If you tainted a key for security reasons, you should not untaint it.").NextSection()

		p.BlueArrow().InfoText("Untaint a key").NextLine()
		p.Text("You may have tainted a key by mistake. If you want to untaint it, use the following command:").NextSection()
		p.Code(fmt.Sprintf("%s key untaint --wallet \"%s\" --pubkey \"%s\"", os.Args[0], keyTaintArgs.wallet, keyTaintArgs.pubKey)).NextSection()
		p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
	}

	return nil
}
