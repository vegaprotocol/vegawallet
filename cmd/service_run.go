package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"code.vegaprotocol.io/go-wallet/logger"
	"code.vegaprotocol.io/go-wallet/service"
	svcstore1 "code.vegaprotocol.io/go-wallet/service/store/v1"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serviceRunArgs struct {
		consoleProxy bool
		noBrowser    bool
	}

	serviceRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Start the vega wallet service",
		Long:  "Start a vega wallet service behind an HTTP server",
		RunE:  runServiceRun,
	}
)

func init() {
	serviceCmd.AddCommand(serviceRunCmd)
	serviceRunCmd.Flags().BoolVarP(&serviceRunArgs.consoleProxy, "console-proxy", "p", false, "Start the vega console proxy and open the console in the default browser")
	serviceRunCmd.Flags().BoolVarP(&serviceRunArgs.noBrowser, "no-browser", "n", false, "Do not open the default browser if the console proxy is stated")
}

func runServiceRun(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	svcStore, err := svcstore1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	cfg, err := svcStore.GetConfig()
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.Level.Level)
	if err != nil {
		return err
	}
	defer log.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := service.NewService(log, cfg, svcStore, handler, Version, VersionHash)
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

	var cproxy *consoleProxy
	if serviceRunArgs.consoleProxy {
		cproxy = newConsoleProxy(log, cfg.Console.LocalPort, cfg.Console.URL, cfg.Nodes.Hosts[0], Version)
		go func() {
			defer cancel()
			err := cproxy.Start()
			if err != nil && err != http.ErrServerClosed {
				log.Error("error starting console proxy server", zap.Error(err))
			}
		}()

		if !serviceRunArgs.noBrowser {
			err := open.Run(cproxy.GetBrowserURL())
			if err != nil {
				log.Error("unable to open the console in the default browser",
					zap.Error(err))
			}
		}
	}

	printStartupMessage(
		cfg.Console.URL,
		fmt.Sprintf("127.0.0.1:%v", cfg.Console.LocalPort),
		fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
	)

	waitSig(ctx, cancel, log)

	err = srv.Stop()
	if err != nil {
		log.Error("error stopping wallet http server", zap.Error(err))
	} else {
		log.Info("wallet http server stopped with success")
	}

	if serviceRunArgs.consoleProxy {
		err = cproxy.Stop()
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
