package cmd

import (
	"errors"
	"fmt"
	"strings"

	"code.vegaprotocol.io/go-wallet/fsutil"
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
	if len(metaArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(metaArgs.passphrase) <= 0 {
		var err error
		metaArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}
	if len(metaArgs.pubkey) <= 0 {
		return errors.New("pubkey is required")
	}
	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootArgs.rootPath, metaArgs.walletOwner, metaArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == metaArgs.pubkey {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return fmt.Errorf("unknown public key: %v", metaArgs.pubkey)
	}

	var meta []wallet.Meta
	if len(metaArgs.metas) > 0 {
		// expect ; separated metas
		metasSplit := strings.Split(metaArgs.metas, ";")
		for _, v := range metasSplit {
			metaVal := strings.Split(v, ":")
			if len(metaVal) != 2 {
				return fmt.Errorf("invalid meta format")
			}
			meta = append(meta, wallet.Meta{Key: metaVal[0], Value: metaVal[1]})
		}
	}

	kp.Meta = meta
	_, err = wallet.Write(wal, rootArgs.rootPath, metaArgs.walletOwner, metaArgs.passphrase)
	if err != nil {
		return err
	}

	return nil
}
