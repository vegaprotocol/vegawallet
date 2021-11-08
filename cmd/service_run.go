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
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/console"
	vglog "code.vegaprotocol.io/vegawallet/libs/zap"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"code.vegaprotocol.io/vegawallet/node"
	"code.vegaprotocol.io/vegawallet/service"
	svcstore "code.vegaprotocol.io/vegawallet/service/store/v1"
	"code.vegaprotocol.io/vegawallet/wallets"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const startupT = ` # Authentication
 - login:                   POST   {{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:                  DELETE {{.WalletServiceLocalAddress}}/api/v1/auth/token

 # Network management
 - network:                 GET    {{.WalletServiceLocalAddress}}/api/v1/network

 # Wallet management
 - create a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets
 - import a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets/import

 # Key pair management
 - generate a key pair:     POST   {{.WalletServiceLocalAddress}}/api/v1/keys
 - list keys:               GET    {{.WalletServiceLocalAddress}}/api/v1/keys
 - describe a key pair:     GET    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid
 - taint a key pair:        PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/taint
 - annotate a key pair:     PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/metadata

 # Commands
 - sign a command:          POST   {{.WalletServiceLocalAddress}}/api/v1/command
 - sign a command (sync):   POST   {{.WalletServiceLocalAddress}}/api/v1/command/sync
 - sign a command (commit): POST   {{.WalletServiceLocalAddress}}/api/v1/command/commit
 - sign data:               POST   {{.WalletServiceLocalAddress}}/api/v1/sign
 - verify data:             POST   {{.WalletServiceLocalAddress}}/api/v1/verify

 # Information
 - get service status:      GET    {{.WalletServiceLocalAddress}}/api/v1/status
 - get the version:         GET    {{.WalletServiceLocalAddress}}/api/v1/version
`

var (
	runServiceLong = cli.LongDesc(`
		Start a Vega wallet service behind an HTTP server.

		To terminate the service, hit ctrl+c.
	`)

	runServiceExample = cli.Examples(`
		# Start the service
		vegawallet service run --network NETWORK

		# Start the service with a log level set to debug
		vegawallet service run --network NETWORK --level debug

		# Start the service with the console proxy and open the console in the 
		# default browser
		vegawallet service run --network NETWORK --console-proxy

		# Start the service with the console proxy without opening the console
		vegawallet service run --network NETWORK --console-proxy --no-browser
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
	cmd.Flags().BoolVar(&f.StartConsole,
		"console-proxy",
		false,
		"Start the Vega console proxy and open the console in the default browser")
	cmd.Flags().BoolVar(&f.NoBrowser,
		"no-browser",
		false,
		"Do not open the default browser when starting the console proxy (requires: --console-proxy)")
	cmd.Flags().StringVar(&f.LogLevel,
		"level",
		"",
		fmt.Sprintf("Set the log level: %v (default: value set by the network configuration)", SupportedLogLevels))

	return cmd
}

type RunServiceFlags struct {
	Network      string
	StartConsole bool
	NoBrowser    bool
	LogLevel     string
}

func (f *RunServiceFlags) Validate() error {
	if len(f.Network) == 0 {
		return flags.FlagMustBeSpecifiedError("network")
	}

	if f.NoBrowser && !f.StartConsole {
		return flags.ParentFlagMustBeSpecifiedError("no-browser", "console-proxy")
	}

	if len(f.LogLevel) != 0 {
		if err := ValidateLogLevel(f.LogLevel); err != nil {
			return err
		}
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

	netStore, err := netstore.InitialiseStore(paths.New(rf.Home))
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

	logLevel := cfg.Level.String()
	if len(f.LogLevel) != 0 {
		logLevel = f.LogLevel
	}
	log, err := Build(rf.Output, logLevel)
	if err != nil {
		return err
	}
	defer vglog.Sync(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svcStore, err := svcstore.InitialiseStore(paths.New(rf.Home))
	if err != nil {
		return fmt.Errorf("couldn't initialise service store: %w", err)
	}

	auth, err := service.NewAuth(log.Named("auth"), svcStore, cfg.TokenExpiry.Get())
	if err != nil {
		return fmt.Errorf("couldn't initialise authentication: %w", err)
	}

	forwarder, err := node.NewForwarder(log.Named("forwarder"), cfg.API.GRPC)
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}

	srv, err := service.NewService(log.Named("service"), cfg, handler, auth, forwarder)
	if err != nil {
		return err
	}

	log = log.Named("command")

	go func() {
		defer cancel()
		err := srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error while starting HTTP server", zap.Error(err))
		}
	}()

	serviceHost := fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port)
	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().Text("HTTP service started at: ").SuccessText(serviceHost).NextLine()
	} else if rf.Output == "json" {
		log.Info(fmt.Sprintf("HTTP service started at: %s", serviceHost))
	}

	var cs *console.Console
	if f.StartConsole {
		cs = console.NewConsole(cfg.Console.LocalPort, cfg.Console.URL, cfg.API.GRPC.Hosts[0])
		go func() {
			defer cancel()
			err := cs.Start()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("error while starting the console proxy", zap.Error(err))
			}
		}()

		consoleLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.Console.LocalPort)
		if rf.Output == flags.InteractiveOutput {
			p.CheckMark().Text("Console proxy pointing to ").Bold(cfg.Console.URL).Text(" started at: ").SuccessText(consoleLocalHost).NextLine()
		} else if rf.Output == "json" {
			log.Info(fmt.Sprintf("console proxy pointing to %s started at: %s", cfg.Console.URL, consoleLocalHost))
		}

		if !f.NoBrowser {
			err := open.Run(cs.GetBrowserURL())
			if err != nil {
				log.Error("unable to open the console in the default browser",
					zap.Error(err))
			}
		}
	}

	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().SuccessText("Starting successful").NextSection()
		p.BlueArrow().InfoText("Available endpoints").NextLine()
		printServiceEndpoints(serviceHost)
		p.NextSection()
		p.BlueArrow().InfoText("Logs").NextLine()
	}

	waitSig(ctx, cancel, log)

	err = srv.Stop()
	if err != nil {
		log.Error("error while stopping HTTP server", zap.Error(err))
	} else {
		log.Info("HTTP server stopped with success")
	}

	if f.StartConsole {
		err = cs.Stop()
		if err != nil {
			log.Error("error while stopping console proxy", zap.Error(err))
		} else {
			log.Info("console proxy stopped with success")
		}
	}

	if rf.Output == flags.InteractiveOutput {
		p.CheckMark().SuccessText("Service stopped").NextLine()
	}

	return nil
}

// waitSig will wait for a sigterm or sigint interrupt.
func waitSig(ctx context.Context, cfunc func(), log *zap.Logger) {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	select {
	case sig := <-gracefulStop:
		log.Info("caught signal", zap.String("signal", fmt.Sprintf("%+v", sig)))
		cfunc()
	case <-ctx.Done():
		// nothing to do
	}
}

func printServiceEndpoints(serviceHost string) {
	params := struct {
		WalletServiceLocalAddress string
	}{
		WalletServiceLocalAddress: serviceHost,
	}

	tmpl, err := template.New("wallet-cmdline").Parse(startupT)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		panic(err)
	}
}
