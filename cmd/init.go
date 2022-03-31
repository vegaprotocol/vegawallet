package cmd

import (
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/service"
	svcstore "code.vegaprotocol.io/vegawallet/service/store/v1"
	"code.vegaprotocol.io/vegawallet/wallets"

	"github.com/spf13/cobra"
)

var (
	initLong = cli.LongDesc(`
		Creates the folders, the configuration files and RSA keys needed by the service
		to operate.
	`)

	initExample = cli.Examples(`
		# Initialise the software
		vegawallet init

		# Re-initialise the software
		vegawallet init --force
	`)
)

type InitHandler func(home string, f *InitFlags) (*InitResponse, error)

func NewCmdInit(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdInit(w, Init, rf)
}

func BuildCmdInit(w io.Writer, handler InitHandler, rf *RootFlags) *cobra.Command {
	f := &InitFlags{}

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialise the software",
		Long:    initLong,
		Example: initExample,
		RunE: func(_ *cobra.Command, _ []string) error {

			if err := f.Validate(); err != nil {
				return err
			}

			resp, err := handler(rf.Home, f)
			if err != nil {
				return err
			}

			switch f.Output {
			case flags.InteractiveOutput:
				PrintInitResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&f.Force,
		"force", "f",
		false,
		"Overwrite exiting wallet configuration at the specified path",
	)

	addOutputFlag(cmd, &f.Output)

	return cmd
}

type InitFlags struct {
	Force  bool
	Output string
}

func (f *InitFlags) Validate() error {
	return flags.ValidateOutput(f.Output)
}

type InitResponse struct {
	RSAKeys struct {
		PublicKeyFilePath  string `json:"publicKeyFilePath"`
		PrivateKeyFilePath string `json:"privateKeyFilePath"`
	} `json:"rsaKeys"`
}

func Init(home string, f *InitFlags) (*InitResponse, error) {
	_, err := wallets.InitialiseStore(home)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	svcStore, err := svcstore.InitialiseStore(paths.New(home))
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise service store: %w", err)
	}

	if err = service.InitialiseService(svcStore, f.Force); err != nil {
		return nil, fmt.Errorf("couldn't initialise the service: %w", err)
	}

	resp := &InitResponse{}
	pubRSAKeysPath, privRSAKeysPath := svcStore.GetRSAKeysPath()
	resp.RSAKeys.PublicKeyFilePath = pubRSAKeysPath
	resp.RSAKeys.PrivateKeyFilePath = privRSAKeysPath

	return resp, nil
}

func PrintInitResponse(w io.Writer, resp *InitResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().Text("Service public RSA keys created at: ").SuccessText(resp.RSAKeys.PublicKeyFilePath).NextLine()
	p.CheckMark().Text("Service private RSA keys created at: ").SuccessText(resp.RSAKeys.PrivateKeyFilePath).NextLine()
	p.CheckMark().SuccessText("Initialisation succeeded").NextSection()

	p.BlueArrow().InfoText("Create a wallet").NextLine()
	p.Text("To create a wallet, use the following command:").NextSection()
	p.Code(fmt.Sprintf("%s create --wallet \"YOUR_USERNAME\"", os.Args[0])).NextSection()
	p.Text("The ").Bold("--wallet").Text(" flag sets the name of your wallet and will be used to login to Vega Console.").NextSection()
	p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
}
