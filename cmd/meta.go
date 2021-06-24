package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	metaArgs struct {
		metas       string
		walletOwner string
		passphrase  string
		pubkey      string
	}

	// metaCmd represents the meta command
	metaCmd = &cobra.Command{
		Use:   "meta",
		Short: "Add metadata to a public key",
		Long:  "Add a list of metadata to a public key",
		RunE:  runMeta,
	}
)

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.Flags().StringVarP(&metaArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	metaCmd.Flags().StringVarP(&metaArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	metaCmd.Flags().StringVarP(&metaArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	metaCmd.Flags().StringVarP(&metaArgs.metas, "metas", "m", "", `A list of metadata e.g: "primary:true;asset;BTC"`)
}

func runMeta(cmd *cobra.Command, args []string) error {
	store, err := wallet.NewFileStoreV1(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(metaArgs.walletOwner) == 0 {
		return errors.New("wallet name is required")
	}
	if len(metaArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(metaArgs.passphrase) == 0 {
		var err error
		metaArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	metas, err := parseMeta(genKeyArgs.metas)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(signArgs.walletOwner, signArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	err = handler.UpdateMeta(genKeyArgs.walletOwner, metaArgs.pubkey, genKeyArgs.passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}
	return nil
}
