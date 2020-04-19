package main

import (
	"os"

	"code.vegaprotocol.io/go-wallet/fsutil"

	"github.com/spf13/cobra"
)

var (
	rootPath    string
	walletOwner string
	passphrase  string
	pubkey      string
	data        string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vegawallet",
	Short: "The vega wallet",
	Long:  `The vega wallet`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootPath, "root-path", fsutil.DefaultVegaDir(), "Root directory of the vegawalle configuration")
}
