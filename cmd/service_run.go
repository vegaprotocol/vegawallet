package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	vgterm "code.vegaprotocol.io/shared/libs/term"
	vglog "code.vegaprotocol.io/shared/libs/zap"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"code.vegaprotocol.io/vegawallet/node"
	"code.vegaprotocol.io/vegawallet/proxy"
	"code.vegaprotocol.io/vegawallet/service"
	svcstore "code.vegaprotocol.io/vegawallet/service/store/v1"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/golang/protobuf/jsonpb"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const MaxConsentRequests = 100

var ErrEnableAutomaticConsentFlagIsRequiredWithoutTTY = errors.New("--automatic-consent flag is required without TTY")

var (
	ErrProgramIsNotInitialised = errors.New("first, you need initialise the program, using the `init` command")

	runServiceLong = cli.LongDesc(`
		Start a Vega wallet service behind an HTTP server.

		By default, every incoming transactions will have to be reviewed in the
		terminal.

		To terminate the service, hit ctrl+c. 

		NOTE: The --output flag is ignored in this command.
	`)

	runServiceExample = cli.Examples(`
		# Start the service
		vegawallet service run --network NETWORK

		# Start the service and open the console in the default browser
		vegawallet service run --network NETWORK --with-console

		# Start the service without opening the console
		vegawallet service run --network NETWORK --with-console --no-browser

		# Start the service and open the token dApp in the default browser
		vegawallet service run --network NETWORK --with-token-dapp

		# Start the service without opening the token dApp
		vegawallet service run --network NETWORK --with-token-dapp --no-browser

		# Start the service with automatic consent of incoming transactions
		vegawallet service run --network NETWORK --automatic-consent
	`)
)

type RunServiceHandler func(io.Writer, *RootFlags, *RunServiceFlags) error

func NewCmdRunService(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdRunService(w, RunService, rf)
}

func BuildCmdRunService(w io.Writer, handler RunServiceHandler, rf *RootFlags) *cobra.Command {
	f := &RunServiceFlags{}

	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Start the Vega wallet service",
		Long:    runServiceLong,
		Example: runServiceExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := f.Validate(); err != nil {
				return err
			}

			if err := handler(w, rf, f); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network configuration to use",
	)
	cmd.Flags().BoolVar(&f.WithConsole,
		"with-console",
		false,
		"Start the Vega console behind a proxy and open it in the default browser",
	)
	cmd.Flags().BoolVar(&f.WithTokenDApp,
		"with-token-dapp",
		false,
		"Start the Vega Token dApp behind a proxy and open it in the default browser",
	)
	cmd.Flags().BoolVar(&f.NoBrowser,
		"no-browser",
		false,
		"Do not open the default browser when starting applications",
	)
	cmd.Flags().BoolVar(&f.EnableAutomaticConsent,
		"automatic-consent",
		false,
		"Automatically approve incoming transaction. Only use this flag when you have absolute trust in incoming transactions! No logs on standard output.",
	)

	autoCompleteNetwork(cmd, rf.Home)

	return cmd
}

type RunServiceFlags struct {
	Network                string
	WithConsole            bool
	WithTokenDApp          bool
	NoBrowser              bool
	EnableAutomaticConsent bool
}

func (f *RunServiceFlags) Validate() error {
	if len(f.Network) == 0 {
		return flags.FlagMustBeSpecifiedError("network")
	}

	if f.NoBrowser && !f.WithConsole && !f.WithTokenDApp {
		return flags.OneOfParentsFlagMustBeSpecifiedError("no-browser", "with-console", "with-token-dapp")
	}

	return nil
}

