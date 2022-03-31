package cmd

import (
	"fmt"
	"io"
	"os"

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

func NewCmdVersion(w io.Writer) *cobra.Command {
	return BuildCmdGetVersion(w, version.GetVersionInfo)
}

func BuildCmdGetVersion(w io.Writer, handler GetVersionHandler) *cobra.Command {
	f := &GetVersionFlags{}

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Get the version of the software",
		Long:    versionLong,
		Example: versionExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := f.Validate(); err != nil {
				return err
			}

			resp := handler()

			switch f.Output {
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

type GetVersionFlags struct {
	Output string
}

func (f *GetVersionFlags) Validate() error {
	return flags.ValidateOutput(f.Output)
}

func PrintGetVersionResponse(w io.Writer, resp *version.GetVersionResponse) {
	p := printer.NewInteractivePrinter(w)

	p.Text("Software version:").NextLine().WarningText(resp.Version).NextSection()
	p.Text("Git hash:").NextLine().WarningText(resp.GitHash).NextSection()

	p.RedArrow().DangerText("Important").NextLine()
	p.Text("This command does NOT give you your wallet version.").NextLine()
	p.Text("To get this information, see the following command:").NextSection()
	p.Code(fmt.Sprintf("%s info --help", os.Args[0])).NextLine()
}
