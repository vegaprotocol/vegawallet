package main

import (
	"os"

	"code.vegaprotocol.io/vegawallet/cmd"
)

func main() {
	writer := &cmd.Writer{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Execute(writer)
}
