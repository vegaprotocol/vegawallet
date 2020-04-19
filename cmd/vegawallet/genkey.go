package main

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/spf13/cobra"
)

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a new keypair for a wallet",
	Long:  "Generate a new keypair for a wallet, this will implicitly generate a new wallet if none exist for the given name",
	RunE:  runGenkey,
}

func init() {
	rootCmd.AddCommand(genkeyCmd)
	genkeyCmd.Flags().StringVarP(&walletOwner, "name", "n", "", "Name of the wallet to use")
	genkeyCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
}

func runGenkey(cmd *cobra.Command, args []string) error {
	if len(walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return errors.New("passphrase is required")
	}

	if ok, err := fsutil.PathExists(rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	if err := wallet.EnsureBaseFolder(rootPath); err != nil {
		return fmt.Errorf("unable to initialization root folder: %v", err)
	}

	_, err := wallet.Read(rootPath, walletOwner, passphrase)
	if err != nil {
		if err != wallet.ErrWalletDoesNotExists {
			// this an invalid key, returning error
			return fmt.Errorf("unable to decrypt wallet: %v", err)
		}
		// wallet do not exit, let's try to create it
		_, err = wallet.Create(rootPath, walletOwner, passphrase)
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

	// now updating the wallet and saving it
	_, err = wallet.AddKeypair(kp, rootPath, walletOwner, passphrase)
	if err != nil {
		return fmt.Errorf("unable to add keypair to wallet: %v", err)
	}

	// print the new keys for user info
	fmt.Printf("new generated keys:\n")
	fmt.Printf("public: %v\n", kp.Pub)
	fmt.Printf("private: %v\n", kp.Priv)

	return nil
}
