package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"github.com/spf13/cobra"
)

var (
	keyAnnotateArgs struct {
		metadata       string
		name           string
		passphraseFile string
		pubkey         string
	}

	keyAnnotateCmd = &cobra.Command{
		Use:   "annotate",
		Short: "Add metadata to a public key",
		Long:  "Add a list of metadata to a public key",
		RunE:  runKeyAnnotate,
	}
)

func init() {
	keyCmd.AddCommand(keyAnnotateCmd)
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.name, "name", "n", "", "Name of the wallet to use")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.metadata, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runKeyAnnotate(_ *cobra.Command, _ []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyAnnotateArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(keyAnnotateArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}

	passphrase, err := getPassphrase(keyAnnotateArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	metadata, err := parseMeta(keyAnnotateArgs.metadata)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(keyAnnotateArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	err = handler.UpdateMeta(keyAnnotateArgs.name, keyAnnotateArgs.pubkey, passphrase, metadata)
	if err != nil {
		return fmt.Errorf("could not update the metadata: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Annotation succeeded").NJump(2)

		p.Text("New metadata:").Jump()
		printMeta(p, metadata)
	}
	return nil
}
