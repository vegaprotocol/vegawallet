package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vgfs "code.vegaprotocol.io/go-wallet/libs/fs"
	"code.vegaprotocol.io/go-wallet/version"
	"code.vegaprotocol.io/go-wallet/wallet"
	wstorev1 "code.vegaprotocol.io/go-wallet/wallet/store/v1"
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
		v, err := version.Check(version.Version)
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
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootArgs.rootPath, "root-path", "r", vgfs.DefaultVegaDir(), "Root directory for the Vega wallet configuration")
	rootCmd.PersistentFlags().BoolVar(&rootArgs.noVersionCheck, "no-version-check", false, "Do not check for new version of the Vega wallet")
}

func getPassphrase(flaggedPassphrase, flaggedPassphraseFile string, confirmInput bool) (string, error) {
	hasPassphraseFileFlag := len(flaggedPassphraseFile) != 0
	hasPassphraseFlag := len(flaggedPassphrase) != 0

	if hasPassphraseFlag && hasPassphraseFileFlag {
		return "", errors.New("can't have both passphrase and passphrase-file flags defined")
	}

	if hasPassphraseFlag {
		return flaggedPassphrase, nil
	} else if hasPassphraseFileFlag {
		rawPassphrase, err := os.ReadFile(flaggedPassphraseFile)
		if err != nil {
			return "", err
		}
		// user might have added \n at the end of the line, let's remove it.
		cleanupPassphrase := strings.Trim(string(rawPassphrase), "\n")
		return cleanupPassphrase, nil
	} else {
		passphrase, err := promptForPassphrase()
		if err != nil {
			return "", fmt.Errorf("could not get passphrase: %w", err)
		}

		if len(passphrase) == 0 {
			return "", fmt.Errorf("passphrase cannot be empty")
		}

		if confirmInput {
			confirmation, err := promptForPassphrase("please confirm passphrase:")
			if err != nil {
				return "", fmt.Errorf("could not get passphrase: %w", err)
			}

			if passphrase != confirmation {
				return "", fmt.Errorf("passphrases do not match")
			}
		}

		return passphrase, nil
	}
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

// newWalletsStore builds a wallets store with the following structure
//
// root-path/
// └── wallets/
//    ├── my-wallet-1
//    └── my-wallet-2
func newWalletsStore(rootPath string) (*wstorev1.Store, error) {
	walletsPath := filepath.Join(rootPath, "wallets")

	return wstorev1.NewStore(walletsPath)
}

func newWalletHandler(rootPath string) (*wallet.Handler, error) {
	store, err := newWalletsStore(rootPath)
	if err != nil {
		return nil, err
	}

	return wallet.NewHandler(store), nil
}
