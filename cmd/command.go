package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	api "code.vegaprotocol.io/protos/vega/api/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	wcommands "code.vegaprotocol.io/vegawallet/commands"
	vglog "code.vegaprotocol.io/vegawallet/libs/zap"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"code.vegaprotocol.io/vegawallet/node"
	"code.vegaprotocol.io/vegawallet/wallets"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

const (
	DefaultForwarderRetryCount = 5
	ForwarderRequestTimeout    = 5 * time.Second
)

var (
	ErrNetworkDoesNotHaveGRPCHostConfigured = errors.New("network does not have gRPC hosts configured")
	ErrDoNotSetPubKeyInRequest              = errors.New("do not set the public key through the request, use --pubkey flag instead")

	sendCommandLong = cli.LongDesc(`
		Send a command to a Vega node via the gRPC API. The command can be sent to 
		any node of a registered network or to a specific node address.

		The command should be a Vega command formatted as a JSON payload, as follows:

		'{"commandName": {"someProperty": "someValue"} }'

		For vote submission, it will look like this:

		'{"voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}'
	`)

	sendCommandExample = cli.Examples(`
		# Send a command to a registered network
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY REQUEST

		# Send a command to a specific Vega node address
		vegawallet command --node-address ADDRESS --wallet WALLET --pubkey PUBKEY REQUEST

		# Send a command with a log level set to debug
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY --level debug REQUEST

		# Send a command with a maximum of 10 retry
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY --retries 10 REQUEST
	`)
)

type SendCommandHandler func(io.Writer, *RootFlags, *SendCommandRequest) error

func NewCmdSendCommand(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdSendCommand(w, SendCommand, rf)
}

func BuildCmdSendCommand(w io.Writer, handler SendCommandHandler, rf *RootFlags) *cobra.Command {
	f := &SendCommandFlags{}

	cmd := &cobra.Command{
		Use:     "command",
		Short:   "Send a command to a Vega node",
		Long:    sendCommandLong,
		Example: sendCommandExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			if err := handler(w, rf, req); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network to which the command is sent",
	)
	cmd.Flags().StringVar(&f.NodeAddress,
		"node-address",
		"",
		"Vega node address to which the command is sent",
	)
	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the public key",
	)
	cmd.Flags().StringVarP(&f.PubKey,
		"pubkey", "k",
		"",
		"Public key of the key pair to use for signing (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVar(&f.LogLevel,
		"level",
		zapcore.InfoLevel.String(),
		fmt.Sprintf("Set the log level: %v", SupportedLogLevels),
	)
	cmd.Flags().Uint64Var(&f.Retries,
		"retries",
		DefaultForwarderRetryCount,
		"Number of retries when contacting the Vega node",
	)

	return cmd
}

type SendCommandFlags struct {
	Network        string
	NodeAddress    string
	Wallet         string
	PubKey         string
	PassphraseFile string
	Retries        uint64
	LogLevel       string
	RawRequest     string
}

func (f *SendCommandFlags) Validate() (*SendCommandRequest, error) {
	req := &SendCommandRequest{
		Retries: f.Retries,
	}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	if len(f.LogLevel) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("level")
	}
	if err := ValidateLogLevel(f.LogLevel); err != nil {
		return nil, err
	}
	req.LogLevel = f.LogLevel

	if len(f.NodeAddress) == 0 && len(f.Network) == 0 {
		return nil, flags.OneOfFlagsMustBeSpecifiedError("network", "node-address")
	}
	if len(f.NodeAddress) != 0 && len(f.Network) != 0 {
		return nil, flags.FlagsMutuallyExclusiveError("network", "node-address")
	}
	req.NodeAddress = f.NodeAddress
	req.Network = f.Network

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	if len(f.PubKey) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("pubkey")
	}
	if len(f.RawRequest) == 0 {
		return nil, flags.ArgMustBeSpecifiedError("request")
	}
	request := &walletpb.SubmitTransactionRequest{}
	if err := jsonpb.UnmarshalString(f.RawRequest, request); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal request: %w", err)
	}
	if len(request.PubKey) != 0 {
		return nil, ErrDoNotSetPubKeyInRequest
	}
	request.PubKey = f.PubKey
	request.Propagate = true
	req.Request = request
	errs := wcommands.CheckSubmitTransactionRequest(req.Request)
	if !errs.Empty() {
		return nil, fmt.Errorf("invalid request: %w", errs)
	}

	return req, nil
}

type SendCommandRequest struct {
	Network     string
	NodeAddress string
	Wallet      string
	Passphrase  string
	Retries     uint64
	LogLevel    string
	Request     *walletpb.SubmitTransactionRequest
}

func SendCommand(w io.Writer, rf *RootFlags, req *SendCommandRequest) error {
	log, err := Build(rf.Output, req.LogLevel)
	if err != nil {
		return err
	}
	defer vglog.Sync(log)

	var hosts []string
	if len(req.Network) != 0 {
		hosts, err = getHostsFromNetwork(rf, req)
		if err != nil {
			return err
		}
	} else {
		hosts = []string{req.NodeAddress}
	}

	forwarder, err := node.NewForwarder(log.Named("forwarder"), network.GRPCConfig{
		Hosts:   hosts,
		Retries: req.Retries,
	})
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}
	defer func() {
		// We can ignore this non-blocking error without logging as it's already
		// logged down stream.
		_ = forwarder.Stop()
	}()

	p := printer.NewInteractivePrinter(w)
	if rf.Output == flags.InteractiveOutput {
		p.BlueArrow().InfoText("Logs").NextLine()
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), ForwarderRequestTimeout)
	defer cancelFn()

	log.Info("retrieving block height")
	blockHeight, err := forwarder.LastBlockHeight(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get last block height: %w", err)
	}
	log.Info(fmt.Sprintf("last block height found: %d", blockHeight))

	store, err := wallets.InitialiseStore(rf.Home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}
	handler := wallets.NewHandler(store)
	err = handler.LoginWallet(req.Wallet, req.Passphrase)
	if err != nil {
		return fmt.Errorf("couldn't login to the wallet %s: %w", req.Wallet, err)
	}
	defer handler.LogoutWallet(req.Wallet)
	tx, err := handler.SignTx(req.Wallet, req.Request, blockHeight)
	if err != nil {
		log.Error("couldn't sign transaction", zap.Error(err))
		return fmt.Errorf("couldn't sign transaction: %w", err)
	}

	log.Info("transaction successfully signed", zap.String("signature", tx.Signature.Value))

	if err = forwarder.SendTx(ctx, tx, api.SubmitTransactionRequest_TYPE_ASYNC); err != nil {
		log.Error("couldn't send transaction", zap.Error(err))
		return fmt.Errorf("couldn't send transaction: %w", err)
	}

	log.Info("transaction successfully sent")

	return nil
}

func getHostsFromNetwork(rf *RootFlags, req *SendCommandRequest) ([]string, error) {
	netStore, err := netstore.InitialiseStore(paths.New(rf.Home))
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise network store: %w", err)
	}
	net, err := network.GetNetwork(netStore, req.Network)
	if err != nil {
		return nil, err
	}

	if len(net.API.GRPC.Hosts) == 0 {
		return nil, ErrNetworkDoesNotHaveGRPCHostConfigured
	}

	return net.API.GRPC.Hosts, nil
}
