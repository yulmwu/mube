package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mube/internal/mubelet/config"
)

func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"node":      cfg.NodeName,
			"runtime":   cfg.Runtime,
			"version":   cfg.Version,
			"capacity":  cfg.Capacity,
			"apiServer": cfg.APIServerEndpoint,
		})
	})

	return r
}
