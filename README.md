<div align="center">

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/mcp-opsgenie/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/mcp-opsgenie/tree/main)
[![GoDoc](https://pkg.go.dev/badge/github.com/giantswarm/mcp-opsgenie.svg)](https://pkg.go.dev/github.com/giantswarm/mcp-opsgenie)

<strong>OpsGenie MCP Server</strong>

*A Model Context Protocol server that enables AI assistants to interact with OpsGenie*

</div>

## Overview

The OpsGenie MCP Server is a [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol) server that provides AI assistants and other MCP clients with standardized access to OpsGenie. This server acts as a bridge between AI tools and your OpsGenie instance, allowing for automated management of alerts, teams, and heartbeats through natural language interactions.

## Features

- **Alert Management**: List, get, acknowledge, and unacknowledge OpsGenie alerts.
- **Team Management**: List and get details for teams.
- **Heartbeat Monitoring**: List and get the status of heartbeats.
- **Powerful Alert Filtering**: Utilize advanced search queries to filter alerts.
- **MCP Compliance**: Fully compatible with the Model Context Protocol for seamless integration.
- **Secure Authentication**: Uses OpsGenie API tokens for secure access.
- **Flexible Configuration**: Configurable API endpoints and logging options.

## Prerequisites

- Go 1.24.2 or later
- An OpsGenie API token with appropriate permissions

## Installation

### Using Go Install

```bash
go install github.com/giantswarm/mcp-opsgenie@latest
```

### Building from Source

```bash
git clone https://github.com/giantswarm/mcp-opsgenie.git
cd mcp-opsgenie
go build -o mcp-opsgenie
```


## Configuration

## Integration with AI Assistants

This MCP server can be integrated with various AI assistants that support the Model Context Protocol:

- **Cursor**: The server can be integrated with Cursor for AI-powered development.
- **VSCode Insiders**: Compatible with VSCode Insiders for enhanced coding assistance.
- **Claude Desktop**: Add the server to your Claude Desktop configuration
- **Custom MCP Clients**: Use any MCP-compatible client to connect

### Example MCP Client Configuration

```json
{
  "servers": {
    "opsgenie": {
      "command": "/path/to/mcp-opsgenie",
      "env": {
        "OPSGENIE_TOKEN": "your-api-token-here"
      }
    }
  }
}
```

### Environment Variables

Set your OpsGenie API token as an environment variable:

```bash
export OPSGENIE_TOKEN="your-opsgenie-api-token-here"
```

## Usage

### Basic Usage

Start the MCP server with default settings:

```bash
mcp-opsgenie
```

### Advanced Usage

The server supports several command-line options:

```bash
$ mcp-opsgenie --help
An MCP (Model Context Protocol) server that connects to OpsGenie's API.
This server enables AI assistants and other MCP clients to interact with your OpsGenie
instance through a standardized protocol.

The server requires an OpsGenie API token to authenticate with the service.

Usage:
  mcp-opsgenie [flags]

Flags:
      --api-url string         Base URL for the OpsGenie API endpoint (default "api.opsgenie.com")
  -h, --help                   help for mcp-opsgenie
      --log-file string        Path to log file (logs is disabled if not specified)
      --token-env-var string   Name of environment variable containing your OpsGenie API token (default "OPSGENIE_TOKEN")
```

### Custom Configuration Examples

```bash
# Use a custom API endpoint
mcp-opsgenie --api-url "https://api.eu.opsgenie.com"

# Use a different environment variable for the token
mcp-opsgenie --token-env-var "CUSTOM_OPSGENIE_TOKEN"

# Enable logging to a file
mcp-opsgenie --log-file "mcp-opsgenie.log"
```

## Available Tools

### `list_alerts`

Retrieve a list of alerts from OpsGenie using advanced search queries.

**Parameters:**
- `query` (optional): Search query for filtering alerts

For comprehensive query documentation, see the [OpsGenie Search Documentation](https://support.atlassian.com/opsgenie/docs/search-queries-for-alerts/).

### `get_alert`

Retrieves a single alert from OpsGenie using its ID.

**Parameters:**
- `id`: Identifier of the alert to be retrieved.

### `acknowledge_alert`

Acknowledges an alert in OpsGenie.

**Parameters:**
- `id`: Identifier of the alert to acknowledge.
- `note` (optional): Note to add to the alert.
- `user` (optional): Display name of the request owner.
- `source` (optional): Display name of the request source.

### `unacknowledge_alert`

Unacknowledges an alert in OpsGenie.

**Parameters:**
- `id`: Identifier of the alert to unacknowledge.
- `note` (optional): Note to add to the alert.
- `user` (optional): Display name of the request owner.
- `source` (optional): Display name of the request source.

### `list_teams`

Retrieve a list of all teams from OpsGenie.

### `get_team`

Retrieves a single team from OpsGenie by its ID or name.

**Parameters:**
- `identifier`: Name or ID of the team to retrieve.
- `identifier_type` (optional): Type of the identifier. Possible values are 'id' and 'name'. Defaults to 'id'.

### `list_heartbeats`

Retrieve a list of all heartbeats from OpsGenie.

### `get_heartbeat`

Retrieves a single heartbeat from OpsGenie by its name.

**Parameters:**
- `name`: Name of the heartbeat to retrieve.
