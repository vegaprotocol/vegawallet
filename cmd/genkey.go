package cmd

import (
	"errors"
	"fmt"
	"strings"

	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	genKeyArgs struct {
		walletOwner string
		passphrase  string
		metas       string
	}

	// genKeyCmd represents the genkey command
	genKeyCmd = &cobra.Command{
		Use:   "genkey",
		Short: "Generate a new keypair for a wallet",
		Long:  "Generate a new keypair for a wallet, this will implicitly generate a new wallet if none exist for the given name",
		RunE:  runGenKey,
	}
)

func init() {
	rootCmd.AddCommand(genKeyCmd)
	genKeyCmd.Flags().StringVarP(&genKeyArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	genKeyCmd.Flags().StringVarP(&genKeyArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	genKeyCmd.Flags().StringVarP(&genKeyArgs.metas, "metas", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runGenKey(cmd *cobra.Command, args []string) error {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)

	if len(genKeyArgs.walletOwner) == 0 {
		return errors.New("wallet name is required")
	}

	walletExists := handler.WalletExists(genKeyArgs.walletOwner)

	if len(genKeyArgs.passphrase) == 0 {
		var (
			err          error
			confirmation string
		)
		genKeyArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		if len(genKeyArgs.passphrase) == 0 {
			return fmt.Errorf("passphrase cannot be empty")
		}

		if !walletExists {
			confirmation, err = promptForPassphrase("please confirm passphrase:")
			if err != nil {
				return fmt.Errorf("could not get passphrase: %v", err)
			}

			if genKeyArgs.passphrase != confirmation {
				return fmt.Errorf("passphrases do not match")
			}
		}
	}

	metas, err := parseMeta(genKeyArgs.metas)
	if err != nil {
		return err
	}

	if !walletExists {
		err := handler.CreateWallet(genKeyArgs.walletOwner, genKeyArgs.passphrase)
		if err != nil {
			return fmt.Errorf("couldn't create wallet: %v", err)
		}
	}

	keyPair, err := handler.GenerateKeyPair(genKeyArgs.walletOwner, genKeyArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not generate a key pair: %v", err)
	}

	err = handler.UpdateMeta(genKeyArgs.walletOwner, keyPair.Pub, genKeyArgs.passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}

	// print the new keys for user info
	fmt.Printf("new generated keys:\n")
	fmt.Printf("public: %v\n", keyPair.Pub)
	fmt.Printf("private: %v\n", keyPair.Priv)

	return nil
}

func parseMeta(metaStr string) ([]wallet.Meta, error) {
	var metas []wallet.Meta

	if len(metaStr) == 0 {
		return metas, nil
	}

	rawMetas := strings.Split(metaStr, ";")
	for _, v := range rawMetas {
		rawMeta := strings.Split(v, ":")
		if len(rawMeta) != 2 {
			return nil, fmt.Errorf("invalid meta format")
		}
		metas = append(metas, wallet.Meta{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metas, nil
}
