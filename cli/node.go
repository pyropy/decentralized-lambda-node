package cli

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pyropy/decentralised-lambda/api"
	"github.com/pyropy/decentralised-lambda/node"
	ufcli "github.com/urfave/cli/v2"
)

var NodeCmd = ufcli.Command{
	Name: "node",
	Subcommands: []*ufcli.Command{
		&NodeStartCmd,
	},
}

var NodeStartCmd = ufcli.Command{
	Name:  "start",
	Usage: "Starts invocation node server",
	Flags: []ufcli.Flag{},
	Action: func(cctx *ufcli.Context) error {
		var apiCfg api.Config
		var nodeCfg node.Config

		err := envconfig.Process("", &apiCfg)
		if err != nil {
			return err
		}

		err = envconfig.Process("", &nodeCfg)
		if err != nil {
			return err
		}

		node, err := node.NewNode(&nodeCfg)
		if err != nil {
			return err
		}

		api := api.NewServer(&apiCfg, node)
		api.Run()

		return nil
	},
}
