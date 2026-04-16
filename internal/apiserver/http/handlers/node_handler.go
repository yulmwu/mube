package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mube/internal/apiserver/store"
	v1 "mube/pkg/api/v1"
)

type NodeHandler struct {
	store store.NodeStore
}

func NewNodeHandler(store store.NodeStore) *NodeHandler {
	return &NodeHandler{store: store}
}

func (h *NodeHandler) Heartbeat(c *gin.Context) {
	var req v1.NodeHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := h.store.UpsertHeartbeat(req, time.Now())
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, node)
}

func (h *NodeHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, v1.NodeListResponse{Items: h.store.List(time.Now())})
}
