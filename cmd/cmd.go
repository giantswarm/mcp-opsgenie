package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/spf13/cobra"

	"github.com/giantswarm/mcp-opsgenie/pkg/mcp"
)

// cmd defines the root command for the MCP OpsGenie server
var cmd = &cobra.Command{
	Use:   name,
	Short: "MCP server providing access to OpsGenie alerts",
	Long: `An MCP (Model Context Protocol) server that connects to OpsGenie's API.
This server enables AI assistants and other MCP clients to interact with your OpsGenie
instance through a standardized protocol.

The server requires an OpsGenie API token to authenticate with the service.`,
	RunE: runner,
}

var (
	// name is the application name used throughout the server
	name    = "mcp-opsgenie"
	version = "0.1.0"

	// apiURL is the OpsGenie API endpoint URL, defaults to the official API URL
	apiURL = string(client.API_URL)
	// envVar is the name of the environment variable containing the OpsGenie API token
	envVar = "OPSGENIE_TOKEN"
	// logFile is the path to the log file
	logFile = ""
)

// init initializes command line flags for the application
func init() {
	cmd.Flags().StringVar(&apiURL, "api-url", apiURL, "Base URL for the OpsGenie API endpoint")
	cmd.Flags().StringVar(&envVar, "token-env-var", envVar, "Name of environment variable containing your OpsGenie API token")
	cmd.Flags().StringVar(&logFile, "log-file", "", "Path to log file (logs is disabled if not specified)")

}

// Execute runs the root command and handles any execution errors
func Execute() {
	err := cmd.Execute()
	if err != nil {
		slog.Error("execution error", "error", err)
		os.Exit(1)
	}
}

// runner is the main execution function that sets up logging, creates the MCP server,
// registers the OpsGenie handler, and starts the stdio server
func runner(c *cobra.Command, args []string) (err error) {
	// Set up logging - default to discard handler (no output)
	logger := slog.DiscardHandler

	// If a log file is specified, create/open it and use it for logging
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		logger = slog.NewTextHandler(file, nil)
	}

	// Set the default logger for the application
	slog.SetDefault(slog.New(logger))
	slog.Info("Starting MCP OpsGenie server", "version", version, "api_url", apiURL)

	// Create a new MCP server instance
	s := server.NewMCPServer(
		name,
		version,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	// Register the OpsGenie handler with the MCP server
	err = mcp.RegisterOpsGenieHandler(s, apiURL, envVar)
	if err != nil {
		return err
	}

	slog.Info("Initialized MCP server successfully, waiting for client connections...")

	// Start the stdio server to handle MCP protocol communication
	err = server.ServeStdio(s)
	if err != nil {
		// If the server was stopped by user (context canceled), exit gracefully
		if errors.Is(err, context.Canceled) {
			slog.Info("MCP OpsGenie server shutdown requested by user")
			return nil
		}

		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}
