package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/console"
	vglog "code.vegaprotocol.io/go-wallet/libs/zap"
	"code.vegaprotocol.io/go-wallet/logger"
	netstore "code.vegaprotocol.io/go-wallet/network/store/v1"
	"code.vegaprotocol.io/go-wallet/node"
	"code.vegaprotocol.io/go-wallet/service"
	svcstore "code.vegaprotocol.io/go-wallet/service/store/v1"
	"code.vegaprotocol.io/go-wallet/wallets"
	"code.vegaprotocol.io/shared/paths"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serviceRunArgs struct {
		startConsole bool
		noBrowser    bool
		network      string
	}

	serviceRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Start the Vega wallet service",
		Long:  "Start a Vega wallet service behind an HTTP server",
		RunE:  runServiceRun,
	}
)

func init() {
	serviceCmd.AddCommand(serviceRunCmd)
	serviceRunCmd.Flags().StringVarP(&serviceRunArgs.network, "network", "n", "", "Name of the network to use")
	serviceRunCmd.Flags().BoolVar(&serviceRunArgs.startConsole, "console-proxy", false, "Start the vega console proxy and open the console in the default browser")
	serviceRunCmd.Flags().BoolVar(&serviceRunArgs.noBrowser, "no-browser", false, "Do not open the default browser if the console proxy is stated")
	_ = serviceRunCmd.MarkFlagRequired("network")
}

func runServiceRun(_ *cobra.Command, _ []string) error {
	p := printer.NewHumanPrinter()

	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	netStore, err := netstore.InitialiseStore(paths.New(rootArgs.home))
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	cfg, err := netStore.GetNetwork(serviceRunArgs.network)
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	encoding := "json"
	if rootArgs.output == "human" {
		encoding = "console"
	}

	log, err := logger.New(cfg.Level.Level, encoding)
	if err != nil {
		return fmt.Errorf("couldn't create logger: %w", err)
	}
	defer vglog.Sync(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svcStore, err := svcstore.InitialiseStore(paths.New(rootArgs.home))
	if err != nil {
		return fmt.Errorf("couldn't initialise service store: %w", err)
	}

	auth, err := service.NewAuth(log.Named("auth"), svcStore, cfg.TokenExpiry.Get())
	if err != nil {
		return fmt.Errorf("couldn't initialise authentication: %w", err)
	}

	forwarder, err := node.NewForwarder(log.Named("forwarder"), cfg.Nodes)
	if err != nil {
		return fmt.Errorf("couldn't initialise the node forwarder: %w", err)
	}

	srv, err := service.NewService(log.Named("service"), cfg, handler, auth, forwarder)
	if err != nil {
		return err
	}
	go func() {
		defer cancel()
		err := srv.Start()
		if err != nil && err != http.ErrServerClosed {
			log.Error("error starting wallet http server", zap.Error(err))
		}
	}()

	serviceHost := fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port)
	if rootArgs.output == "human" {
		p.CheckMark().Text("HTTP service started at: ").SuccessText(serviceHost).Jump()
	} else if rootArgs.output == "json" {
		log.Info(fmt.Sprintf("HTTP service started at: %s", serviceHost))
	}

	var cs *console.Console
	if serviceRunArgs.startConsole {
		cs = console.NewConsole(cfg.Console.LocalPort, cfg.Console.URL, cfg.Nodes.Hosts[0])
		go func() {
			defer cancel()
			err := cs.Start()
			if err != nil && err != http.ErrServerClosed {
				log.Error("error starting console proxy server", zap.Error(err))
			}
		}()

		consoleLocalHost := fmt.Sprintf("http://127.0.0.1:%v", cfg.Console.LocalPort)
		if rootArgs.output == "human" {
			p.CheckMark().Text("Console proxy pointing to ").Bold(cfg.Console.URL).Text(" started at: ").SuccessText(consoleLocalHost).Jump()
		} else if rootArgs.output == "json" {
			log.Info(fmt.Sprintf("console proxy pointing to %s started at: %s", cfg.Console.URL, consoleLocalHost))
		}

		if !serviceRunArgs.noBrowser {
			err := open.Run(cs.GetBrowserURL())
			if err != nil {
				log.Error("unable to open the console in the default browser",
					zap.Error(err))
			}
		}
	}

	if rootArgs.output == "human" {
		p.CheckMark().SuccessText("Starting successful").NJump(2)
		p.BlueArrow().InfoText("Available endpoints").Jump()
		printEndpoints(serviceHost)
		p.NJump(2)
		p.BlueArrow().InfoText("Logs").Jump()
	}

	waitSig(ctx, cancel, log)

	err = srv.Stop()
	if err != nil {
		log.Error("error stopping wallet http server", zap.Error(err))
	} else {
		log.Info("wallet http server stopped with success")
	}

	if serviceRunArgs.startConsole {
		err = cs.Stop()
		if err != nil {
			log.Error("error stopping console proxy server", zap.Error(err))
		} else {
			log.Info("console proxy server stopped with success")
		}
	}

	return nil
}

// waitSig will wait for a sigterm or sigint interrupt.
func waitSig(ctx context.Context, cfunc func(), log *zap.Logger) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	select {
	case sig := <-gracefulStop:
		log.Info("caught signal", zap.String("wallet", fmt.Sprintf("%+v", sig)))
		cfunc()
	case <-ctx.Done():
		// nothing to do
	}
}
