package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/version"
	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

var (
	requestTimeout = 30 * time.Second

	rootExamples = cli.Examples(`
		# Specify a custom Vega home directory
		vegawallet --home PATH_TO_DIR COMMAND

		# Change the output to JSON
		vegawallet --output json COMMAND

		# Disable colors on output using environment variable
		NO_COLOR=1 vegawallet COMMAND

		# Disable the verification of the software version
		vegawallet --no-version-check COMMAND
	`)
)

type CheckVersionHandler func() (*semver.Version, error)

func NewCmdRoot(w io.Writer) *cobra.Command {
	vh := func() (*semver.Version, error) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		v, err := version.Check(version.BuildReleasesRequestFromGithub(ctx), version.Version)
		if err != nil {
			return nil, fmt.Errorf("couldn't check latest releases: %w", err)
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
		Example:       rootExamples,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// The `__complete` command is being run to build up the auto-completion
			// file. We should skip any verification to not temper with the process.
			// Any additional printing will end up in the auto-completion registry.
			// The `completion` command output the completion script for a given
			// shell, that should not be tempered with. We should skip it as well.
			if cmd.Name() == "__complete" || cmd.Name() == "completion" {
				return nil
			}

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
					p.CrossMark().DangerText(err.Error()).NextSection()
					return nil
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

	_ = cmd.MarkPersistentFlagDirname("home")
	_ = cmd.RegisterFlagCompletionFunc("output", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return flags.AvailableOutputs, cobra.ShellCompDirectiveDefault
	})

	// Root commands
	cmd.AddCommand(NewCmdInit(w, f))
	cmd.AddCommand(NewCmdCompletion(w))
	cmd.AddCommand(NewCmdVersion(w, f))

	// Sub-commands
	cmd.AddCommand(NewCmdCommand(w, f))
	cmd.AddCommand(NewCmdKey(w, f))
	cmd.AddCommand(NewCmdNetwork(w, f))
	cmd.AddCommand(NewCmdService(w, f))
	cmd.AddCommand(NewCmdTx(w, f))
	cmd.AddCommand(NewCmdMessage(w, f))

	// Wallet commands
	// We don't have a wrapper sub-command for wallet commands.
	cmd.AddCommand(NewCmdCreateWallet(w, f))
	cmd.AddCommand(NewCmdDeleteWallet(w, f))
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
	return flags.ValidateOutput(f.Output)
}
