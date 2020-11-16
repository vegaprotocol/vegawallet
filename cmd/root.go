package cmd

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

var (
	rootArgs struct {
		rootPath string
	}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "vegawallet",
		Short: "The Vega wallet",
		Long:  `The Vega wallet`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootArgs.rootPath, "root-path", fsutil.DefaultVegaDir(), "Root directory for the Vega wallet configuration")
}

func promptForPassphrase(msg ...string) (string, error) {
	if len(msg) <= 0 {
		fmt.Print("please enter passphrase:")
	} else {
		fmt.Print(msg[0])
	}
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Println()

	return string(password), nil
}
