package cmd

import (
	"io"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/version"

	"github.com/spf13/cobra"
)

var (
	versionLong = cli.LongDesc(`
		Get the version of the software.

		This is NOT related to the wallet version. To get information about the wallet,
		use the "info" command.
	`)

	versionExample = cli.Examples(`
		# Get the version of the software
		vegawallet version
	`)
)

type GetVersionHandler func() *version.GetVersionResponse

func NewCmdVersion(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdGetVersion(w, version.GetVersionInfo, rf)
}

func BuildCmdGetVersion(w io.Writer, handler GetVersionHandler, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Get the version of the software",
		Long:    versionLong,
		Example: versionExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			resp := handler()

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintGetVersionResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	return cmd
}

func PrintGetVersionResponse(w io.Writer, resp *version.GetVersionResponse) {
	p := printer.NewInteractivePrinter(w)

	p.Text("Version:").NextLine().WarningText(resp.Version).NextSection()
	p.Text("Git hash:").NextLine().WarningText(resp.GitHash).NextSection()
}
