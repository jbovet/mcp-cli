# MCP CLI

A command-line interface for interacting with the [MCP Registry Service](https://github.com/modelcontextprotocol/registry) and connecting directly to MCP servers.

## Features

- **Registry Service Integration:**
  - List registered MCP servers
  - Get detailed server information by ID or name
  - Check service health and status
  - Service ping functionality
  - Pagination support for listing servers
  - Pattern matching for server searches

- **Direct MCP Server Connection:**
  - Connect to MCP servers via stdio transport
  - Connect to HTTP-based MCP servers
  - Interactive mode for real-time server exploration
  - Tool execution with argument parsing
  - Resource reading capabilities
  - Prompt listing and interaction

## Installation

Clone the repository and build using Make:

```sh
git clone https://github.com/jbovet/mcp-cli.git
cd mcp-cli
make build
```

The binary will be created in `bin/mcp-cli`.

## Usage

### Registry Service Commands

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

### Direct MCP Server Connection

#### Stdio Transport (Local Processes)

```sh
# Connect to a Python MCP server
mcp-cli connect --type stdio --command "python" --args "server.py"

# Connect to a Node.js MCP server with environment variables
mcp-cli connect --type stdio --command "node" --args "server.js" --env "DEBUG=1"

# Connect using command string parsing
mcp-cli connect --type stdio --command "python /path/to/server.py --port 8080"

# Interactive mode for exploration
mcp-cli connect --type stdio --command "python server.py" --interactive
```

#### HTTP Transport (Remote Servers)

```sh
# Connect to HTTP MCP server
mcp-cli connect --type http --url "http://localhost:8080/mcp"

# Connect to HTTPS MCP server
mcp-cli connect --type streamable --url "https://api.example.com/mcp"

# Interactive mode with HTTP server
mcp-cli connect --type http --url "http://localhost:8080/mcp" --interactive
```

#### Interactive Mode Commands

When running with `--interactive`, you can use these commands:

```sh
# List server capabilities
tools         # List available tools
resources     # List available resources  
prompts       # List available prompts

# Execute operations
call <tool-name> [args...]    # Call a tool with arguments
read <resource-uri>           # Read a resource by URI

# Navigation
help          # Show help message
quit/exit     # Exit interactive mode
```

#### Example Interactive Session

```sh
> tools
Available tools (3):
1. validate_api - Validate API specifications
2. format_code - Format source code
3. analyze_logs - Analyze log files

> call validate_api openapi.yaml
Tool 'validate_api' result:
✅ Tool executed successfully:
    API specification is valid
    Found 12 endpoints
    No validation errors

> resources
Available resources (2):
1. file://config.json (Configuration) - Server configuration
2. file://schema.yaml (Schema) - API schema definition

> read file://config.json
Resource content:
  Text (application/json): {"server": {"port": 8080}}

> quit
Goodbye!
```

### Global Flags

- `--url`: Base URL of the MCP Registry Service (default: http://localhost:8080)
- `--verbose, -v`: Enable verbose output

### Registry Service Options

- `--limit`: Maximum number of servers to return (1-100)
- `--cursor`: Pagination cursor for fetching next page

### Connect Command Options

- `--type`: Transport type (`stdio`, `http`, `streamable`)
- `--url`: Server URL for HTTP-based connections
- `--command`: Command to execute for stdio connections
- `--args`: Arguments for the command (can be repeated)
- `--env`: Environment variables for the command (can be repeated)
- `--timeout`: Connection timeout (default: 60s)
- `--interactive`: Run in interactive mode

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
  connect.go      - MCP server connection command
  get.go         - Registry "get" command group
  server.go      - Individual server details
  servers.go     - Server listing
  health.go      - Health check command
  ping.go        - Ping command
pkg/        - Core packages
  client/   - Registry API client implementation
  models/   - Data models
  adapter/  - MCP server adapters (stdio, HTTP)
    adapter.go    - Core adapter interfaces
    stdio.go      - Stdio transport implementation
    http.go       - HTTP transport implementation
    factory.go    - Adapter factory and utilities
bin/        - Build output
```

## Use Cases

### Registry Integration
- Discover available MCP servers in the registry
- Get detailed information about server capabilities
- Monitor registry service health

### Direct Server Testing
- Test MCP server implementations during development
- Validate tool and resource functionality
- Interactive debugging and exploration
- Automated testing and validation scripts

### CI/CD Integration
- Health checks in deployment pipelines
- Automated server capability validation
- Registry service monitoring

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.
