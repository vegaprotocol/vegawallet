package main

import (
	"encoding/base64"
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/spf13/cobra"
)

var (
	sig string
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the signature",
	Long: "Verify the signature for a blob of data",
	RunE: runVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&walletOwner, "name", "n", "", "Name of the wallet to use")
	verifyCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	verifyCmd.Flags().StringVarP(&data, "message", "m", "", "Message to be verified (base64)")
	verifyCmd.Flags().StringVarP(&sig, "signature", "s", "", "Signature to be verified (base64)")
	verifyCmd.Flags().StringVarP(&pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runVerify(cmd *cobra.Command, args []string) error {
	if len(walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(passphrase) <= 0 {
		return errors.New("passphrase is required")
	}
	if len(pubkey) <= 0 {
		return errors.New("pubkey is required")
	}
	if len(data) <= 0 {
		return errors.New("data is required")
	}
	if len(sig) <= 0 {
		return errors.New("data is required")
	}

	if ok, err := fsutil.PathExists(rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootPath, walletOwner, passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	dataBuf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid base64 encoded data: %v", err)
	}
	sigBuf, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return fmt.Errorf("invalid base64 encoded data: %v", err)
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

	alg, err := crypto.NewSignatureAlgorithm(crypto.Ed25519)
	if err != nil {
		return fmt.Errorf("unable to instanciate signature algorithm: %v", err)
	}
	verified, err := wallet.Verify(alg, kp, dataBuf, sigBuf)
	if err != nil {
		return fmt.Errorf("unable to verify: %v", err)
	}
	fmt.Printf("%v\n", verified)

	return nil
}
