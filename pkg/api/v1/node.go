package v1

import "time"

const (
	NodeStateUnknown  = "Unknown"
	NodeStateReady    = "Ready"
	NodeStateNotReady = "NotReady"
)

type NodeHeartbeatRequest struct {
	Name     string `json:"name" binding:"required"`
	Runtime  string `json:"runtime" binding:"required"`
	Version  string `json:"version" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,gte=1"`
}

type NodeStatus struct {
	Name          string    `json:"name"`
	IP            string    `json:"ip"`
	Port          int       `json:"port"`
	State         string    `json:"state"`
	Runtime       string    `json:"runtime"`
	Version       string    `json:"version"`
	Capacity      int       `json:"capacity"`
	LastHeartbeat time.Time `json:"lastHeartbeat"`
}

type NodeListResponse struct {
	Items []NodeStatus `json:"items"`
}
