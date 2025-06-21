// Package mcp provides handlers for integrating OpsGenie with the Model Context Protocol (MCP).
// It enables AI assistants to query and interact with OpsGenie alerts through standardized MCP tools.
package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/giantswarm/mcp-opsgenie/pkg/opsgenie"
)

// opsgenieHandler handles MCP tool requests for OpsGenie operations.
// It encapsulates the OpsGenie alert client and provides methods to interact with alerts.
type opsgenieHandler struct {
	alertClient *opsgenie.AlertClient
}

// RegisterOpsGenieHandler registers the OpsGenie MCP tools with the provided MCP server.
// It creates an alert client using the specified API URL and environment variable for authentication,
// then registers the available tools (currently 'list_alerts') with the server.
//
// Parameters:
//   - s: The MCP server instance to register tools with
//   - apiUrl: The OpsGenie API URL endpoint
//   - envVar: The name of the environment variable containing the OpsGenie API key
//
// Returns an error if the alert client cannot be created or if tool registration fails.
func RegisterOpsGenieHandler(s *server.MCPServer, apiUrl, envVar string) error {
	alertClient, err := opsgenie.NewAlertClient(apiUrl, envVar)
	if err != nil {
		return fmt.Errorf("failed to create OpsGenie alert client: %w", err)
	}

	// Initialize the handler with the alert client
	handler := &opsgenieHandler{
		alertClient: alertClient,
	}

	handler.registerAlertTools(s)

	return nil
}
