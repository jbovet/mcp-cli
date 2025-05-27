# MCP CLI

A command-line interface for interacting with the MCP Registry Service.

## Features

- List registered MCP servers
- Get detailed server information by ID or name
- Check service health and status
- Service ping functionality
- Pagination support for listing servers
- Pattern matching for server searches

## Installation

Clone the repository and build using Make:

```sh
git clone https://github.com/jbovet/mcp-cli.git
cd mcp-cli
make build
```

The binary will be created in `bin/mcp-cli`.

## Usage

### Basic Commands

```sh
# List all servers
mcp-cli get servers

# Get server details by ID
mcp-cli get server 123e4567-e89b-12d3-a456-426614174000

# Get server details by name
mcp-cli get server io.github.owner/server-name

# Search servers by name pattern
mcp-cli get server redis --name --show-matches

# Check service health
mcp-cli health

# Ping service
mcp-cli ping
```

### Global Flags

- `--url`: Base URL of the MCP Registry Service (default: http://localhost:8080)
- `--verbose, -v`: Enable verbose output

### Server Listing Options

- `--limit`: Maximum number of servers to return (1-100)
- `--cursor`: Pagination cursor for fetching next page

## Development

### Requirements

- Go 1.24.3 or later
- Make

### Common Tasks

```sh
# Install dependencies
make deps

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

### Project Structure

```
cmd/        - Command implementations
pkg/        - Core packages
  client/   - API client implementation
  models/   - Data models
bin/        - Build output
```

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.