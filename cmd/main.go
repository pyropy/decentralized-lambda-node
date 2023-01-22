package main

import (
	"github.com/pyropy/decentralised-lambda/api"
	"github.com/pyropy/decentralised-lambda/node"
)

func main() {
	apiCfg := api.Config{
		Host: "localhost",
		Port: "6969",
	}
	nodeCfg := node.Config{
		IPFSEndpoint:     "/ip4/0.0.0.0/tcp/5001",
		BacalhauEndpoint: "http://0.0.0.0:58859",
	}
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

	node, err := node.NewNode(&nodeCfg)
	if err != nil {
		panic(err)
	}

	api := api.NewServer(&apiCfg, node)
	api.Run()
}
