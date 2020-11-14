package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/spf13/cobra"
)

var (
	signArgs struct {
		walletOwner string
		passphrase  string
		message     string
		pubkey      string
	}

	// signCmd represents the sign command
	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a blob of data",
		Long:  "Sign a blob of dara base64 encoded",
		RunE:  runSign,
	}
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringVarP(&signArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	signCmd.Flags().StringVarP(&signArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	signCmd.Flags().StringVarP(&signArgs.message, "message", "m", "", "Message to be signed (base64)")
	signCmd.Flags().StringVarP(&signArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runSign(cmd *cobra.Command, args []string) error {
	if len(signArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(signArgs.passphrase) <= 0 {
		var err error
		signArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}
	if len(signArgs.pubkey) <= 0 {
		return errors.New("pubkey is required")
	}
	if len(signArgs.message) <= 0 {
		return errors.New("data is required")
	}

	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootArgs.rootPath, signArgs.walletOwner, signArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	dataBuf, err := base64.StdEncoding.DecodeString(signArgs.message)
	if err != nil {
		return fmt.Errorf("invalid base64 encoded data: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == signArgs.pubkey {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return fmt.Errorf("unknown public key: %v", signArgs.pubkey)
	}
	if kp.Tainted {
		return fmt.Errorf("key is tainted: %v", signArgs.pubkey)
	}

	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("unable to instanciate signature algorithm: %v", err)
	}
	sig, err := wallet.Sign(alg, kp, dataBuf)
	if err != nil {
		return fmt.Errorf("unable to sign: %v", err)
	}
	fmt.Printf("%v\n", base64.StdEncoding.EncodeToString(sig))

	return nil
}
