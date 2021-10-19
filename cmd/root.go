package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/version"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// requestTimeout is the maximum time the program will wait for a response
	// after issuing a request.
	requestTimeout = 30 * time.Second

	rootArgs struct {
		output         string
		home           string
		noVersionCheck bool
	}

	rootCmd = &cobra.Command{
		Use:               os.Args[0],
		Short:             "The Vega wallet",
		Long:              "The Vega wallet",
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
		return checkVersion()
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

func checkVersion() error {
	if !rootArgs.noVersionCheck {
		p := printer.NewHumanPrinter()
		if version.IsUnreleased() {
			p.CrossMark().DangerText("You are running an unreleased version of the Vega wallet. Use it at your own risk!").NJump(2)
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
			defer cancel()
			v, err := version.Check(version.BuildReleasesRequestFromGithub(ctx), version.Version)
			if err != nil {
				return fmt.Errorf("could not check Vega wallet version updates: %w", err)
			}
			if v != nil {
				p.Text("Version ").SuccessText(v.String()).Text(" is available. Your current version is ").DangerText(version.Version).Text(".").Jump()
				p.Text("Download the latest version at: ").Underline(version.GetReleaseURL(v)).NJump(2)
			}
		}
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if rootArgs.output == "human" && !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		} else {
			if rootArgs.output == "human" {
				p := printer.NewHumanPrinter()
				p.CrossMark().DangerText(err.Error()).Jump()
			} else if rootArgs.output == "json" {
				jsonErr := vgjson.Print(struct {
					Error string `json:"error"`
				}{
					Error: err.Error(),
				})
				if jsonErr != nil {
					_, _ = fmt.Fprintf(os.Stderr, "couldn't format JSON: %v\n", jsonErr)
					_, _ = fmt.Fprintf(os.Stderr, "original error: %v\n", err)
				}
			} else {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootArgs.output, "output", "o", "human", "Specify the output format: json,human")
	rootCmd.PersistentFlags().StringVar(&rootArgs.home, "home", "", "Specify the location of a custom Vega home")
	rootCmd.PersistentFlags().BoolVar(&rootArgs.noVersionCheck, "no-version-check", false, "Do not check for new version of the Vega wallet")
}

func getPassphrase(flaggedPassphraseFile string, confirmInput bool) (string, error) {
	hasPassphraseFileFlag := len(flaggedPassphraseFile) != 0

	if hasPassphraseFileFlag {
		passphraseDir, passphraseFileName := filepath.Split(flaggedPassphraseFile)
		rawPassphrase, err := fs.ReadFile(os.DirFS(passphraseDir), passphraseFileName)
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
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()

	return string(password), nil
}
