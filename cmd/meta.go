package cmd

import (
	"errors"
	"fmt"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
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
	metaCmd.Flags().StringVarP(&metaArgs.metas, "metas", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runMeta(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
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

	metas, err := parseMeta(metaArgs.metas)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(metaArgs.walletOwner, metaArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	err = handler.UpdateMeta(metaArgs.walletOwner, metaArgs.pubkey, metaArgs.passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}

	fmt.Printf("The meta have been updated.\n")
	return nil
}
