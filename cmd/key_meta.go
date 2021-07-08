package cmd

import (
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	keyMetaArgs struct {
		meta       string
		name       string
		passphrase string
		pubkey     string
	}

	keyMetaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Add metadata to a public key",
		Long:  "Add a list of metadata to a public key",
		RunE:  runMeta,
	}
)

func init() {
	keyCmd.AddCommand(keyMetaCmd)
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.name, "name", "n", "", "Name of the wallet to use")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.meta, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runMeta(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(keyMetaArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(keyMetaArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(keyMetaArgs.passphrase) == 0 {
		var err error
		keyMetaArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	meta, err := parseMeta(keyMetaArgs.meta)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(keyMetaArgs.name, keyMetaArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	err = handler.UpdateMeta(keyMetaArgs.name, keyMetaArgs.pubkey, keyMetaArgs.passphrase, meta)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}

	fmt.Printf("The meta have been updated.\n")
	return nil
}
