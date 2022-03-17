package cmd

import (
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"fmt"
	"io"
	"os"
	"text/template"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"github.com/spf13/cobra"
)

const startupT = ` # Authentication
 - login:                   POST   {{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:                  DELETE {{.WalletServiceLocalAddress}}/api/v1/auth/token

 # Network management
 - network:                 GET    {{.WalletServiceLocalAddress}}/api/v1/network

 # Wallet management
 - create a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets
 - import a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets/import

 # Key pair management
 - generate a key pair:     POST   {{.WalletServiceLocalAddress}}/api/v1/keys
 - list keys:               GET    {{.WalletServiceLocalAddress}}/api/v1/keys
 - describe a key pair:     GET    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid
 - taint a key pair:        PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/taint
 - annotate a key pair:     PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/metadata

 # Commands
 - sign a command:          POST   {{.WalletServiceLocalAddress}}/api/v1/command
 - sign a command (sync):   POST   {{.WalletServiceLocalAddress}}/api/v1/command/sync
 - sign a command (commit): POST   {{.WalletServiceLocalAddress}}/api/v1/command/commit
 - sign data:               POST   {{.WalletServiceLocalAddress}}/api/v1/sign
 - verify data:             POST   {{.WalletServiceLocalAddress}}/api/v1/verify

 # Information
 - get service status:      GET    {{.WalletServiceLocalAddress}}/api/v1/status
 - get the version:         GET    {{.WalletServiceLocalAddress}}/api/v1/version
`

var (
	listEndpointsLong = cli.LongDesc(`
		List the Vega wallet service HTTP endpoints
	`)

	listEndpointsExample = cli.Examples(`
		# List service endpoints
		vegawallet endpoints --network NETWORK
	`)
)

type ListEndpointsHandler func(io.Writer, *RootFlags, *ListEndpointsFlags) error

func NewCmdListEndpoints(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdListEndpoints(w, ListEndpoints, rf)
}

func BuildCmdListEndpoints(w io.Writer, handler ListEndpointsHandler, rf *RootFlags) *cobra.Command {
	f := &ListEndpointsFlags{}

	cmd := &cobra.Command{
		Use:     "endpoints",
		Short:   "List endpoints",
		Long:    listEndpointsLong,
		Example: listEndpointsExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := f.Validate(); err != nil {
				return err
			}

			if err := handler(w, rf, f); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network configuration to use",
	)

	return cmd
}

type ListEndpointsFlags struct {
	Network string
}

func (f *ListEndpointsFlags) Validate() error {
	if len(f.Network) == 0 {
		return flags.FlagMustBeSpecifiedError("network")
	}

	return nil
}

func ListEndpoints(w io.Writer, rf *RootFlags, f *ListEndpointsFlags) error {
	p := printer.NewInteractivePrinter(w)

	vegaPaths := paths.New(rf.Home)
	netStore, err := netstore.InitialiseStore(vegaPaths)
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	cfg, err := netStore.GetNetwork(f.Network)
	if err != nil {
		return fmt.Errorf("couldn't initialise network store: %w", err)
	}

	serviceHost := fmt.Sprintf("http://%v:%v", cfg.Host, cfg.Port)

	p.BlueArrow().InfoText("Available endpoints").NextLine()
	printServiceEndpoints(serviceHost)
	p.NextLine()

	return nil
}

func printServiceEndpoints(serviceHost string) {
	params := struct {
		WalletServiceLocalAddress string
	}{
		WalletServiceLocalAddress: serviceHost,
	}

	tmpl, err := template.New("wallet-cmdline").Parse(startupT)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		panic(err)
	}
}