func RunService(w io.Writer, rf *RootFlags, f *RunServiceFlags) error {
	p := printer.NewInteractivePrinter(w)

	store, err := wallets.InitialiseStore(rf.Home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	vegaPaths := paths.New(rf.Home)
	netStore, err := netstore.InitialiseStore(vegaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	exists, err := netStore.NetworkExists(f.Network)
	if err != nil {
		return fmt.Errorf("couldn't verify the network existence: %w", err)
	}
	if !exists {
		return network.NewNetworkDoesNotExistError(f.Network)
	}

	cfg, err := netStore.GetNetwork(f.Network)
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	if err := verifyNetworkConfig(cfg, f); err != nil {
		return err
	}
	svcLog, svcLogPath, err := BuildJSONLogger(cfg.Level.String(), vegaPaths, paths.WalletServiceLogsHome)
	if err != nil {
		return err
	}
	defer vglog.Sync(svcLog)

	svcLog = svcLog.Named("service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svcStore, err := svcstore.InitialiseStore(paths.New(rf.Home))
	if err != nil {
		return fmt.Errorf("couldn't initialise service store: %w", err)
	}

	if isInit, err := service.IsInitialised(svcStore); err != nil {
		return fmt.Errorf("couldn't verify service initialisation state: %w", err)
	} else if !isInit {
		return ErrProgramIsNotInitialised
	}

	auth, err := service.NewAuth(svcLog.Named("auth"), svcStore, cfg.TokenExpiry.Get())
	if err != nil {
		return fmt.Errorf("couldn't initialise authentication: %w", err)
	}

	forwarder, err := node.NewForwarder(svcLog.Named("forwarder"), cfg.API.GRPC)
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}

	cliLog, cliLogPath, err := BuildJSONLogger(cfg.Level.String(), vegaPaths, paths.WalletCLILogsHome)
	if err != nil {
		return err
	}
	defer vglog.Sync(cliLog)

	cliLog = cliLog.Named("command")

	consentRequests := make(chan service.ConsentRequest, MaxConsentRequests)
	defer close(consentRequests)
	sentTransactions := make(chan service.SentTransaction)
	defer close(sentTransactions)

	var policy service.Policy
	if vgterm.HasTTY() {
		cliLog.Info("TTY detected")
		if f.EnableAutomaticConsent {
			cliLog.Info("Automatic consent enabled")
			policy = service.NewAutomaticConsentPolicy()
		} else {
			cliLog.Info("Explicit consent enabled")
			policy = service.NewExplicitConsentPolicy(ctx, consentRequests, sentTransactions)
		}
	} else {
		cliLog.Info("No TTY detected")
		if !f.EnableAutomaticConsent {
			cliLog.Error("Explicit consent can't be used when no TTY is attached to the process")
			return ErrEnableAutomaticConsentFlagIsRequiredWithoutTTY
		}
		cliLog.Info("Automatic consent enabled")
		policy = service.NewAutomaticConsentPolicy()
	}

	srv, err := service.NewService(svcLog.Named("api"), cfg, handler, auth, forwarder, policy)
	if err != nil {
		return err
	}

	go func() {
		defer cancel()
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cliLog.Error("Error while starting HTTP server", zap.Error(err))
		}
	}()

	serviceHost := fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port)
	if !f.EnableAutomaticConsent {
		p.CheckMark().Text("HTTP service started at: ").SuccessText(serviceHost).NextLine()
	}
	cliLog.Info(fmt.Sprintf("HTTP service started at: %s", serviceHost))

	defer func() {
		if err = srv.Stop(); err != nil {
			cliLog.Error("Error while stopping HTTP server", zap.Error(err))
			cliLog.Info("HTTP server stopped with success")
		}
	}()

	var cs *proxy.Proxy
	if f.WithConsole {
		cs = startConsole(cliLog, f, cfg, cancel, p)
		defer func() {
			if err = cs.Stop(); err != nil {
				cliLog.Error("Error while stopping console proxy", zap.Error(err))
			} else {
				cliLog.Info("Console proxy stopped with success")
			}
		}()
	}

	var tokenDApp *proxy.Proxy
	if f.WithTokenDApp {
		tokenDApp = startTokenDApp(cliLog, f, cfg, cancel, p)
		defer func() {
			if err = tokenDApp.Stop(); err != nil {
				cliLog.Error("Error while stopping token dApp proxy", zap.Error(err))
			} else {
				cliLog.Info("Token dApp proxy stopped with success")
			}
		}()
	}

	if !f.EnableAutomaticConsent {
		p.CheckMark().Text("Service logs located at: ").SuccessText(svcLogPath).NextLine()
		p.CheckMark().Text("CLI logs located at: ").SuccessText(cliLogPath).NextLine()
		p.CheckMark().SuccessText("Starting successful").NextSection()
		p.NextLine()
	}

	waitSig(ctx, cancel, cliLog, consentRequests, sentTransactions, p)

	return nil
}

func verifyNetworkConfig(cfg *network.Network, f *RunServiceFlags) error {
	if err := cfg.EnsureCanConnectGRPCNode(); err != nil {
		return err
	}
	if f.WithConsole {
		if err := cfg.EnsureCanConnectConsole(); err != nil {
			return err
		}
	}
	if f.WithTokenDApp {
		if err := cfg.EnsureCanConnectTokenDApp(); err != nil {
			return err
		}
	}
	return nil
}

