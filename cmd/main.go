package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/pyropy/decentralised-lambda/api"
	"github.com/pyropy/decentralised-lambda/node"
)

func main() {
	var apiCfg api.Config
	var nodeCfg node.Config

	err := envconfig.Process("", &apiCfg)
	if err != nil {
		panic(err)
	}

	err = envconfig.Process("", &nodeCfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(nodeCfg)

	node, err := node.NewNode(&nodeCfg)
	if err != nil {
		panic(err)
	}

	api := api.NewServer(&apiCfg, node)
	api.Run()
}
