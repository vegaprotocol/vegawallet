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
	"code.vegaprotocol.io/go-wallet/logger"
	"code.vegaprotocol.io/go-wallet/service"
	svcstore1 "code.vegaprotocol.io/go-wallet/service/store/v1"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serviceRunArgs struct {
		startConsole bool
		noBrowser    bool
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
	serviceRunCmd.Flags().BoolVarP(&serviceRunArgs.startConsole, "console-proxy", "p", false, "Start the vega console proxy and open the console in the default browser")
	serviceRunCmd.Flags().BoolVarP(&serviceRunArgs.noBrowser, "no-browser", "n", false, "Do not open the default browser if the console proxy is stated")
}

func runServiceRun(_ *cobra.Command, _ []string) error {
	p := printer.NewHumanPrinter()

	handler, err := newWalletHandler(rootArgs.vegaHome)
	if err != nil {
		return err
	}

	svcStore, err := svcstore1.NewStore(rootArgs.vegaHome)
	if err != nil {
		return err
	}

	cfg, err := svcStore.GetConfig()
	if err != nil {
		return err
	}

	encoding := "json"
	if rootArgs.output == "human" {
		encoding = "console"
	}

	log, err := logger.New(cfg.Level.Level, encoding)
	if err != nil {
		return err
	}
	defer log.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := service.NewService(log, cfg, svcStore, handler)
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
			log.Info(fmt.Sprintf("Console proxy pointing to %s started at: %s", cfg.Console.URL, consoleLocalHost))
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
		log.Info("Caught signal", zap.String("name", fmt.Sprintf("%+v", sig)))
		cfunc()
	case <-ctx.Done():
		// nothing to do
	}
}
