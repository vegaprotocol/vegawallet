package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	wcommands "code.vegaprotocol.io/go-wallet/commands"
	vglog "code.vegaprotocol.io/go-wallet/libs/zap"
	"code.vegaprotocol.io/go-wallet/logger"
	"code.vegaprotocol.io/go-wallet/network"
	netstore "code.vegaprotocol.io/go-wallet/network/store/v1"
	"code.vegaprotocol.io/go-wallet/node"
	"code.vegaprotocol.io/go-wallet/wallets"
	api "code.vegaprotocol.io/protos/vega/api/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/shared/paths"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var (
	commandArgs struct {
		network        string
		wallet         string
		passphraseFile string
		nodeAddress    string
		retries        uint64
		pubKey         string
	}

	commandCmd = &cobra.Command{
		Use:   "command",
		Short: "Send a command to the Vega network",
		Long:  "Send a command to the Vega network.",
		Args:  cobra.ExactArgs(1),
		RunE:  runCommand,
	}
)

func init() {
	rootCmd.AddCommand(commandCmd)
	commandCmd.Flags().StringVarP(&serviceRunArgs.network, "network", "n", "", "Name of the network to use")
	commandCmd.Flags().StringVarP(&commandArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	commandCmd.Flags().StringVarP(&commandArgs.pubKey, "pubkey", "", "", "The public key to use from the wallet")
	commandCmd.Flags().StringVarP(&commandArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	commandCmd.Flags().StringVar(&commandArgs.nodeAddress, "node-address", "0.0.0.0:3002", "Address of the Vega node to use")
	commandCmd.Flags().Uint64Var(&commandArgs.retries, "retries", 5, "Number of retries when contacting the Vega node")
	_ = commandCmd.MarkFlagRequired("wallet")
	_ = commandCmd.MarkFlagRequired("pubkey")
}

func runCommand(_ *cobra.Command, pos []string) error {
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

	err = handler.LoginWallet(commandArgs.wallet, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't login to the wallet %s: %w", commandArgs.wallet, err)
	}
	defer handler.LogoutWallet(commandArgs.wallet)

	encoding := "json"
	if rootArgs.output == "human" {
		encoding = "console"
	}

	log, err := logger.New(zapcore.InfoLevel, encoding)
	if err != nil {
		return fmt.Errorf("couldn't create logger: %w", err)
	}
	defer vglog.Sync(log)

	if len(commandArgs.nodeAddress) != 0 && len(commandArgs.network) != 0 {
		return errors.New("can't have both node address and network flag set")
	}

	var hosts []string
	if len(commandArgs.nodeAddress) != 0 {
		hosts = []string{commandArgs.nodeAddress}
	} else if len(commandArgs.network) != 0 {
		netStore, err := netstore.InitialiseStore(paths.New(rootArgs.home))
		if err != nil {
			return fmt.Errorf("couldn't initialise network store: %w", err)
		}
		exists, err := netStore.NetworkExists(commandArgs.network)
		if err != nil {
			return fmt.Errorf("couldn't verify network existance: %w", err)
		}
		if !exists {
			return fmt.Errorf("network %s does not exist", commandArgs.network)
		}
		net, err := netStore.GetNetwork(commandArgs.network)
		if err != nil {
			return fmt.Errorf("couldn't get network %s: %w", commandArgs.network, err)
		}
		hosts = net.API.GRPC.Hosts
	} else {
		return errors.New("should set node address or network flag")
	}

	forwarder, err := node.NewForwarder(log.Named("forwarder"), network.GRPCConfig{
		Hosts:   hosts,
		Retries: commandArgs.retries,
	})
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}
	defer func() {
		// We can ignore this non-blocking error without logging as it's already
		// logged down stream.
		_ = forwarder.Stop()
	}()

	ctx, cfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cfunc()

	blockHeight, err := forwarder.LastBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get last block height: %w", err)
	}

	log.Info(fmt.Sprintf("last block height found: %d", blockHeight))

	tx, err := handler.SignTx(commandArgs.wallet, wReq, blockHeight)
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
