package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giantswarm/mcp-toolkit/health"
	"github.com/giantswarm/mcp-toolkit/httpx"
	"github.com/giantswarm/mcp-toolkit/logging"
	"github.com/giantswarm/mcp-toolkit/middleware/responsecap"
	"github.com/giantswarm/mcp-toolkit/middleware/timeout"
	"github.com/giantswarm/mcp-toolkit/tracing"
	"github.com/mark3labs/mcp-go/server"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/spf13/cobra"

	"github.com/giantswarm/mcp-opsgenie/pkg/mcp"
)

const (
	serverName              = "mcp-opsgenie"
	transportStdio          = "stdio"
	transportSSE            = "sse"
	transportStreamableHTTP = "streamable-http"
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

	cmd.Flags().StringVar(&apiURL, "api-url", string(client.API_URL), "Base URL for the OpsGenie API endpoint")
	cmd.Flags().StringVar(&envVar, "token-env-var", "OPSGENIE_TOKEN", "Name of environment variable containing your OpsGenie API token")
	cmd.Flags().StringVar(&logFile, "log-file", "", "Path to log file. If empty: logs go to stderr for HTTP transports, are discarded for stdio.")

	cmd.Flags().StringVar(&transport, "transport", transportStdio,
		fmt.Sprintf("Transport type: %s, %s, or %s", transportStdio, transportSSE, transportStreamableHTTP))
	cmd.Flags().StringVar(&httpAddr, "http-addr", ":8080", "HTTP server address (for sse and streamable-http transports)")
	cmd.Flags().StringVar(&sseEndpoint, "sse-endpoint", "/sse", "SSE endpoint path (for sse transport)")
	cmd.Flags().StringVar(&messageEndpoint, "message-endpoint", "/message", "Message endpoint path (for sse transport)")
	cmd.Flags().StringVar(&httpEndpoint, "http-endpoint", "/mcp", "HTTP endpoint path (for streamable-http transport)")

	return cmd
}

// runServeWithVersion contains the main server logic with support for multiple transports and explicit version
func runServeWithVersion(apiURL, envVar, logFile, transport, httpAddr, sseEndpoint, messageEndpoint, httpEndpoint, version string) error {
	shutdownCtx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger, closeLog, err := buildLogger(transport, logFile)
	if err != nil {
		return err
	}
	defer closeLog()
	slog.SetDefault(logger)

	logger.Info("starting MCP OpsGenie server", "version", version, "api_url", apiURL, "transport", transport)

	shutdownOTEL, err := tracing.Init(shutdownCtx, serverName, version)
	if err != nil {
		logger.Warn("otel init failed; continuing without tracing", "error", err)
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = shutdownOTEL(ctx)
		}()
	}

	mcpSrv := server.NewMCPServer(
		serverName,
		version,
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
		server.WithToolHandlerMiddleware(timeout.New(30*time.Second)),
		server.WithToolHandlerMiddleware(responsecap.New(responsecap.Options{})),
	)

	if err := mcp.RegisterOpsGenieHandler(mcpSrv, apiURL, envVar); err != nil {
		return err
	}

	logger.Info("MCP server initialized; awaiting client connections")

	hc := health.New()

	switch transport {
	case transportStdio:
		return runStdioServer(mcpSrv, logger)
	case transportSSE:
		hc.SetReady(true)
		return runSSEServer(shutdownCtx, mcpSrv, hc, httpAddr, sseEndpoint, messageEndpoint, logger)
	case transportStreamableHTTP:
		hc.SetReady(true)
		return runStreamableHTTPServer(shutdownCtx, mcpSrv, hc, httpAddr, httpEndpoint, logger)
	default:
		return fmt.Errorf("unsupported transport type: %s (supported: %s, %s, %s)",
			transport, transportStdio, transportSSE, transportStreamableHTTP)
	}
}

// buildLogger returns the slog logger and a close function. For stdio with
// no --log-file, logs are discarded (stdio drives MCP protocol on
// stdin/stdout; stderr is safe but quiet by default for parity with the
// pre-toolkit behaviour). For HTTP transports, logs default to stderr via
// the toolkit's logging.New (JSON in-cluster, text locally). When
// --log-file is set it overrides both, writing to that file.
func buildLogger(transport, logFile string) (*slog.Logger, func(), error) {
	noop := func() {}
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			return nil, noop, err
		}
		return logging.New(logging.Options{Output: f}), func() { _ = f.Close() }, nil
	}
	if transport == transportStdio {
		return slog.New(slog.DiscardHandler), noop, nil
	}
	return logging.New(logging.Options{}), noop, nil
}

// runStdioServer runs the server with STDIO transport.
func runStdioServer(mcpSrv *server.MCPServer, logger *slog.Logger) error {
	if err := server.ServeStdio(mcpSrv); err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("MCP OpsGenie server shutdown requested")
			return nil
		}
		return fmt.Errorf("server stopped with error: %w", err)
	}
	logger.Info("server stopped normally")
	return nil
}

// runSSEServer runs the SSE transport on its own *http.Server with /healthz
// and /readyz mounted on the same mux.
func runSSEServer(ctx context.Context, mcpSrv *server.MCPServer, hc *health.Health, addr, sseEndpoint, messageEndpoint string, logger *slog.Logger) error {
	sseServer := server.NewSSEServer(mcpSrv,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
	)
	mux := http.NewServeMux()
	mux.Handle(sseEndpoint, sseServer)
	mux.Handle(messageEndpoint, sseServer)
	hc.Mount(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Info("SSE server starting", "addr", addr, "sse", sseEndpoint, "message", messageEndpoint)
	return httpx.Run(ctx, srv, 30*time.Second)
}

// runStreamableHTTPServer runs the streamable-HTTP transport on its own
// *http.Server with /healthz and /readyz mounted on the same mux.
func runStreamableHTTPServer(ctx context.Context, mcpSrv *server.MCPServer, hc *health.Health, addr, endpoint string, logger *slog.Logger) error {
	streamServer := server.NewStreamableHTTPServer(mcpSrv,
		server.WithEndpointPath(endpoint),
	)
	mux := http.NewServeMux()
	mux.Handle(endpoint, streamServer)
	hc.Mount(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Info("streamable-HTTP server starting", "addr", addr, "endpoint", endpoint)
	return httpx.Run(ctx, srv, 30*time.Second)
}
