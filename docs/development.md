## Development

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/giantswarm/mcp-opsgenie.git
   cd mcp-opsgenie
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up your environment:
   ```bash
   export OPSGENIE_TOKEN="your-development-token"
   ```

4. Run the server:
   ```bash
   go run main.go --log-file "debug.log"
   ```

### Project Structure

```
├── cmd/                    # Command-line interface
│   └── cmd.go             # Main command implementation
├── pkg/
│   ├── mcp/               # MCP protocol handlers
│   │   └── opsgenie_handler.go
│   └── opsgenie/          # OpsGenie API client
├── docs/                  # Documentation
├── main.go               # Application entry point
├── go.mod                # Go module definition
└── Makefile              # Build automation
```

### Adding New Features

1. **New Tools**: Add new MCP tools in `pkg/mcp/opsgenie_handler.go`
2. **API Methods**: Extend the OpsGenie client in `pkg/opsgenie/`
3. **Configuration**: Add new command-line flags in `cmd/cmd.go`

### Testing

Run tests with:

```bash
go test ./...
```

### Building

Build the binary:

```bash
go build -o mcp-opsgenie
```

## Security

- **API Token**: Store your OpsGenie API token securely using environment variables
- **Permissions**: Ensure your API token has only the necessary permissions
- **Network**: Run the server in a secure network environment
- **Logging**: Be cautious about logging sensitive information

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/giantswarm/mcp-opsgenie/issues)
- **Documentation**: Check the [docs/](docs/) directory for additional documentation
- **OpsGenie API**: Refer to the [OpsGenie API Documentation](https://docs.opsgenie.com/docs/api-overview)

## Acknowledgments

- Built with [mcp-go](https://github.com/mark3labs/mcp-go) - Go implementation of the Model Context Protocol
- Uses [opsgenie-go-sdk-v2](https://github.com/opsgenie/opsgenie-go-sdk-v2) - Official OpsGenie Go SDK
- Follows the [Model Context Protocol](https://github.com/modelcontextprotocol) specification