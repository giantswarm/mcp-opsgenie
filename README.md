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
- **Multi-Transport Support**: Connect via stdio, Server-Sent Events (SSE), or Streamable HTTP.
- **Enhanced CLI**: Version management, self-update capability, and comprehensive help system.
- **MCP Compliance**: Fully compatible with the Model Context Protocol for seamless integration.
- **Secure Authentication**: Uses OpsGenie API tokens for secure access.
- **Flexible Configuration**: Configurable API endpoints, transport options, and logging.
- **Backwards Compatible**: Maintains compatibility with existing configurations and scripts.

## Prerequisites

- Go 1.24.4 or later
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

### Self-Update

The server includes a built-in self-update mechanism:

```bash
mcp-opsgenie self-update
```

## Configuration

### Environment Variables

Set your OpsGenie API token as an environment variable:

```bash
export OPSGENIE_TOKEN="your-opsgenie-api-token-here"
```

## Usage

### CLI Commands

The MCP server provides several commands:

```bash
$ mcp-opsgenie --help
An MCP (Model Context Protocol) server that connects to OpsGenie's API.
This server enables AI assistants and other MCP clients to interact with your OpsGenie
instance through a standardized protocol.

The server requires an OpsGenie API token to authenticate with the service.

When run without subcommands, it starts the MCP server (equivalent to 'mcp-opsgenie serve').

Usage:
  mcp-opsgenie [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  self-update Update mcp-opsgenie to the latest version
  serve       Start the MCP OpsGenie server
  version     Print the version number of mcp-opsgenie

Flags:
      --api-url string            Base URL for the OpsGenie API endpoint (default "api.opsgenie.com")
  -h, --help                      help for mcp-opsgenie
      --http-addr string          HTTP server address (for sse and streamable-http transports) (default ":8080")
      --http-endpoint string      HTTP endpoint path (for streamable-http transport) (default "/mcp")
      --log-file string           Path to log file (logs is disabled if not specified)
      --message-endpoint string   Message endpoint path (for sse transport) (default "/message")
      --sse-endpoint string       SSE endpoint path (for sse transport) (default "/sse")
      --token-env-var string      Name of environment variable containing your OpsGenie API token (default "OPSGENIE_TOKEN")
      --transport string          Transport type: stdio, sse, or streamable-http (default "stdio")
  -v, --version                   version for mcp-opsgenie

Use "mcp-opsgenie [command] --help" for more information about a command.
```

### Basic Usage (Backwards Compatible)

Start the MCP server with default settings using stdio transport:

```bash
# Both commands are equivalent and start the server
mcp-opsgenie
mcp-opsgenie serve
```

### Multi-Transport Support

The server supports three transport types for different deployment scenarios:

#### Standard I/O (Default)
Best for MCP client integrations and development:

```bash
mcp-opsgenie serve --transport stdio
```

#### Server-Sent Events (SSE)
Ideal for web applications and browser-based clients:

```bash
mcp-opsgenie serve --transport sse --http-addr :8080
```

#### Streamable HTTP
Perfect for HTTP-based integrations and REST-like interactions:

```bash
mcp-opsgenie serve --transport streamable-http --http-addr :8080
```

### Advanced Configuration Examples

```bash
# Use a custom API endpoint
mcp-opsgenie --api-url "https://api.eu.opsgenie.com"

# Use a different environment variable for the token
mcp-opsgenie --token-env-var "CUSTOM_OPSGENIE_TOKEN"

# Enable logging to a file
mcp-opsgenie --log-file "mcp-opsgenie.log"

# Run with SSE transport on custom port with custom endpoints
mcp-opsgenie serve \
  --transport sse \
  --http-addr :9090 \
  --sse-endpoint /events \
  --message-endpoint /messages

# Run with streamable HTTP transport
mcp-opsgenie serve \
  --transport streamable-http \
  --http-addr :8080 \
  --http-endpoint /api/mcp
```

### Version Management

```bash
# Check current version
mcp-opsgenie version
mcp-opsgenie --version

# Update to latest version
mcp-opsgenie self-update
```

## Integration with AI Assistants

This MCP server can be integrated with various AI assistants that support the Model Context Protocol:

- **Cursor**: The server can be integrated with Cursor for AI-powered development.
- **VSCode Insiders**: Compatible with VSCode Insiders for enhanced coding assistance.
- **Claude Desktop**: Add the server to your Claude Desktop configuration
- **Custom MCP Clients**: Use any MCP-compatible client to connect

### Example MCP Client Configuration

#### Standard I/O Transport (Recommended)
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

#### SSE Transport
```json
{
  "servers": {
    "opsgenie": {
      "url": "http://localhost:8080/sse",
      "transport": "sse"
    }
  }
}
```

#### Streamable HTTP Transport
```json
{
  "servers": {
    "opsgenie": {
      "url": "http://localhost:8080/mcp",
      "transport": "http"
    }
  }
}
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

## Deployment

### Docker

You can run the server in a Docker container:

```dockerfile
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mcp-opsgenie

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mcp-opsgenie .
EXPOSE 8080
CMD ["./mcp-opsgenie", "serve", "--transport", "sse", "--http-addr", ":8080"]
```

### Systemd Service

Create a systemd service for automatic startup:

```ini
[Unit]
Description=OpsGenie MCP Server
After=network.target

[Service]
Type=simple
User=mcp
ExecStart=/usr/local/bin/mcp-opsgenie serve --log-file /var/log/mcp-opsgenie.log
Environment=OPSGENIE_TOKEN=your-token-here
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Development

### Building

```bash
go build -v .
```

### Testing

```bash
make test
```

### Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to submit issues, feature requests, and pull requests.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

- Create an issue in this repository
- Check the [OpsGenie API documentation](https://docs.opsgenie.com/docs/api-overview)
- Review the [Model Context Protocol specification](https://github.com/modelcontextprotocol)
