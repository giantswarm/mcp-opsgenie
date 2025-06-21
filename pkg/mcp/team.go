package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (h *opsgenieHandler) registerTeamTools(s *server.MCPServer) {
	listTeamsTool := mcp.NewTool("list_teams",
		mcp.WithDescription("Retrieve a list of all teams from OpsGenie."),
	)
	s.AddTool(listTeamsTool, h.ListTeams)

	getTeamTool := mcp.NewTool("get_team",
		mcp.WithDescription("Retrieves a single team from OpsGenie by its ID or name."),
		mcp.WithString("identifier",
			mcp.Description("Name or ID of the team to retrieve."),
			mcp.Required(),
		),
		mcp.WithString("identifier_type",
			mcp.Description("Type of the identifier. Possible values are 'id' and 'name'. Defaults to 'id'."),
		),
	)
	s.AddTool(getTeamTool, h.GetTeam)
}

// ListTeams retrieves all teams from OpsGenie.
func (h *opsgenieHandler) ListTeams(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	teams, err := h.teamClient.ListTeams(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve teams from OpsGenie: %v", err)), nil
	}

	data, err := json.Marshal(teams)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize teams to JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// GetTeam retrieves a single OpsGenie team by its name or ID.
func (h *opsgenieHandler) GetTeam(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	identifier := request.GetString("identifier", "")
	if identifier == "" {
		return mcp.NewToolResultError("the 'identifier' parameter is required"), nil
	}
	identifierType := request.GetString("identifier_type", "id")

	team, err := h.teamClient.GetTeam(ctx, identifier, identifierType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve team with identifier '%s' from OpsGenie: %v", identifier, err)), nil
	}

	data, err := json.Marshal(team)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize team to JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
