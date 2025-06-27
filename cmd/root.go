package cmd

import (
	"os"

	"github.com/spf13/cobra"
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

	// If no subcommand is provided, run the serve command by default
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "serve")
	}

	err := rootCmd.Execute()
	if err != nil {
		// Cobra itself usually prints the error. Exiting with a non-zero status code
		// indicates that an error occurred during execution.
		os.Exit(1)
	}
}

// init is a special Go function that is executed when the package is initialized.
// It is used here to add subcommands to the root command.
func init() {
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newSelfUpdateCmd())
	rootCmd.AddCommand(newServeCmd())

	// Example of how to define persistent flags (global for the application):
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mcp-opsgenie/config.yaml)")

	// Example of how to define local flags (only run when this action is called directly):
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
