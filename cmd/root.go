package cmd

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/version"
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
	v, err := version.Check(Version)
	if err != nil {
		fmt.Printf("could not check vega wallet version updates: %v\n", err)
	}
	if v != nil {
		fmt.Printf("A new version %v of vega wallet is available, you can download it at %v.\n",
			v, version.GetReleaseURL(v))
	}

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
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()

	return string(password), nil
}
