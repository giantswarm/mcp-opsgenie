// Package mcp provides handlers for integrating OpsGenie with the Model Context Protocol (MCP).
// It enables AI assistants to query and interact with OpsGenie alerts through standardized MCP tools.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/giantswarm/mcp-opsgenie/pkg/opsgenie"
)

// listAlertQueryDescription contains comprehensive documentation for OpsGenie alert search queries.
// This documentation is sourced from the official OpsGenie documentation:
// https://support.atlassian.com/opsgenie/docs/search-queries-for-alerts/
const listAlertQueryDescription = `Search query for filtering alerts.

## Field reference for alert search

You can search using field:value combinations with most alert fields:

Field | Description
-------------------
createdAt | Unix timestamp (ms) or DD-MM-YYYY. e.g. createdAt:1470394841148, createdAt:15-05-2020
lastOccurredAt, snoozedUntil | Unix timestamps in ms
alertId | id of the alert, e.g. b9a2fb13-1b76-4b41-be28-eed2c61978fa
tinyId | Short internal ID (not recommended)
alias, count, message, description, source, entity, status, owner, acknowledgedBy, closedBy, recipients | Use exact strings. Status can be open or closed; boolean fields: isSeen, acknowledged, snoozed
teams, integration.name, integration.type, tag, actions | Filter by team names, integration details, tags, etc
details.key, details.value | Nested details; e.g. details.key:Impact

## Condition operators

Use relational operators for numeric/timestamp fields.

Examples:

- count > 5
- count <= 4
- lastOccurredAt < 1470394841148

## Logical operators

Combine conditions using AND, OR, and parentheses.

Examples:

- message:(lorem OR ipsum)
- description:(lorem AND ipsum)
- message:lorem AND count >= 3
- (message:(lorem OR ipsum)) AND count >= 3
- status:open AND (count >= 3 OR entity:lipsum)

Negate with NOT.

Examples:

- NOT message:lorem
- NOT status:open

## Wildcards (*)

Use * only at the end of a word:

message: lorem* matches words beginning with "lorem" ("Lorem ipsum", "Lorem123"), but not within words ("dolorlorem")

Wildcards are not supported for teams and users — use full names

## Null queries

Check for presence or absence of fields:

Supported fields: source, entity, tag, actions, owner, teams, acknowledgedBy, closedBy, recipients, details.key, details.value, integration.name, integration.type.

Examples:

- owner:null — alerts without an owner
- teams is null — no teams assigned
- details.key is not null — alerts with details.key set
- tag !: null — alerts with a tag

## Tip

Combine any of the above to fine-tune your alert search queries.`

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

	// Define the list_alerts tooListAlertsmprehensive documentation
	tool := mcp.NewTool("list_alerts",
		mcp.WithDescription("Retrieve a list of alerts from OpsGenie"),
		mcp.WithString("query",
			mcp.Description(listAlertQueryDescription),
		),
	)

	// Register the tool with the MCP server
	s.AddTool(tool, handler.ListAlerts)

	return nil
}

// ListAlerts retrieves alerts from OpsGenie based on the provided search query.
// This method implements the MCP tool handler interface for the 'list_alerts' tool.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation and timeouts
//   - request: The MCP tool call request containing the search query parameter
//
// Returns:
//   - A CallToolResult containing the serialized alerts data on success
//   - A CallToolResult with error information on failure
//   - An error is only returned for internal MCP framework issues (always nil in this implementation)
func (h *opsgenieHandler) ListAlerts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract the query parameter (defaults to empty string if not provided)
	query := request.GetString("query", "")

	// Fetch alerts from OpsGenie using the provided query
	alerts, err := h.alertClient.ListAlerts(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve alerts from OpsGenie: %v", err)), nil
	}

	// Serialize the alerts to JSON for the MCP response
	data, err := json.Marshal(alerts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize alerts to JSON: %v", err)), nil
	}

	// Return the serialized alerts as a text result
	return mcp.NewToolResultText(string(data)), nil
}
