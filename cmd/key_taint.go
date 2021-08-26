package cmd

import (
	"errors"
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"github.com/spf13/cobra"
)

var (
	keyTaintArgs struct {
		name           string
		passphrase     string
		passphraseFile string
		pubkey         string
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
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.name, "name", "n", "", "Name of the wallet to use")
	keyTaintCmd.Flags().StringVar(&keyTaintArgs.passphrase, "passphrase", "", "Passphrase to access the wallet")
	keyTaintCmd.Flags().StringVar(&keyTaintArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runKeyTaint(_ *cobra.Command, _ []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyTaintArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	passphrase, err := getPassphrase(keyTaintArgs.passphrase, keyTaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.TaintKey(keyTaintArgs.name, keyTaintArgs.pubkey, passphrase)
	if err != nil {
		return fmt.Errorf("could not taint the key: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().Text("Key pair has been tainted").Jump()
		p.CheckMark().SuccessText("Tainting succeeded").NJump(2)

		p.RedArrow().DangerText("Important").Jump()
		p.Text("If you tainted a key for security reasons, you should not untaint it.").NJump(2)

		p.BlueArrow().InfoText("Untaint a key").Jump()
		p.Text("You may have tainted a key by mistake. If you want to untaint it, use the following command:").NJump(2)
		p.Code(fmt.Sprintf("%s key untaint --name \"%s\" --pubkey \"%s\"", os.Args[0], keyTaintArgs.name, keyTaintArgs.pubkey)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	}

	return nil
}
