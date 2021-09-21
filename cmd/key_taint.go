package cmd

import (
	"errors"
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallets"
	"github.com/spf13/cobra"
)

var (
	keyTaintArgs struct {
		wallet         string
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
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runKeyTaint(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(keyTaintArgs.wallet) == 0 {
		return errors.New("wallet is required")
	}

	passphrase, err := getPassphrase(keyTaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.TaintKey(keyTaintArgs.wallet, keyTaintArgs.pubkey, passphrase)
	if err != nil {
		return fmt.Errorf("could not taint the key: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Tainting succeeded").NJump(2)

		p.RedArrow().DangerText("Important").Jump()
		p.Text("If you tainted a key for security reasons, you should not untaint it.").NJump(2)

		p.BlueArrow().InfoText("Untaint a key").Jump()
		p.Text("You may have tainted a key by mistake. If you want to untaint it, use the following command:").NJump(2)
		p.Code(fmt.Sprintf("%s key untaint --wallet \"%s\" --pubkey \"%s\"", os.Args[0], keyTaintArgs.wallet, keyTaintArgs.pubkey)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	}

	return nil
}
