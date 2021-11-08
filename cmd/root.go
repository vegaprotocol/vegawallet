package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/version"
	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

var requestTimeout = 30 * time.Second

type CheckVersionHandler func() (*semver.Version, error)

func NewCmdRoot(w io.Writer) *cobra.Command {
	vh := func() (*semver.Version, error) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		v, err := version.Check(version.BuildReleasesRequestFromGithub(ctx), version.Version)
		if err != nil {
			return nil, fmt.Errorf("couldn't check latest Vega wallet releases: %w", err)
		}
		return v, nil
	}

	return BuildCmdRoot(w, vh)
}

func BuildCmdRoot(w io.Writer, vh CheckVersionHandler) *cobra.Command {
	f := &RootFlags{}

	cmd := &cobra.Command{
		Use:           os.Args[0],
		Short:         "The Vega wallet",
		Long:          "The Vega wallet",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := f.Validate(); err != nil {
				return err
			}

			if !f.NoVersionCheck && f.Output == flags.InteractiveOutput {
				p := printer.NewInteractivePrinter(w)
				if version.IsUnreleased() {
					p.CrossMark().DangerText("You are running an unreleased version of the Vega wallet (").DangerText(version.Version).DangerText("). Use it at your own risk!").NextSection()
				}

				v, err := vh()
				if err != nil {
					return err
				}

				if v != nil {
					p.Text("Version ").SuccessText(v.String()).Text(" is available. Your current version is ").DangerText(version.Version).Text(".").NextLine()
					p.Text("Download the latest version at: ").Underline(version.GetReleaseURL(v)).NextSection()
				}
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&f.Output,
		"output", "o",
		flags.InteractiveOutput,
		fmt.Sprintf("Specify the output format: %v", flags.AvailableOutputs),
	)
	cmd.PersistentFlags().StringVar(&f.Home,
		"home",
		"",
		"Specify the location of a custom Vega home",
	)
	cmd.PersistentFlags().BoolVar(&f.NoVersionCheck,
		"no-version-check",
		false,
		"Do not check for new version of the Vega wallet",
	)

	// Root commands
	cmd.AddCommand(NewCmdSendCommand(w, f))
	cmd.AddCommand(NewCmdInit(w, f))
	cmd.AddCommand(NewCmdSignMessage(w, f))
	cmd.AddCommand(NewCmdVerifyMessage(w, f))
	cmd.AddCommand(NewCmdVersion(w, f))

	// Sub-commands
	cmd.AddCommand(NewCmdKey(w, f))
	cmd.AddCommand(NewCmdNetwork(w, f))
	cmd.AddCommand(NewCmdService(w, f))

	// Wallet commands
	// We don't have a wrapper sub-command for wallet commands.
	cmd.AddCommand(NewCmdGetInfoWallet(w, f))
	cmd.AddCommand(NewCmdImportWallet(w, f))
	cmd.AddCommand(NewCmdListWallets(w, f))

	return cmd
}

type RootFlags struct {
	Output         string
	Home           string
	NoVersionCheck bool
}

func (f *RootFlags) Validate() error {
	if err := flags.ValidateOutput(f.Output); err != nil {
		// This flag has special treatment because error reporting depends on it,
		// and we need to differentiate output errors from the rest to select the
		// right way to print the data.
		// As a result, we wrap generic errors in a specific one
		return NewInvalidOutputError(err)
	}

	return nil
}
