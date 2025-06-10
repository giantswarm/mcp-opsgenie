package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "mcp-opsgenie",
	Short: "A Model Context Protocol (MCP) server for OpsGenie",
	Long: `mcp-opsgenie is a Model Context Protocol (MCP) server that provides
access to OpsGenie. This server is designed to be used with the MCP framework.`,
	RunE: runner,
}

var (
	logFile = ""
)

func init() {
	cmd.Flags().StringVar(&logFile, "log-file", "", "File to write logs to (log disabled by default)")

}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		slog.Error("execution error", "error", err)
		os.Exit(1)
	}
}

func runner(c *cobra.Command, args []string) (err error) {
	logger := slog.DiscardHandler

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		logger = slog.NewTextHandler(file, nil)
	}

	slog.SetDefault(slog.New(logger))
	slog.Info("starting mcp-opsgenie server")

	return nil
}
