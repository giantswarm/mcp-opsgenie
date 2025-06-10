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
	logFileName = ""
)

func init() {
	cmd.Flags().StringVar(&logFileName, "log-file", "", "File to write logs to (default is stdout)")

}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		slog.Error("execution error", "error", err)
		os.Exit(1)
	}
}

func runner(c *cobra.Command, args []string) (err error) {
	logFile := os.Stdout

	if logFileName != "" {
		logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer logFile.Close()
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, nil)))
	slog.Info("starting mcp-opsgenie server")

	return nil
}
