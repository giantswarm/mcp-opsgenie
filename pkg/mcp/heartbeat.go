package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (h *opsgenieHandler) registerHeartbeatTools(s *server.MCPServer) {
	listHeartbeatsTool := mcp.NewTool("list_heartbeats",
		mcp.WithDescription("Retrieve a list of all heartbeats from OpsGenie."),
	)
	s.AddTool(listHeartbeatsTool, h.ListHeartbeats)

	getHeartbeatTool := mcp.NewTool("get_heartbeat",
		mcp.WithDescription("Retrieves a single heartbeat from OpsGenie by its name."),
		mcp.WithString("name",
			mcp.Description("Name of the heartbeat to retrieve."),
			mcp.Required(),
		),
	)
	s.AddTool(getHeartbeatTool, h.GetHeartbeat)
}

// ListHeartbeats retrieves all heartbeats from OpsGenie.
func (h *opsgenieHandler) ListHeartbeats(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	heartbeats, err := h.heartbeatClient.ListHeartbeats(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve heartbeats from OpsGenie: %v", err)), nil
	}

	data, err := json.Marshal(heartbeats)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize heartbeats to JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// GetHeartbeat retrieves a single OpsGenie heartbeat by its name.
func (h *opsgenieHandler) GetHeartbeat(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	if name == "" {
		return mcp.NewToolResultError("the 'name' parameter is required"), nil
	}

	heartbeat, err := h.heartbeatClient.GetHeartbeat(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve heartbeat with name '%s' from OpsGenie: %v", name, err)), nil
	}

	data, err := json.Marshal(heartbeat)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize heartbeat to JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
