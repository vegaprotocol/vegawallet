package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the vega wallet service",
	Long:  "Start a vega wallet service behind an http server",
	RunE:  runServiceRun,
}

func init() {
	serviceCmd.AddCommand(runCmd)
}

func runServiceRun(cmd *cobra.Command, args []string) error {
	cfg, err := wallet.LoadConfig(rootPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	srv, err := wallet.NewService(log, cfg, rootPath)
	if err != nil {
		return err
	}
	go func() {
		defer cancel()
		err := srv.Start()
		if err != nil {
			log.Error("error starting wallet http server", zap.Error(err))
		}
	}()

	waitSig(ctx, log)

	err = srv.Stop()
	if err != nil {
		log.Error("error stopping wallet http server", zap.Error(err))
	} else {
		log.Info("wallet http server stopped with success")
	}

	return nil
}

// waitSig will wait for a sigterm or sigint interrupt.
func waitSig(ctx context.Context, log *zap.Logger) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	select {
	case sig := <-gracefulStop:
		log.Info("Caught signal", zap.String("name", fmt.Sprintf("%+v", sig)))
	case <-ctx.Done():
		// nothing to do
	}
}
