package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
)

func printKeyPair(p *printer.HumanPrinter, keyPair wallet.KeyPair) {
	p.Text("Private key:").NextLine().WarningText(keyPair.PrivateKey()).NextLine()
	p.Text("Public key:").NextLine().WarningText(keyPair.PublicKey()).NextLine()
	p.Text("Algorithm:").NextLine().WarningText(fmt.Sprintf("%s (version %v)", keyPair.AlgorithmName(), keyPair.AlgorithmVersion())).NextLine()
	p.Text("Tainted:").NextLine().WarningText(fmt.Sprintf("%v", keyPair.IsTainted())).NextLine()
	p.Text("Metadata:").NextLine()
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
		p.WarningText(fmt.Sprintf("%-*s", padding, m.Key)).Text(" | ").WarningText(m.Value).NextLine()
	}
}
