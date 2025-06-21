// Package opsgenie provides a client for interacting with the OpsGenie API.
package opsgenie

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/team"
	"github.com/sirupsen/logrus"
)

// TeamClient is a wrapper around the OpsGenie team client.
type TeamClient struct {
	*team.Client
}

// NewTeamClient creates a new TeamClient instance.
func NewTeamClient(apiUrl, envVar string) (*TeamClient, error) {
	logger := logrus.New()
	logger.Out = io.Discard

	config := &client.Config{
		OpsGenieAPIURL: client.ApiUrl(apiUrl),
		ApiKey:         os.Getenv(envVar),
		Logger:         logger,
	}

	teamClient, err := team.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsGenie team client: %w", err)
	}

	t := &TeamClient{
		Client: teamClient,
	}

	return t, nil
}

// ListTeams retrieves all teams from OpsGenie.
func (c *TeamClient) ListTeams(ctx context.Context) ([]team.ListedTeams, error) {
	result, err := c.Client.List(ctx, &team.ListTeamRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return result.Teams, nil
}

// GetTeam retrieves a single team by its ID or name.
func (c *TeamClient) GetTeam(ctx context.Context, identifier, identifierType string) (*team.GetTeamResult, error) {
	if identifier == "" {
		return nil, fmt.Errorf("team identifier cannot be empty")
	}

	var idType team.Identifier
	switch identifierType {
	case "name":
		idType = team.Name
	default:
		idType = team.Id
	}

	getTeamRequest := &team.GetTeamRequest{
		IdentifierValue: identifier,
		IdentifierType:  idType,
	}

	result, err := c.Client.Get(ctx, getTeamRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get team %s: %w", identifier, err)
	}

	return result, nil
}
