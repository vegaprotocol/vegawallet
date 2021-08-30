package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	vgfs "code.vegaprotocol.io/go-wallet/libs/fs"
	vgjson "code.vegaprotocol.io/go-wallet/libs/json"
	"code.vegaprotocol.io/go-wallet/version"
	"code.vegaprotocol.io/go-wallet/wallet"
	wstorev1 "code.vegaprotocol.io/go-wallet/wallet/store/v1"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	rootArgs struct {
		output         string
		rootPath       string
		noVersionCheck bool
	}

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:               os.Args[0],
		Short:             "The Vega wallet",
		Long:              `The Vega wallet`,
		PersistentPreRunE: rootPreRun,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}
)

func rootPreRun(_ *cobra.Command, _ []string) error {
	err := parseOutputFlag()
	if err != nil {
		return err
	}
	if rootArgs.output == "human" {
		checkVersion()
	}
	return nil
}

func parseOutputFlag() error {
	if rootArgs.output == "human" && !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return errors.New("output \"human\" is not script-friendly, use \"json\" instead")
	}

	supportedOutput := []string{"json", "human"}
	for _, output := range supportedOutput {
		if rootArgs.output == output {
			return nil
		}
	}

	return fmt.Errorf("unsupported output \"%s\"", rootArgs.output)
}

func checkVersion() {
	if !rootArgs.noVersionCheck {
		v, err := version.Check(version.Version)
		if err != nil {
			fmt.Printf("could not check vega wallet version updates: %v\n", err)
		}
		if v != nil {
			p := termenv.ColorProfile()
			xVersion := termenv.String(v.String()).Foreground(p.Color("6"))
			xURL := termenv.String(version.GetReleaseURL(v)).Underline()
			fmt.Printf("A new version %s of vega wallet is available!\nDownload it at %v\n\n", xVersion, xURL)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if rootArgs.output == "human" && !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			fmt.Println(err)
		} else {
			if rootArgs.output == "human" {
				p := printer.NewHumanPrinter()
				p.CrossMark().DangerText(err.Error()).Jump()
			} else if rootArgs.output == "json" {
				jsonErr := vgjson.Print(struct {
					Error string
				}{
					Error: err.Error(),
				})
				if jsonErr != nil {
					fmt.Printf("couldn't format JSON: %v\n", jsonErr)
					fmt.Printf("original error: %v\n", err)
				}
			} else {
				fmt.Println(err)
			}
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootArgs.output, "output", "o", "human", "Specify the output format: json,human")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.rootPath, "root-path", "r", vgfs.DefaultVegaDir(), "Root directory for the Vega wallet configuration")
	rootCmd.PersistentFlags().BoolVar(&rootArgs.noVersionCheck, "no-version-check", false, "Do not check for new version of the Vega wallet")
}

func getPassphrase(flaggedPassphraseFile string, confirmInput bool) (string, error) {
	hasPassphraseFileFlag := len(flaggedPassphraseFile) != 0

	if hasPassphraseFileFlag {
		rawPassphrase, err := os.ReadFile(flaggedPassphraseFile)
		if err != nil {
			return "", err
		}
		// user might have added \n at the end of the line, let's remove it.
		cleanupPassphrase := strings.Trim(string(rawPassphrase), "\n")
		return cleanupPassphrase, nil
	} else {
		if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			return "", errors.New("passphrase-file flag required without TTY")
		}

		passphrase, err := promptForPassphrase()
		if err != nil {
			return "", fmt.Errorf("could not get passphrase: %w", err)
		}

		if len(passphrase) == 0 {
			return "", fmt.Errorf("passphrase cannot be empty")
		}

		if confirmInput {
			confirmation, err := promptForPassphrase("Confirm passphrase: ")
			if err != nil {
				return "", fmt.Errorf("could not get passphrase: %w", err)
			}

			if passphrase != confirmation {
				return "", fmt.Errorf("passphrases do not match")
			}
		}
		fmt.Println()

		return passphrase, nil
	}
}

func promptForPassphrase(msg ...string) (string, error) {
	if len(msg) <= 0 {
		fmt.Print("Enter passphrase: ")
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
