package main

import "github.com/giantswarm/mcp-opsgenie/cmd"

const version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
