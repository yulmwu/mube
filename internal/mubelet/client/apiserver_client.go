package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	v1 "mube/pkg/api/v1"
)

type APIServerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIServerClient(baseURL string) *APIServerClient {
	return &APIServerClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *APIServerClient) SendHeartbeat(ctx context.Context, req v1.NodeHeartbeatRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal heartbeat: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/nodes/heartbeat", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build heartbeat request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call api server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status code from api server: %d", resp.StatusCode)
	}

	return nil
}
