package cmd

import (
	"os"
	"text/template"
)

const startupT = ` # Authentication
 - login:                   POST   {{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:                  DELETE {{.WalletServiceLocalAddress}}/api/v1/auth/token

 # Network management
 - network:                 GET    {{.WalletServiceLocalAddress}}/api/v1/network

 # Wallet management
 - create a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets
 - import a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets/import
 - download a wallet:       GET    {{.WalletServiceLocalAddress}}/api/v1/wallets

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

func printEndpoints(serviceHost string) {
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
