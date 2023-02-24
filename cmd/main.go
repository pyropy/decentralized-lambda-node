package main

import (
	"github.com/pyropy/decentralised-lambda/api"
	"github.com/pyropy/decentralised-lambda/node"
)

func main() {
	apiCfg := api.DefaultConfig()
	nodeCfg := node.DefaultConfig()
	//var apiCfg api.Config
	//var nodeCfg node.Config
	//
	//err := envconfig.Process("", &apiCfg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = envconfig.Process("", &nodeCfg)
	//if err != nil {
	//	panic(err)
	//}
	//

	node, err := node.NewNode(nodeCfg)
	if err != nil {
		panic(err)
	}

	api := api.NewServer(apiCfg, node)
	api.Run()
}
