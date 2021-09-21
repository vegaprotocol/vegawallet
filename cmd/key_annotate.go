package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallets"
	"github.com/spf13/cobra"
)

var (
	keyAnnotateArgs struct {
		metadata       string
		wallet         string
		passphraseFile string
		pubkey         string
		clear          bool
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
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	keyAnnotateCmd.Flags().StringVarP(&keyAnnotateArgs.metadata, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
	keyAnnotateCmd.Flags().BoolVar(&keyAnnotateArgs.clear, "clear", false, "Clear the metadata")
}

func runKeyAnnotate(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(keyAnnotateArgs.wallet) == 0 {
		return errors.New("wallet is required")
	}
	if len(keyAnnotateArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(keyAnnotateArgs.metadata) == 0 && !keyAnnotateArgs.clear {
		return errors.New("meta is required")
	}
	if len(keyAnnotateArgs.metadata) != 0 && keyAnnotateArgs.clear {
		return errors.New("can't have `--meta` and `--clear` flags at the same time")
	}

	passphrase, err := getPassphrase(keyAnnotateArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	var metadata []wallet.Meta
	if !keyAnnotateArgs.clear {
		metadata, err = parseMeta(keyAnnotateArgs.metadata)
		if err != nil {
			return err
		}
	}

	err = handler.LoginWallet(keyAnnotateArgs.wallet, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	err = handler.UpdateMeta(keyAnnotateArgs.wallet, keyAnnotateArgs.pubkey, passphrase, metadata)
	if err != nil {
		return fmt.Errorf("could not update the metadata: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		if keyAnnotateArgs.clear {
			p.CheckMark().SuccessText("Annotation cleared").Jump()
		} else {
			p.CheckMark().SuccessText("Annotation succeeded").NJump(2)
			p.Text("New metadata:").Jump()
			printMeta(p, metadata)
		}
	}
	return nil
}
