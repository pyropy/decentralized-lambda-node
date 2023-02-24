package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"
	"net/http"
)

func (a *Server) HandleInvokeLambda(c *gin.Context) {
	cidParam := c.Param("cid")

	wasmCid, err := cid.Parse(cidParam)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	jobID, err := a.node.InvokeLambdaFunction(ctx, wasmCid, c.Request)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobID": jobID})
}
