package cmd

import (
	"os"
	"text/template"
)

const startupT = ` - get service status:      GET    {{.WalletServiceLocalAddress}}/api/v1/status
 - config:                  GET    {{.WalletServiceLocalAddress}}/api/v1/config
 - login:                   POST   {{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:                  DELETE {{.WalletServiceLocalAddress}}/api/v1/auth/token
 - create a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets
 - import a wallet:         POST   {{.WalletServiceLocalAddress}}/api/v1/wallets/import
 - create a key:            POST   {{.WalletServiceLocalAddress}}/api/v1/keys
 - list keys:               GET    {{.WalletServiceLocalAddress}}/api/v1/keys
 - get a key:               GET    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid
 - taint a key:             PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/taint
 - annotate a key:          PUT    {{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/metadata
 - sign data:               POST   {{.WalletServiceLocalAddress}}/api/v1/sign
 - verify data:             POST   {{.WalletServiceLocalAddress}}/api/v1/verify
 - sign a command:          POST   {{.WalletServiceLocalAddress}}/api/v1/command
 - sign a command (sync):   POST   {{.WalletServiceLocalAddress}}/api/v1/command/sync
 - sign a command (commit): POST   {{.WalletServiceLocalAddress}}/api/v1/command/commit
 - download a wallet:       GET    {{.WalletServiceLocalAddress}}/api/v1/wallets
 - get the version:         GET    {{.WalletServiceLocalAddress}}/api/v1/wallets
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
