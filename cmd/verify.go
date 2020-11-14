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
	verifyArgs struct {
		walletOwner string
		passphrase  string
		sig         string
		message     string
		pubkey      string
	}

	// verifyCmd represents the verify command
	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify the signature",
		Long:  "Verify the signature for a blob of data",
		RunE:  runVerify,
	}
)

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&verifyArgs.walletOwner, "name", "n", "", "Name of the wallet to use")
	verifyCmd.Flags().StringVarP(&verifyArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	verifyCmd.Flags().StringVarP(&verifyArgs.message, "message", "m", "", "Message to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.sig, "signature", "s", "", "Signature to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runVerify(cmd *cobra.Command, args []string) error {
	if len(verifyArgs.walletOwner) <= 0 {
		return errors.New("wallet name is required")
	}
	if len(verifyArgs.passphrase) <= 0 {
		var err error
		verifyArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}
	if len(verifyArgs.pubkey) <= 0 {
		return errors.New("pubkey is required")
	}
	if len(verifyArgs.message) <= 0 {
		return errors.New("message is required")
	}
	if len(verifyArgs.sig) <= 0 {
		return errors.New("data is required")
	}

	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		return fmt.Errorf("invalid root directory path: %v", err)
	}

	wal, err := wallet.Read(rootArgs.rootPath, verifyArgs.walletOwner, verifyArgs.passphrase)
	if err != nil {
		return fmt.Errorf("unable to decrypt wallet: %v", err)
	}

	dataBuf, err := base64.StdEncoding.DecodeString(verifyArgs.message)
	if err != nil {
		return fmt.Errorf("invalid base64 encoded data: %v", err)
	}
	sigBuf, err := base64.StdEncoding.DecodeString(verifyArgs.sig)
	if err != nil {
		return fmt.Errorf("invalid base64 encoded data: %v", err)
	}

	var kp *wallet.Keypair
	for i, v := range wal.Keypairs {
		if v.Pub == verifyArgs.pubkey {
			kp = &wal.Keypairs[i]
		}
	}
	if kp == nil {
		return fmt.Errorf("unknown public key: %v", verifyArgs.pubkey)
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
