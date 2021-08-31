package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"github.com/spf13/cobra"
)

var (
	keyUntaintArgs struct {
		name           string
		passphraseFile string
		pubKey         string
	}

	keyUntaintCmd = &cobra.Command{
		Use:   "untaint",
		Short: "Untaint a public key",
		Long:  "Untaint a public key",
		RunE:  runKeyUntaint,
	}
)

func init() {
	keyCmd.AddCommand(keyUntaintCmd)
	keyUntaintCmd.Flags().StringVarP(&keyUntaintArgs.name, "name", "n", "", "Name of the wallet to use")
	keyUntaintCmd.Flags().StringVarP(&keyUntaintArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyUntaintCmd.Flags().StringVarP(&keyUntaintArgs.pubKey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runKeyUntaint(_ *cobra.Command, _ []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyUntaintArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	passphrase, err := getPassphrase(keyUntaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.UntaintKey(keyUntaintArgs.name, keyUntaintArgs.pubKey, passphrase)
	if err != nil {
		return fmt.Errorf("could not untaint the key: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Untainting succeeded").NJump(2)

		p.RedArrow().DangerText("Important").Jump()
		p.Text("If you tainted a key for security reasons, you should not use it.").Jump()
	}

	return nil
}
