package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"code.vegaprotocol.io/go-wallet/fsutil"
)

var (
	rootPath    string
	walletOwner string
	passphrase  string
	pubkey string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vegawallet",
	Short: "The vega wallet",
	Long: `The vega wallet`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootPath, "root-path", fsutil.DefaultVegaDir(), "config file (default is $HOME/.vega)")
}
