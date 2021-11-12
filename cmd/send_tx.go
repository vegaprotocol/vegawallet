package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	api "code.vegaprotocol.io/protos/vega/api/v1"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	vglog "code.vegaprotocol.io/vegawallet/libs/zap"
	"code.vegaprotocol.io/vegawallet/network"
	"code.vegaprotocol.io/vegawallet/node"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/cobra"
)

var (
	sendTxLong = cli.LongDesc(`
		Send a transaction to a Vega node via the gRPC API. The command can be sent to 
		any node of a registered network or to a specific node address.

		The transaction is base64-encoded.
	`)

	sendTxExample = cli.Examples(`
		# Send a command to a registered network
		vegawallet send tx --network NETWORK BASE64_TRANSACTION

		# Send a command to a specific Vega node address
		vegawallet send tx --node-address ADDRESS BASE64_TRANSACTION

		# Send a command with a log level set to debug
		vegawallet send tx --network NETWORK --level debug BASE64_TRANSACTION

		# Send a command with a maximum of 10 retry
		vegawallet send tx --network NETWORK --retries 10 BASE64_TRANSACTION
	`)
)

type SendTxHandler func(io.Writer, *RootFlags, *SendTxRequest) error

func NewCmdSendTx(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdSendTx(w, SendTx, rf)
}

func BuildCmdSendTx(w io.Writer, handler SendTxHandler, rf *RootFlags) *cobra.Command {
	f := &SendTxFlags{}

	cmd := &cobra.Command{
		Use:     "command",
		Short:   "Send a transaction to a Vega node",
		Long:    sendTxLong,
		Example: sendTxExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if aLen := len(args); aLen == 0 {
				return flags.ArgMustBeSpecifiedError("transaction")
			} else if aLen > 1 {
				return flags.TooManyArgsError("transaction")
			}
			f.RawTx = args[0]

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

type SendTxFlags struct {
	Network     string
	NodeAddress string
	Retries     uint64
	LogLevel    string
	RawTx       string
}

func (f *SendTxFlags) Validate() (*SendTxRequest, error) {
	req := &SendTxRequest{
		Retries: f.Retries,
	}

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

	if len(f.RawTx) == 0 {
		return nil, flags.ArgMustBeSpecifiedError("transaction")
	}
	decodedTx, err := base64.StdEncoding.DecodeString(f.RawTx)
	if err != nil {
		return nil, flags.MustBase64EncodedError("transaction")
	}
	tx := &commandspb.Transaction{}
	if err := proto.Unmarshal(decodedTx, tx); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal transaction: %w", err)
	}
	req.Tx = tx

	return req, nil
}

type SendTxRequest struct {
	Network     string
	NodeAddress string
	Retries     uint64
	LogLevel    string
	Tx          *commandspb.Transaction
}

func SendTx(w io.Writer, rf *RootFlags, req *SendTxRequest) error {
	log, err := Build(rf.Output, req.LogLevel)
	if err != nil {
		return err
	}
	defer vglog.Sync(log)

	var hosts []string
	if len(req.Network) != 0 {
		hosts, err = getHostsFromNetwork(rf, req.Network)
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
		if err = forwarder.Stop(); err != nil {
			log.Warn("couldn't stop the forwarder", zap.Error(err))
		}
	}()

	p := printer.NewInteractivePrinter(w)
	if rf.Output == flags.InteractiveOutput {
		p.BlueArrow().InfoText("Logs").NextLine()
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), ForwarderRequestTimeout)
	defer cancelFn()

	if err = forwarder.SendTx(ctx, req.Tx, api.SubmitTransactionRequest_TYPE_ASYNC); err != nil {
		log.Error("couldn't send transaction", zap.Error(err))
		return fmt.Errorf("couldn't send transaction: %w", err)
	}

	log.Info("transaction successfully sent")

	return nil
}
