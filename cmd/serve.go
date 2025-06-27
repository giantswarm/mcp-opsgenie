package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/giantswarm/mcp-opsgenie/pkg/mcp"
	"github.com/mark3labs/mcp-go/server"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

// newServeCmd creates the Cobra command for starting the MCP server.
func newServeCmd() *cobra.Command {
	var (
		// OpsGenie configuration
		apiURL  string
		envVar  string
		logFile string

		// Transport options
		transport       string
		httpAddr        string
		sseEndpoint     string
		messageEndpoint string
		httpEndpoint    string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP OpsGenie server",
		Long: `Start the MCP OpsGenie server to provide tools for interacting
with OpsGenie alerts, teams, and heartbeats via the Model Context Protocol.

Supports multiple transport types:
  - stdio: Standard input/output (default)
  - sse: Server-Sent Events over HTTP
  - streamable-http: Streamable HTTP transport

The server requires an OpsGenie API token to authenticate with the service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServeWithVersion(apiURL, envVar, logFile, transport, httpAddr, sseEndpoint, messageEndpoint, httpEndpoint, cmd.Root().Version)
		},
	}

	// Add flags for OpsGenie configuration
	cmd.Flags().StringVar(&apiURL, "api-url", string(client.API_URL), "Base URL for the OpsGenie API endpoint")
	cmd.Flags().StringVar(&envVar, "token-env-var", "OPSGENIE_TOKEN", "Name of environment variable containing your OpsGenie API token")
	cmd.Flags().StringVar(&logFile, "log-file", "", "Path to log file (logs is disabled if not specified)")

	// Transport flags
	cmd.Flags().StringVar(&transport, "transport", "stdio", "Transport type: stdio, sse, or streamable-http")
	cmd.Flags().StringVar(&httpAddr, "http-addr", ":8080", "HTTP server address (for sse and streamable-http transports)")
	cmd.Flags().StringVar(&sseEndpoint, "sse-endpoint", "/sse", "SSE endpoint path (for sse transport)")
	cmd.Flags().StringVar(&messageEndpoint, "message-endpoint", "/message", "Message endpoint path (for sse transport)")
	cmd.Flags().StringVar(&httpEndpoint, "http-endpoint", "/mcp", "HTTP endpoint path (for streamable-http transport)")

	return cmd
}

// runServeWithVersion contains the main server logic with support for multiple transports and explicit version
func runServeWithVersion(apiURL, envVar, logFile, transport, httpAddr, sseEndpoint, messageEndpoint, httpEndpoint, version string) error {
	// Setup graceful shutdown - listen for both SIGINT and SIGTERM
	shutdownCtx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
	mcpSrv := server.NewMCPServer(
		"mcp-opsgenie",
		version, // Use version parameter instead of rootCmd.Version
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	// Register the OpsGenie handler with the MCP server
	err := mcp.RegisterOpsGenieHandler(mcpSrv, apiURL, envVar)
	if err != nil {
		return err
	}

	slog.Info("Initialized MCP server successfully, waiting for client connections...")

	fmt.Printf("Starting MCP OpsGenie server with %s transport...\n", transport)

	// Start the appropriate server based on transport type
	switch transport {
	case "stdio":
		return runStdioServer(mcpSrv)
	case "sse":
		return runSSEServer(mcpSrv, httpAddr, sseEndpoint, messageEndpoint, shutdownCtx)
	case "streamable-http":
		return runStreamableHTTPServer(mcpSrv, httpAddr, httpEndpoint, shutdownCtx)
	default:
		return fmt.Errorf("unsupported transport type: %s (supported: stdio, sse, streamable-http)", transport)
	}
}

// runStdioServer runs the server with STDIO transport
func runStdioServer(mcpSrv *mcpserver.MCPServer) error {
	// Start the server in a goroutine so we can handle shutdown signals
	serverDone := make(chan error, 1)
	go func() {
		defer close(serverDone)
		if err := mcpserver.ServeStdio(mcpSrv); err != nil {
			serverDone <- err
		}
	}()

	// Wait for server completion
	select {
	case err := <-serverDone:
		if err != nil {
			// If the server was stopped by user (context canceled), exit gracefully
			if errors.Is(err, context.Canceled) {
				slog.Info("MCP OpsGenie server shutdown requested by user")
				return nil
			}
			return fmt.Errorf("server stopped with error: %w", err)
		} else {
			fmt.Println("Server stopped normally")
		}
	}

	fmt.Println("Server gracefully stopped")
	return nil
}

// runSSEServer runs the server with SSE transport
func runSSEServer(mcpSrv *mcpserver.MCPServer, addr, sseEndpoint, messageEndpoint string, ctx context.Context) error {
	// Create SSE server with custom endpoints
	sseServer := mcpserver.NewSSEServer(mcpSrv,
		mcpserver.WithSSEEndpoint(sseEndpoint),
		mcpserver.WithMessageEndpoint(messageEndpoint),
	)

	fmt.Printf("SSE server starting on %s\n", addr)
	fmt.Printf("  SSE endpoint: %s\n", sseEndpoint)
	fmt.Printf("  Message endpoint: %s\n", messageEndpoint)

	// Start server in goroutine
	serverDone := make(chan error, 1)
	go func() {
		defer close(serverDone)
		if err := sseServer.Start(addr); err != nil {
			serverDone <- err
		}
	}()

	// Wait for either shutdown signal or server completion
	select {
	case <-ctx.Done():
		fmt.Println("Shutdown signal received, stopping SSE server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()
		if err := sseServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down SSE server: %w", err)
		}
	case err := <-serverDone:
		if err != nil {
			return fmt.Errorf("SSE server stopped with error: %w", err)
		} else {
			fmt.Println("SSE server stopped normally")
		}
	}

	fmt.Println("SSE server gracefully stopped")
	return nil
}

// runStreamableHTTPServer runs the server with Streamable HTTP transport
func runStreamableHTTPServer(mcpSrv *mcpserver.MCPServer, addr, endpoint string, ctx context.Context) error {
	// Create Streamable HTTP server with custom endpoint
	httpServer := mcpserver.NewStreamableHTTPServer(mcpSrv,
		mcpserver.WithEndpointPath(endpoint),
	)

	fmt.Printf("Streamable HTTP server starting on %s\n", addr)
	fmt.Printf("  HTTP endpoint: %s\n", endpoint)

	// Start server in goroutine
	serverDone := make(chan error, 1)
	go func() {
		defer close(serverDone)
		if err := httpServer.Start(addr); err != nil {
			serverDone <- err
		}
	}()

	// Wait for either shutdown signal or server completion
	select {
	case <-ctx.Done():
		fmt.Println("Shutdown signal received, stopping HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down HTTP server: %w", err)
		}
	case err := <-serverDone:
		if err != nil {
			return fmt.Errorf("HTTP server stopped with error: %w", err)
		} else {
			fmt.Println("HTTP server stopped normally")
		}
	}

	fmt.Println("HTTP server gracefully stopped")
	return nil
}
