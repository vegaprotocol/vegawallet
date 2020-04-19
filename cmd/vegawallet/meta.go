package main

import (
	"errors"
	"fmt"
	"strings"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
)

var (
	metas string
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Add metadata to a public key",
	Long:  "Add a list of metadata to a public key",
	RunE:  runMeta,
}

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.Flags().StringVarP(&walletOwner, "name", "n", "", "Name of the wallet to use")
	metaCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	metaCmd.Flags().StringVarP(&pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	metaCmd.Flags().StringVarP(&metas, "metas", "m", "", `A list of metadata e.g: "primary:true;asset;BTC"`)
}

func runMeta(cmd *cobra.Command, args []string) error {
	if len(walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return errors.New("passphrase is required")
	}
	if len(pubkey) <= 0 {
		return errors.New("pubkey is required")
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

	var meta []wallet.Meta
	if len(metas) > 0 {
		// expect ; separated metas
		metasSplit := strings.Split(metas, ";")
		for _, v := range metasSplit {
			metaVal := strings.Split(v, ":")
			if len(metaVal) != 2 {
				return fmt.Errorf("invalid meta format")
			}
			meta = append(meta, wallet.Meta{Key: metaVal[0], Value: metaVal[1]})
		}

	}

	kp.Meta = meta
	_, err = wallet.Write(wal, rootPath, walletOwner, passphrase)
	if err != nil {
		return err
	}

	return nil
}
