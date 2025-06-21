// Package opsgenie provides a client for interacting with the OpsGenie Alert API.
// It offers functionality to retrieve alerts with pagination support and query filtering.
package opsgenie

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/sirupsen/logrus"
)

const (
	// maxAlertsPerRequest is the maximum number of alerts that can be fetched in a single API request.
	// This limit is enforced by the OpsGenie API.
	// Reference: https://docs.opsgenie.com/docs/alert-api#list-alerts
	maxAlertsPerRequest = 100

	// maxTotalAlerts is the maximum total number of alerts that can be fetched across all paginated requests.
	// This limit is enforced by the OpsGenie API.
	// Reference: https://docs.opsgenie.com/docs/alert-api#list-alerts
	maxTotalAlerts = 20000
)

// AlertClient is a wrapper around the OpsGenie alert client that provides
// enhanced functionality for fetching and managing alerts.
type AlertClient struct {
	*alert.Client
}

// NewAlertClient creates a new AlertClient instance configured with the provided API URL and API key.
// The API key is retrieved from the environment variable specified by envVar.
//
// Parameters:
//   - apiUrl: The OpsGenie API URL endpoint (e.g., "https://api.opsgenie.com")
//   - envVar: The name of the environment variable containing the API key
//
// Returns:
//   - *AlertClient: A configured alert client ready for use
//   - error: An error if the client creation fails or if the API key is missing
func NewAlertClient(apiUrl, envVar string) (*AlertClient, error) {
	logger := logrus.New()
	logger.Out = io.Discard

	config := &client.Config{
		OpsGenieAPIURL: client.ApiUrl(apiUrl),
		ApiKey:         os.Getenv(envVar),
		Logger:         logger,
	}

	alertClient, err := alert.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsGenie alert client: %w", err)
	}

	a := &AlertClient{
		Client: alertClient,
	}

	return a, nil
}

// ListAlerts retrieves alerts from OpsGenie based on the provided query string.
// The method handles pagination automatically, fetching all matching alerts up to the maximum limit.
//
// The query parameter supports OpsGenie's query syntax for filtering alerts.
// Examples:
//   - "status:open" - fetch only open alerts
//   - "tag:critical" - fetch alerts with the "critical" tag
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - query: OpsGenie query string for filtering alerts (empty string fetches all alerts)
//
// Returns:
//   - []alert.Alert: A slice of alerts matching the query criteria
//   - error: An error if the API request fails or if the context is cancelled
func (a *AlertClient) ListAlerts(ctx context.Context, query string) ([]alert.Alert, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	alerts := make([]alert.Alert, 0, maxAlertsPerRequest)
	offset := 0

	slog.Info("fetching alerts",
		"query", query,
		"max_per_request", maxAlertsPerRequest,
		"max_total", maxTotalAlerts)

	// Paginate through all available alerts until we reach the limit or no more alerts exist
	for offset < maxTotalAlerts {
		// Prepare the list request with pagination parameters
		listRequest := &alert.ListAlertRequest{
			Offset: offset,
			Limit:  maxAlertsPerRequest,
			Sort:   alert.CreatedAt, // Sort by creation time
			Order:  alert.Desc,      // Most recent alerts first
			Query:  query,
		}

		response, err := a.Client.List(ctx, listRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to list alerts: %w", err)
		}

		// If no alerts are returned, we've reached the end of available data
		if len(response.Alerts) == 0 {
			break
		}

		// Append the fetched alerts to our result set
		alerts = append(alerts, response.Alerts...)
		offset += maxAlertsPerRequest
	}

	slog.Info("fetched alerts", "count", len(alerts))

	return alerts, nil
}
