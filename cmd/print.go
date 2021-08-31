package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallet"
)

func printKeyPair(p *printer.HumanPrinter, keyPair wallet.KeyPair) {
	p.Text("Private key:").Jump().WarningText(keyPair.PrivateKey()).Jump()
	p.Text("Public key:").Jump().WarningText(keyPair.PublicKey()).Jump()
	p.Text("Algorithm:").Jump().WarningText(fmt.Sprintf("%s (version %v)", keyPair.AlgorithmName(), keyPair.AlgorithmVersion())).Jump()
	p.Text("Tainted:").Jump().WarningText(fmt.Sprintf("%v", keyPair.IsTainted())).Jump()
	p.Text("Metadata:").Jump()
	printMeta(p, keyPair.Meta())
}

func printMeta(p *printer.HumanPrinter, meta []wallet.Meta) {
	padding := 0
	for _, m := range meta {
		keyLen := len(m.Key)
		if keyLen > padding {
			padding = keyLen
		}
	}

	for _, m := range meta {
		p.WarningText(fmt.Sprintf("%-*s", padding, m.Key)).Text(" | ").WarningText(m.Value).Jump()
	}
}

