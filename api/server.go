package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pyropy/decentralised-lambda/node"
)

type Server struct {
	node   *node.Node
	config *Config
}

func NewServer(cfg *Config, node *node.Node) *Server {
	return &Server{config: cfg, node: node}
}

func (a *Server) Run() {
	router := gin.Default()
	router.POST("/invoke/:cid", a.HandleInvokeLambda)

	addr := fmt.Sprintf("%s:%s", a.config.Host, a.config.Port)
	router.Run(addr)
}
