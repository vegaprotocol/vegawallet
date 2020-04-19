package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/fsutil"
)

// taintCmd represents the taint command
var taintCmd = &cobra.Command{
	Use:   "taint",
	Short: "Taint a public key",
	Long: "Taint a public key",
	RunE: runTaint,
}

func init() {
	rootCmd.AddCommand(taintCmd)
	taintCmd.Flags().StringVarP(&walletOwner, "name", "n", "", "Name of the wallet to use")
	taintCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	taintCmd.Flags().StringVarP(&pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func  runTaint(cmd *cobra.Command, args []string) error {
	if len(walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return errors.New("passphrase is required")
	}

	if ok, err := fsutil.PathExists(rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, walletOwner, passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == pubkey {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return fmt.Errorf("unknown public key: %v", pubkey)
	}

	if kp.Tainted {
		return fmt.Errorf("key %v is already tainted", pubkey)
	}

	kp.Tainted = true

	_, err = wallet.Write(wal, rootPath, walletOwner, passphrase)
	if err != nil {
		return err
	}

	return nil
}
