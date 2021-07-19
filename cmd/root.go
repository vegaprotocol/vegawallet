package cmd

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/fsutil"
	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"code.vegaprotocol.io/go-wallet/version"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	rootArgs struct {
		rootPath       string
		noVersionCheck bool
	}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:              "vegawallet",
		Short:            "The Vega wallet",
		Long:             `The Vega wallet`,
		PersistentPreRun: checkVersion,
	}
)

func checkVersion(cmd *cobra.Command, args []string) {
	if !rootArgs.noVersionCheck {
		v, err := version.Check(Version)
		if err != nil {
			fmt.Printf("could not check vega wallet version updates: %v\n", err)
		}
		if v != nil {
			fmt.Printf("A new version %v of vega wallet is available, you can download it at %v.\n",
				v, version.GetReleaseURL(v))
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootArgs.rootPath, "root-path", fsutil.DefaultVegaDir(), "Root directory for the Vega wallet configuration")
	rootCmd.PersistentFlags().BoolVar(&rootArgs.noVersionCheck, "no-version-check", false, "Do not check for new version of the Vega wallet")
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

func getStore() (*storev1.Store, error) {
	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return nil, err
	}

	if err := store.Initialise(); err != nil {
		return nil, err
	}
	return store, nil
}