func startConsole(log *zap.Logger, f *RunServiceFlags, cfg *network.Network, cancel context.CancelFunc, p *printer.InteractivePrinter) *proxy.Proxy {
	cs := proxy.NewProxy(cfg.Console.LocalPort, cfg.Console.URL, cfg.API.GRPC.Hosts[0])
	go func() {
		defer cancel()
		if err := cs.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error while starting the console proxy", zap.Error(err))
		}
	}()

	consoleLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.Console.LocalPort)
	if !f.EnableAutomaticConsent {
		p.CheckMark().Text("Console proxy pointing to ").Bold(cfg.Console.URL).Text(" started at: ").SuccessText(consoleLocalHost).NextLine()
	}
	log.Info(fmt.Sprintf("console proxy pointing to %s started at: %s", cfg.Console.URL, consoleLocalHost))

	if !f.NoBrowser {
		if err := open.Run(cs.GetBrowserURL()); err != nil {
			log.Error("unable to open the application in the default browser", zap.Error(err))
		}
	}

	return cs
}

func startTokenDApp(log *zap.Logger, f *RunServiceFlags, cfg *network.Network, cancel context.CancelFunc, p *printer.InteractivePrinter) *proxy.Proxy {
	tokenDApp := proxy.NewProxy(cfg.TokenDApp.LocalPort, cfg.TokenDApp.URL, cfg.API.GRPC.Hosts[0])
	go func() {
		defer cancel()
		if err := tokenDApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error while starting the token dApp proxy", zap.Error(err))
		}
	}()

	tokenDAppLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.TokenDApp.LocalPort)
	if !f.EnableAutomaticConsent {
		p.CheckMark().Text("token dApp proxy pointing to ").Bold(cfg.TokenDApp.URL).Text(" started at: ").SuccessText(tokenDAppLocalHost).NextLine()
	}
	log.Info(fmt.Sprintf("token dApp proxy pointing to %s started at: %s", cfg.TokenDApp.URL, tokenDAppLocalHost))

	if !f.NoBrowser {
		if err := open.Run(tokenDApp.GetBrowserURL()); err != nil {
			log.Error("unable to open the token dApp in the default browser", zap.Error(err))
		}
	}
	return tokenDApp
}

// waitSig will wait for a sigterm or sigint interrupt.
func waitSig(
	ctx context.Context,
	cancelFunc context.CancelFunc,
	log *zap.Logger,
	consentRequests chan service.ConsentRequest,
	sentTransactions chan service.SentTransaction,
	p *printer.InteractivePrinter,
) {
	gracefulStop := make(chan os.Signal, 1)
	defer close(gracefulStop)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	go func() {
		if err := handleConsentRequests(ctx, log, consentRequests, sentTransactions, p); err != nil {
			cancelFunc()
		}
	}()

	for {
		select {
		case sig := <-gracefulStop:
			log.Info("caught signal", zap.String("signal", fmt.Sprintf("%+v", sig)))
			cancelFunc()
			return
		case <-ctx.Done():
			// nothing to do
			return
		}
	}
}

func handleConsentRequests(ctx context.Context, log *zap.Logger, consentRequests chan service.ConsentRequest, sentTransactions chan service.SentTransaction, p *printer.InteractivePrinter) error {
	for {
		select {
		case <-ctx.Done():
			// nothing to do
			return nil
		case consentRequest := <-consentRequests:
			m := jsonpb.Marshaler{Indent: "    "}
			marshalledTx, err := m.MarshalToString(consentRequest.Tx)
			if err != nil {
				log.Error("couldn't marshal transaction from consent request", zap.Error(err))
				return err
			}

			p.BlueArrow().Text("New transaction received: ").NextLine()
			p.InfoText(marshalledTx).NextLine()

			if flags.DoYouApproveTx() {
				log.Info("user approved the signing of the transaction", zap.Any("transaction", marshalledTx))
				consentRequest.Confirmation <- service.ConsentConfirmation{Decision: true}
				p.CheckMark().SuccessText("Transaction approved").NextLine()

				sentTx := <-sentTransactions
				log.Info("transaction sent", zap.Any("ID", sentTx.TxID), zap.Any("hash", sentTx.TxHash))
				if sentTx.Error != nil {
					log.Error("transaction failed", zap.Any("transaction", marshalledTx))
					p.BangMark().DangerText("Transaction failed").NextLine()
					p.BangMark().DangerText("Error: ").DangerText(sentTx.Error.Error()).NextSection()
				} else {
					log.Info("transaction sent", zap.Any("hash", sentTx.TxHash))
					p.CheckMark().Text("Transaction with hash ").SuccessText(sentTx.TxHash).Text(" sent!").NextSection()
				}
			} else {
				log.Info("user rejected the signing of the transaction", zap.Any("transaction", marshalledTx))
				consentRequest.Confirmation <- service.ConsentConfirmation{Decision: false}
				p.BangMark().DangerText("Transaction rejected").NextSection()
			}
		}
	}
}
