package main

import (
	"github.com/pyropy/decentralised-lambda/cli"
	ufcli "github.com/urfave/cli/v2"
)

func main() {
	app := &ufcli.App{
		Name: "decentralised-lambda",
		Commands: []*ufcli.Command{
			&cli.WasmCmd,
			&cli.NodeCmd,
		},
	}

	app.Setup()
	cli.RunApp(app)
}
