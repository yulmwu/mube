package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mube/internal/apiserver/http/handlers"
)

func NewRouter(nodeHandler *handlers.NodeHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/nodes", nodeHandler.List)
		v1.POST("/nodes/heartbeat", nodeHandler.Heartbeat)
	}

	return r
}
