package cmd

import (
	"os"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/spf13/cobra"
)

// Variables for root command flags (same as serve command for backwards compatibility)
var (
	// OpsGenie configuration
	rootApiURL  string
	rootEnvVar  string
	rootLogFile string

	// Transport options
	rootTransport       string
	rootHttpAddr        string
	rootSseEndpoint     string
	rootMessageEndpoint string
	rootHttpEndpoint    string
)

// rootCmd represents the base command for the mcp-opsgenie application.
// It is the entry point when the application is called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "mcp-opsgenie",
	Short: "MCP server providing access to OpsGenie alerts",
	Long: `An MCP (Model Context Protocol) server that connects to OpsGenie's API.
This server enables AI assistants and other MCP clients to interact with your OpsGenie
instance through a standardized protocol.

The server requires an OpsGenie API token to authenticate with the service.

When run without subcommands, it starts the MCP server (equivalent to 'mcp-opsgenie serve').`,
	// SilenceUsage prevents Cobra from printing the usage message on errors that are handled by the application.
	// This is useful for providing cleaner error output to the user.
	SilenceUsage: true,
}

// SetVersion sets the version for the root command.
// This function is typically called from the main package to inject the application version at build time.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute is the main entry point for the CLI application.
// It initializes and executes the root command, which in turn handles subcommands and flags.
// This function is called by main.main().
func Execute() {
	// SetVersionTemplate defines a custom template for displaying the version.
	// This is used when the --version flag is invoked.
	rootCmd.SetVersionTemplate(`{{printf "mcp-opsgenie version %s\n" .Version}}`)

	// Check if no subcommand was provided and run serve logic (backwards compatibility)
	if len(os.Args) == 1 {
		// Run serve logic directly with root command flag values
		err := runServeWithVersion(rootApiURL, rootEnvVar, rootLogFile, rootTransport, rootHttpAddr, rootSseEndpoint, rootMessageEndpoint, rootHttpEndpoint, rootCmd.Version)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	err := rootCmd.Execute()
	if err != nil {
		// Cobra itself usually prints the error. Exiting with a non-zero status code
		// indicates that an error occurred during execution.
		os.Exit(1)
	}
}

// init is a special Go function that is executed when the package is initialized.
// It is used here to add subcommands to the root command and define flags.
func init() {
	// Add subcommands
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newSelfUpdateCmd())
	rootCmd.AddCommand(newServeCmd())

	// Add flags to root command for backwards compatibility (same as serve command)
	rootCmd.Flags().StringVar(&rootApiURL, "api-url", string(client.API_URL), "Base URL for the OpsGenie API endpoint")
	rootCmd.Flags().StringVar(&rootEnvVar, "token-env-var", "OPSGENIE_TOKEN", "Name of environment variable containing your OpsGenie API token")
	rootCmd.Flags().StringVar(&rootLogFile, "log-file", "", "Path to log file (logs is disabled if not specified)")

	// Transport flags
	rootCmd.Flags().StringVar(&rootTransport, "transport", "stdio", "Transport type: stdio, sse, or streamable-http")
	rootCmd.Flags().StringVar(&rootHttpAddr, "http-addr", ":8080", "HTTP server address (for sse and streamable-http transports)")
	rootCmd.Flags().StringVar(&rootSseEndpoint, "sse-endpoint", "/sse", "SSE endpoint path (for sse transport)")
	rootCmd.Flags().StringVar(&rootMessageEndpoint, "message-endpoint", "/message", "Message endpoint path (for sse transport)")
	rootCmd.Flags().StringVar(&rootHttpEndpoint, "http-endpoint", "/mcp", "HTTP endpoint path (for streamable-http transport)")
}
