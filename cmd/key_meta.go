package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	keyMetaArgs struct {
		meta           string
		name           string
		passphrase     string
		passphraseFile string
		pubkey         string
	}

	keyMetaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Add metadata to a public key",
		Long:  "Add a list of metadata to a public key",
		RunE:  runKeyMeta,
	}
)

func init() {
	keyCmd.AddCommand(keyMetaCmd)
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.name, "name", "n", "", "Name of the wallet to use")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	keyMetaCmd.Flags().StringVar(&keyMetaArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	keyMetaCmd.Flags().StringVarP(&keyMetaArgs.meta, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runKeyMeta(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyMetaArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(keyMetaArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}

	passphrase, err := getPassphrase(keyMetaArgs.passphrase, keyMetaArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	meta, err := parseMeta(keyMetaArgs.meta)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(keyMetaArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	err = handler.UpdateMeta(keyMetaArgs.name, keyMetaArgs.pubkey, passphrase, meta)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}

	fmt.Printf("The meta have been updated.\n")
	return nil
}
