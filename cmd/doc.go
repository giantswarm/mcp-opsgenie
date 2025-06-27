// Package cmd contains the command-line interface (CLI) for mcp-opsgenie.
//
// This package implements the Cobra-based command structure that provides
// subcommands for starting the MCP server, checking version information,
// and performing self-updates.
//
// The main commands available are:
//   - serve: Start the MCP OpsGenie server with various transport options
//   - version: Display version information
//   - self-update: Update the application to the latest release from GitHub
//
// The package follows the standard Cobra CLI patterns and provides a clean
// separation between the CLI interface and the core server functionality.
package cmd
