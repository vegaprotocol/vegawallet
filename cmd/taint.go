package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	taintArgs struct {
		walletOwner string
		passphrase  string
		pubkey      string
	}
	// taintCmd represents the taint command
	taintCmd = &cobra.Command{
		Use:   "taint",
		Short: "Taint a public key",
		Long:  "Taint a public key",
		RunE:  runTaint,
	}
)

func init() {
	rootCmd.AddCommand(taintCmd)
	taintCmd.Flags().StringVarP(&taintArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	taintCmd.Flags().StringVarP(&taintArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	taintCmd.Flags().StringVarP(&taintArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runTaint(cmd *cobra.Command, args []string) error {
	if len(taintArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(taintArgs.passphrase) <= 0 {
		var err error
		taintArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootArgs.rootPath, taintArgs.walletOwner, taintArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == taintArgs.pubkey {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return fmt.Errorf("unknown public key: %v", taintArgs.pubkey)
	}

	if kp.Tainted {
		return fmt.Errorf("key %v is already tainted", taintArgs.pubkey)
	}

	kp.Tainted = true

	_, err = wallet.Write(wal, rootArgs.rootPath, taintArgs.walletOwner, taintArgs.passphrase)
	if err != nil {
		return err
	}

	return nil
}
