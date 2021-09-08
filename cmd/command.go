package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/protos/commands"
	"code.vegaprotocol.io/protos/vega/api"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	commandArgs struct {
		name           string
		passphraseFile string
		nodeAddress    string
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
	commandCmd.Flags().StringVarP(&commandArgs.nodeAddress, "node-address", "", "0.0.0.0:3002", `Address of the vega node to use"`)
	commandCmd.MarkFlagRequired("name")
	commandCmd.MarkFlagRequired("pubkey")
}

func runCommand(_ *cobra.Command, pos []string) error {
	if len(pos) != 1 {
		return errors.New("invalid number of arguments, require at most 1 command to be signed and sent by the wallet")
	}

	command := walletpb.SubmitTransactionRequest{}
	err := jsonpb.UnmarshalString(pos[0], &command)
	if err != nil {
		return fmt.Errorf("invalid command input: %w", err)
	}

	command.PubKey = commandArgs.pubKey

	errs := wallet.CheckSubmitTransactionRequest(&command)
	if !errs.Empty() {
		return fmt.Errorf("invalid command payload: %w", err)
	}

	passphrase, err := getPassphrase(importArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	store, err := newWalletsStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	w, err := store.GetWallet(commandArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not open wallet: %w", err)
	}

	clt, err := getClient(commandArgs.nodeAddress)
	if err != nil {
		return fmt.Errorf("could not connect to vega node: %w", err)
	}

	height, err := getHeight(clt)
	if err != nil {
		return err
	}

	tx, err := signTx(w, &command, height)
	if err != nil {
		return fmt.Errorf("could not sign the transaction: %w", err)
	}

	return sendTx(clt, tx)
}

func sendTx(clt api.TradingServiceClient, tx *commandspb.Transaction) error {
	ctx, cfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cfunc()
	req := api.SubmitTransactionV2Request{
		Tx:   tx,
		Type: api.SubmitTransactionV2Request_TYPE_ASYNC,
	}
	_, err := clt.SubmitTransactionV2(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}
	return nil
}

func getHeight(clt api.TradingServiceClient) (uint64, error) {
	ctx, cfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cfunc()
	resp, err := clt.LastBlockHeight(ctx, &api.LastBlockHeightRequest{})
	if err != nil {
		return 0, fmt.Errorf("could not get last block: %w", err)
	}

	return resp.Height, nil
}

func signTx(w wallet.Wallet, req *walletpb.SubmitTransactionRequest, height uint64) (*commandspb.Transaction, error) {
	data := commands.NewInputData(height)
	wallet.WrapRequestCommandIntoInputData(data, req)
	marshalledData, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}

	pubKey := req.GetPubKey()
	signature, err := w.SignTxV2(pubKey, marshalledData)
	if err != nil {
		return nil, err
	}

	return commands.NewTransaction(pubKey, marshalledData, signature), nil
}

func getClient(address string) (api.TradingServiceClient, error) {
	tdconn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return api.NewTradingServiceClient(tdconn), nil
}
