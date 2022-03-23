package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ErrProgramIsNotInitialised = errors.New("first, you need initialise the program, using the `init` command")

	runServiceLong = cli.LongDesc(`
		Start a Vega wallet service behind an HTTP server.

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

	autoCompleteNetwork(cmd, rf.Home)

	return cmd
}

type RunServiceFlags struct {
	Network       string
	WithConsole   bool
	WithTokenDApp bool
	NoBrowser     bool
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

	svcLog, svcLogPath, err := BuildJSONLogger(cfg.Level.String(), paths.WalletServiceLogsHome)
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

	pendingConsents := make(chan service.ConsentRequest, 1)
	consentConfirmations := make(chan service.ConsentConfirmation, 1)

	var policy service.Policy
	switch rf.Output {
	case flags.InteractiveOutput:
		policy = service.NewExplicitConsentPolicy(pendingConsents, consentConfirmations)
	case flags.JSONOutput:
		policy = service.NewAutomaticConsentPolicy(pendingConsents, consentConfirmations)
	}

	srv, err := service.NewService(svcLog.Named("api"), cfg, handler, auth, forwarder, policy)
	if err != nil {
		return err
	}

	cliLog, cliLogPath, err := BuildJSONLogger(cfg.Level.String(), paths.WalletCLILogsHome)
	if err != nil {
		return err
	}
	defer vglog.Sync(cliLog)

	cliLog = cliLog.Named("command")

	go func() {
		defer cancel()
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cliLog.Error("Error while starting HTTP server", zap.Error(err))
		}
	}()

	serviceHost := fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port)
	if rf.Output == flags.InteractiveOutput {
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
		cs = startConsole(cliLog, rf, !f.NoBrowser, cfg, cancel, p)
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
		tokenDApp = startTokenDApp(cliLog, rf, !f.NoBrowser, cfg, cancel, p)
		defer func() {
			if err = tokenDApp.Stop(); err != nil {
				cliLog.Error("Error while stopping token dApp proxy", zap.Error(err))
			} else {
				cliLog.Info("Token dApp proxy stopped with success")
			}
		}()
	}

	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().SuccessText("Starting successful").NextSection()
		p.CheckMark().SuccessText("Service logs output to: ").Bold(svcLogPath).NextSection()
		p.CheckMark().SuccessText("CLI logs output to: ").Bold(cliLogPath).NextSection()
		p.NextLine()
	}

	waitSig(ctx, cancel, cliLog, pendingConsents, consentConfirmations, p)

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

func startConsole(log *zap.Logger, rf *RootFlags, openBrowser bool, cfg *network.Network, cancel context.CancelFunc, p *printer.InteractivePrinter) *proxy.Proxy {
	cs := proxy.NewProxy(cfg.Console.LocalPort, cfg.Console.URL, cfg.API.GRPC.Hosts[0])
	go func() {
		defer cancel()
		if err := cs.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error while starting the console proxy", zap.Error(err))
		}
	}()

	consoleLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.Console.LocalPort)
	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().Text("Console proxy pointing to ").Bold(cfg.Console.URL).Text(" started at: ").SuccessText(consoleLocalHost).NextLine()
	}
	log.Info(fmt.Sprintf("console proxy pointing to %s started at: %s", cfg.Console.URL, consoleLocalHost))

	if openBrowser {
		if err := open.Run(cs.GetBrowserURL()); err != nil {
			log.Error("unable to open the application in the default browser", zap.Error(err))
		}
	}

	return cs
}

func startTokenDApp(log *zap.Logger, rf *RootFlags, openBrowser bool, cfg *network.Network, cancel context.CancelFunc, p *printer.InteractivePrinter) *proxy.Proxy {
	tokenDApp := proxy.NewProxy(cfg.TokenDApp.LocalPort, cfg.TokenDApp.URL, cfg.API.GRPC.Hosts[0])
	go func() {
		defer cancel()
		if err := tokenDApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error while starting the token dApp proxy", zap.Error(err))
		}
	}()

	tokenDAppLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.TokenDApp.LocalPort)
	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().Text("token dApp proxy pointing to ").Bold(cfg.TokenDApp.URL).Text(" started at: ").SuccessText(tokenDAppLocalHost).NextLine()
	}
	log.Info(fmt.Sprintf("token dApp proxy pointing to %s started at: %s", cfg.TokenDApp.URL, tokenDAppLocalHost))

	if openBrowser {
		if err := open.Run(tokenDApp.GetBrowserURL()); err != nil {
			log.Error("unable to open the token dApp in the default browser", zap.Error(err))
		}
	}
	return tokenDApp
}

// waitSig will wait for a sigterm or sigint interrupt.
func waitSig(ctx context.Context, cfunc func(), log *zap.Logger, pendingSigRequests chan service.ConsentRequest, sigRequestsResponses chan service.ConsentConfirmation, p *printer.InteractivePrinter) {
	gracefulStop := make(chan os.Signal, 1)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	for {
		select {
		case sig := <-gracefulStop:
			log.Info("caught signal", zap.String("signal", fmt.Sprintf("%+v", sig)))
			cfunc()
			return
		case <-ctx.Done():
			// nothing to do
			return
		case ev := <-pendingSigRequests:
			txStr := ev.String()
			p.CheckMark().Text("Received TX sign request: ").WarningText(txStr).NextLine()
			reader := bufio.NewReader(os.Stdin)
			p.CheckMark().WarningText("Please accept or decline sign request: (y/n)").NextLine()
			answer, err := reader.ReadString('\n')
			if err != nil {
				log.Info("failed to read user input")
				cfunc()
				return
			}
			if answer == "y" || answer == "Y" {
				log.Info("user approved signature for transaction", zap.String("transaction", txStr))
				sigRequestsResponses <- service.ConsentConfirmation{Decision: true, TxStr: txStr}
				p.CheckMark().WarningText("Sign request accepted").NextLine()
			} else {
				log.Info("user declined signature for transaction", zap.String("transaction", txStr))
				sigRequestsResponses <- service.ConsentConfirmation{Decision: false, TxStr: txStr}
				p.CheckMark().WarningText("Sign request rejected").NextLine()
			}
		}
	}
}
