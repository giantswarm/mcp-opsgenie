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

| Field | Example Value | Description |
|-------|---------------|-------------|
| createdAt | 1470394841148 | Unix timestamp in milliseconds (Fri, 05 Aug 2016 11:00:41.148 GMT) |
| createdAt | 15-05-2020 | DD-MM-YYYY format |
| lastOccurredAt | 1470394841148 | Unix timestamp in milliseconds |
| snoozedUntil | 1470394841148 | Unix timestamp in milliseconds |
| alertId | b9a2fb13-1b76-4b41-be28-eed2c61978fa | Full alert ID |
| tinyId | 28 | Short ID (not recommended, it rolls) |
| alias | host_down | Alert alias |
| count | 5 | Number of alert occurrences |
| message | "Server apollo average" | Alert message text |
| description | "Monitoring tool is reporting..." | Alert description |
| source | john.smith@opsgenie.com | Alert source |
| entity | entity1 | Related entity |
| status | open | Alert status (open or closed) |
| owner | john.smith@opsgenie.com | Alert owner username |
| acknowledgedBy | john.smith@opsgenie.com | Who acknowledged the alert |
| closedBy | john.smith@opsgenie.com | Who closed the alert |
| recipients | john.smith@opsgenie.com | Alert recipients |
| isSeen | true | Whether alert has been seen (true/false) |
| acknowledged | true | Whether alert is acknowledged (true/false) |
| snoozed | false | Whether alert is snoozed (true/false) |
| teams | team1 | Team name |
| integration.name | "API Integration" | Integration name |
| integration.type | API | Integration type |
| tag | EC2 | Alert tag |
| actions | start | Available actions |
| details.key | Impact | Custom detail key |
| details.value | External | Custom detail value |

## Query Operators

**Comparison Operators (for numeric/timestamp fields):**
- Greater than: count > 5
- Less than: count < 10
- Greater than or equal: count >= 3
- Less than or equal: count <= 4
- Less than timestamp: lastOccurredAt < 1470394841148

**Logical Operators:**
- AND: message:(error AND critical)
- OR: message:(error OR warning)
- NOT: NOT status:closed
- Parentheses for grouping: (message:error OR description:critical) AND status:open

**Complex Query Examples:**
- Multiple conditions: message:error AND count >= 3
- Grouped conditions: (message:error OR message:warning) AND status:open
- Status with count: status:open AND (count >= 3 OR entity:database)
- Negation: NOT message:test AND status:open

## Wildcards

**Rules:**
- Use * only at the END of words
- Works for: message:error* (matches "error", "errors", "error123")
- Doesn't work for: message:*error or message:err*or
- Not supported for: teams, users (use full names)

**Examples:**
- message:database* (matches "database", "databases", "database_error")
- source:app* (matches "app1", "application", "app_server")

## Null Value Queries

**Check for empty/missing fields:**
- owner:null (alerts without owner)
- teams is null (no teams assigned)
- details.key is not null (has custom details)
- tag !: null (has at least one tag)

**Supported null check fields:**
source, entity, tag, actions, owner, teams, acknowledgedBy, closedBy, recipients, details.key, details.value, integration.name, integration.type

## Common Query Patterns

**Find open alerts:** status:open
**Find high-priority alerts:** message:(critical OR high OR urgent)
**Find unassigned alerts:** owner:null AND status:open
**Find recent alerts:** createdAt > 1640995200000
**Find alerts by team:** teams:infrastructure
**Find alerts with tags:** tag !: null
**Find alerts without acknowledgment:** acknowledgedBy:null AND status:open

## Tips for AI Usage

1. Always use exact field names as listed above
2. Use quotes for multi-word values: message:"database connection error"
3. Combine multiple conditions with AND/OR for precise filtering
4. Use parentheses to group complex conditions
5. Remember that wildcards only work at the end of words
6. Status field only accepts "open" or "closed" as values`

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
