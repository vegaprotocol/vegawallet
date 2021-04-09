package cmd

import (
	"os"
	"text/template"
)

const startupT = `Wallet http service started at: http://{{.WalletServiceLocalAddress}}
Wallet console proxy started at: http://{{.ConsoleProxyLocalAddress}} proxying {{.ConsoleProxyProxiedAddress}}

Available endpoints:
 - status:          GET    http://{{.WalletServiceLocalAddress}}/api/v1/status
 - login:           POST   http://{{.WalletServiceLocalAddress}}/api/v1/auth/token
 - logout:          DELETE http://{{.WalletServiceLocalAddress}}/api/v1/auth/token
 - create wallet:   POST   http://{{.WalletServiceLocalAddress}}/api/v1/wallets
 - create key:      POST   http://{{.WalletServiceLocalAddress}}/api/v1/keys
 - list keys:       GET    http://{{.WalletServiceLocalAddress}}/api/v1/keys
 - get key:         GET    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid
 - taint key:       PUT    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/taint
 - update meta:     PUT    http://{{.WalletServiceLocalAddress}}/api/v1/keys/:keyid/metadata
 - sign data:       POST   http://{{.WalletServiceLocalAddress}}/api/v1/sign
 - sign:            POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages
 - sign sync:       POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages/sync
 - sign commit:     POST   http://{{.WalletServiceLocalAddress}}/api/v1/messages/commit
 - download wallet: GET    http://{{.WalletServiceLocalAddress}}/api/v1/wallets
`

func printStartupMessage(consoleURL, consoleLocalHost, serviceHost string) {
	params := struct {
		ConsoleProxyLocalAddress,
		ConsoleProxyProxiedAddress,
		WalletServiceLocalAddress string
	}{
		ConsoleProxyLocalAddress:   consoleLocalHost,
		ConsoleProxyProxiedAddress: consoleURL,
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
