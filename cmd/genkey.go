package cmd

import (
	"errors"
	"fmt"
	"strings"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/spf13/cobra"
)

var (
	genkeyArgs struct {
		walletOwner string
		passphrase  string
		metas       string
	}

	// genkeyCmd represents the genkey command
	genkeyCmd = &cobra.Command{
		Use:   "genkey",
		Short: "Generate a new keypair for a wallet",
		Long:  "Generate a new keypair for a wallet, this will implicitly generate a new wallet if none exist for the given name",
		RunE:  runGenkey,
	}
)

func init() {
	rootCmd.AddCommand(genkeyCmd)
	genkeyCmd.Flags().StringVarP(&genkeyArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	genkeyCmd.Flags().StringVarP(&genkeyArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	genkeyCmd.Flags().StringVarP(&genkeyArgs.metas, "metas", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runGenkey(cmd *cobra.Command, args []string) error {
	if len(genkeyArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(genkeyArgs.passphrase) <= 0 {
		var (
			err          error
			confirmation string
		)
		genkeyArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}

		// if wallet does not exists
		// ask for passphrase confirmation + check it's not empty
		if !wallet.WalletFileExists(rootArgs.rootPath, genkeyArgs.walletOwner) {
			confirmation, err = promptForPassphrase("please confirm passphrase:")
			if err != nil {
				return fmt.Errorf("could not get passphrase: %v", err)
			}

			if genkeyArgs.passphrase != confirmation {
				return fmt.Errorf("passphrases do not match")
			}

			if len(genkeyArgs.passphrase) <= 0 {
				return fmt.Errorf("passphrase cannot be empty")
			}
		}
	}

	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(rootArgs.rootPath); err != nil {
			return fmt.Errorf("error creating root directory: %v", err)
		}
	}

	if err := wallet.EnsureBaseFolder(rootArgs.rootPath); err != nil {
		return fmt.Errorf("unable to initialization root folder: %v", err)
	}

	_, err := wallet.Read(rootArgs.rootPath, genkeyArgs.walletOwner, genkeyArgs.passphrase)
	if err != nil {
		if err != wallet.ErrWalletDoesNotExists {
			// this an invalid key, returning error
			return fmt.Errorf("unable to decrypt wallet: %v", err)
		}
		// wallet do not exit, let's try to create it
		_, err = wallet.Create(rootArgs.rootPath, genkeyArgs.walletOwner, genkeyArgs.passphrase)
		if err != nil {
			return fmt.Errorf("unable to create wallet: %v", err)
		}
	}

	// at this point we have a valid wallet
	// let's generate the keypair
	// defaulting to ed25519 for now
	algo := crypto.NewEd25519()
	kp, err := wallet.GenKeypair(algo.Name())
	if err != nil {
		return fmt.Errorf("unable to generate new key pair: %v", err)
	}

	var meta []wallet.Meta
	if len(genkeyArgs.metas) > 0 {
		// expect ; separated metas
		metasSplit := strings.Split(genkeyArgs.metas, ";")
		for _, v := range metasSplit {
			metaVal := strings.Split(v, ":")
			if len(metaVal) != 2 {
				return fmt.Errorf("invalid meta format")
			}
			meta = append(meta, wallet.Meta{Key: metaVal[0], Value: metaVal[1]})
		}
	}

	kp.Meta = meta

	// now updating the wallet and saving it
	_, err = wallet.AddKeypair(kp, rootArgs.rootPath, genkeyArgs.walletOwner, genkeyArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to add keypair to wallet: %v", err)
	}

	// print the new keys for user info
	fmt.Printf("new generated keys:\n")
	fmt.Printf("public: %v\n", kp.Pub)
	fmt.Printf("private: %v\n", kp.Priv)

	return nil
}
