package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	wcommands "code.vegaprotocol.io/go-wallet/commands"
	"code.vegaprotocol.io/go-wallet/logger"
	"code.vegaprotocol.io/go-wallet/node"
	"code.vegaprotocol.io/go-wallet/wallets"
	api "code.vegaprotocol.io/protos/vega/api/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var (
	commandArgs struct {
		name           string
		passphraseFile string
		nodeAddress    string
		retries        uint64
		pubKey         string
	}

	commandCmd = &cobra.Command{
		Use:   "command",
		Short: "Send a command to the vega network",
		Long:  "Import a wallet using the mnemonic.",
		RunE:  runCommand,
	}
)

func init() {
	rootCmd.AddCommand(commandCmd)
	commandCmd.Flags().StringVarP(&commandArgs.name, "name", "n", "", "Name of the wallet to use")
	commandCmd.Flags().StringVarP(&commandArgs.pubKey, "pubkey", "", "", "The public key to use from the wallet")
	commandCmd.Flags().StringVarP(&commandArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	commandCmd.Flags().StringVar(&commandArgs.nodeAddress, "node-address", "0.0.0.0:3002", "Address of the Vega node to use")
	commandCmd.Flags().Uint64Var(&commandArgs.retries, "retries", 5, "Number of retries when contacting the Vega node")
	commandCmd.MarkFlagRequired("name")
	commandCmd.MarkFlagRequired("pubkey")
}

func runCommand(_ *cobra.Command, pos []string) error {
	if len(pos) != 1 {
		return errors.New("invalid number of arguments, require at most 1 command to be signed and sent by the wallet")
	}

	wReq := &walletpb.SubmitTransactionRequest{}
	err := jsonpb.UnmarshalString(pos[0], wReq)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal request: %w", err)
	}

	wReq.PubKey = commandArgs.pubKey

	errs := wcommands.CheckSubmitTransactionRequest(wReq)
	if !errs.Empty() {
		return fmt.Errorf("invalid request: %w", err)
	}

	passphrase, err := getPassphrase(importArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	err = handler.LoginWallet(commandArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't login to the wallet %s: %w", commandArgs.name, err)
	}
	defer handler.LogoutWallet(commandArgs.name)

	encoding := "json"
	if rootArgs.output == "human" {
		encoding = "console"
	}

	log, err := logger.New(zapcore.InfoLevel, encoding)
	if err != nil {
		return err
	}
	defer log.Sync()

	forwarder, err := node.NewForwarder(log.Named("forwarder"), node.NodesConfig{
		Hosts:   []string{commandArgs.nodeAddress},
		Retries: commandArgs.retries,
	})
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}
	defer forwarder.Stop()

	ctx, cfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cfunc()

	blockHeight, err := forwarder.LastBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get last block height: %w", err)
	}

	log.Info(fmt.Sprintf("last block height found: %d", blockHeight))

	tx, err := handler.SignTx(commandArgs.name, wReq, blockHeight)
	if err != nil {
		return fmt.Errorf("couldn't sign transaction: %w", err)
	}

	log.Info("transaction successfully signed", zap.String("signature", tx.Signature.Value))

	if err = forwarder.SendTx(ctx, tx, api.SubmitTransactionRequest_TYPE_ASYNC); err != nil {
		return fmt.Errorf("couldn't send transaction: %w", err)
	}

	log.Info("transaction successfully sent")

	return nil
}
