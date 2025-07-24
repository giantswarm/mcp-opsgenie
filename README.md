<div align="center">

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/mcp-opsgenie/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/mcp-opsgenie/tree/main)
[![GoDoc](https://pkg.go.dev/badge/github.com/giantswarm/mcp-opsgenie.svg)](https://pkg.go.dev/github.com/giantswarm/mcp-opsgenie)

<strong>OpsGenie MCP Server</strong>

*A Model Context Protocol server that enables AI assistants to interact with OpsGenie alerts*

</div>

## Overview

The OpsGenie MCP Server is a [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol) server that provides AI assistants and other MCP clients with standardized access to OpsGenie's alert management system. This server acts as a bridge between AI tools and your OpsGenie instance, allowing for automated alert querying and management through natural language interactions.

## Features

- **Alert Querying**: Retrieve and filter OpsGenie alerts using powerful search queries
- **Comprehensive Search**: Support for all OpsGenie alert fields and advanced search operators
- **MCP Compliance**: Fully compatible with the Model Context Protocol standard
- **Secure Authentication**: Uses OpsGenie API tokens for secure access
- **Flexible Configuration**: Configurable API endpoints and logging options

## Prerequisites

- Go 1.24.2 or later
- An OpsGenie API token with appropriate permissions

|Tool|OpsGenie Permission|
|-----|----------|
|`list_alerts`|Read|
|`get_alert`|Read|
|`acknowledge_alert`|Update|
|`unacknowledge_alert`|Update|
|`list_heartbeats`|Read|
|`get_heartbeat`|Read|
|`list_teams`|Read|
|`get_team`|Read|


## Installation

### Building from Source

```bash
git clone https://github.com/giantswarm/mcp-opsgenie.git
cd mcp-opsgenie
go build -o mcp-opsgenie
```

### Using Go Install

```bash
go install github.com/giantswarm/mcp-opsgenie@latest
```

## Configuration

### Environment Variables

Set your OpsGenie API token as an environment variable:

```bash
export OPSGENIE_TOKEN="your-opsgenie-api-token-here"
```

## Usage

### Basic Usage

Start the MCP server with default settings:

```bash
./mcp-opsgenie
```

### Advanced Usage

The server supports several command-line options:

```bash
./mcp-opsgenie --help

MCP server providing access to OpsGenie alerts

An MCP (Model Context Protocol) server that connects to OpsGenie's API.
This server enables AI assistants and other MCP clients to interact with your OpsGenie
instance through a standardized protocol.

The server requires an OpsGenie API token to authenticate with the service.

Usage:
  mcp-opsgenie [flags]

Flags:
      --api-url string         Base URL for the OpsGenie API endpoint (default "https://api.opsgenie.com")
  -h, --help                   help for mcp-opsgenie
      --log-file string        Path to log file (logs is disabled if not specified)
      --token-env-var string   Name of environment variable containing your OpsGenie API token (default "OPSGENIE_TOKEN")
```

### Custom Configuration Examples

```bash
# Use a custom API endpoint
./mcp-opsgenie --api-url "https://api.eu.opsgenie.com"

# Use a different environment variable for the token
./mcp-opsgenie --token-env-var "CUSTOM_OPSGENIE_TOKEN"

# Enable logging to a file
./mcp-opsgenie --log-file "mcp-opsgenie.log"
```

## Available Tools

### `list_alerts`

Retrieve a list of alerts from OpsGenie using advanced search queries.

**Parameters:**
- `query` (optional): Search query for filtering alerts

**Example Queries:**

```
# Get all open alerts
status:open

# Find critical alerts
message:(critical OR high OR urgent)

# Find unassigned alerts
owner:null AND status:open

# Find alerts from the last 24 hours (using timestamp)
createdAt > 1640995200000

# Find database-related alerts
message:database* OR entity:database

# Complex query with multiple conditions
(message:error OR message:warning) AND status:open AND teams:infrastructure
```

For comprehensive query documentation, see the [OpsGenie Search Documentation](https://support.atlassian.com/opsgenie/docs/search-queries-for-alerts/).

## Integration with AI Assistants

This MCP server can be integrated with various AI assistants that support the Model Context Protocol:

1. **Claude Desktop**: Add the server to your Claude Desktop configuration
2. **Custom MCP Clients**: Use any MCP-compatible client to connect
3. **Development Tools**: Integrate with IDEs and development environments that support MCP

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
