package cmd

import (
	"os"
	"text/template"
)

const startupT = ` - status:          GET    http://{{.WalletServiceLocalAddress}}/api/v1/status
 - login:           POST   http://{{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:          DELETE http://{{.WalletServiceLocalAddress}}/api/v1/auth/token
 - create wallet:   POST   http://{{.WalletServiceLocalAddress}}/api/v1/wallets
 - import wallet:   POST   http://{{.WalletServiceLocalAddress}}/api/v1/wallets/import
 - create key:      POST   http://{{.WalletServiceLocalAddress}}/api/v1/keys
 - list keys:       GET    http://{{.WalletServiceLocalAddress}}/api/v1/keys
 - get key:         GET    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid
 - taint key:       PUT    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/taint
 - update meta:     PUT    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/metadata
 - sign data:       POST   http://{{.WalletServiceLocalAddress}}/api/v1/sign
 - sign v2:         POST   http://{{.WalletServiceLocalAddress}}/api/v1/command
 - sign sync v2:    POST   http://{{.WalletServiceLocalAddress}}/api/v1/command/sync
 - sign commit v2:  POST   http://{{.WalletServiceLocalAddress}}/api/v1/command/commit
 - sign:            POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages
 - sign sync:       POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages/sync
 - sign commit:     POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages/commit
 - download wallet: GET    http://{{.WalletServiceLocalAddress}}/api/v1/wallets
`

func printEndpoints(serviceHost string) {
	params := struct {
		WalletServiceLocalAddress string
	}{
		WalletServiceLocalAddress:  serviceHost,
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
