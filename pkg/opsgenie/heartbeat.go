// Package opsgenie provides a client for interacting with the OpsGenie API.
package opsgenie

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/heartbeat"
	"github.com/sirupsen/logrus"
)

// HeartbeatClient is a wrapper around the OpsGenie heartbeat client.
type HeartbeatClient struct {
	*heartbeat.Client
}

// NewHeartbeatClient creates a new HeartbeatClient instance.
func NewHeartbeatClient(apiUrl, envVar string) (*HeartbeatClient, error) {
	logger := logrus.New()
	logger.Out = io.Discard

	config := &client.Config{
		OpsGenieAPIURL: client.ApiUrl(apiUrl),
		ApiKey:         os.Getenv(envVar),
		Logger:         logger,
	}

	heartbeatClient, err := heartbeat.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsGenie heartbeat client: %w", err)
	}

	h := &HeartbeatClient{
		Client: heartbeatClient,
	}

	return h, nil
}

// ListHeartbeats retrieves all heartbeats from OpsGenie.
func (c *HeartbeatClient) ListHeartbeats(ctx context.Context) ([]heartbeat.Heartbeat, error) {
	result, err := c.Client.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list heartbeats: %w", err)
	}

	return result.Heartbeats, nil
}

// GetHeartbeat retrieves a single heartbeat by its name.
func (c *HeartbeatClient) GetHeartbeat(ctx context.Context, name string) (*heartbeat.Heartbeat, error) {
	if name == "" {
		return nil, fmt.Errorf("heartbeat name cannot be empty")
	}

	result, err := c.Client.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get heartbeat %s: %w", name, err)
	}

	return &result.Heartbeat, nil
}
