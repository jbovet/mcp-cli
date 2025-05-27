# MCP Registry CLI

A command-line interface for interacting with the Model Context Protocol (MCP) Registry Service.

## Features

- **List Servers**: Browse all registered MCP servers with pagination support
- **Server Details**: Get comprehensive information about specific servers
- **Health Checks**: Monitor service health and status
- **Service Ping**: Check service version and availability
- **Table Output**: Clean, formatted output using tables
- **Error Handling**: Comprehensive error handling with informative messages
- **Verbose Mode**: Optional verbose output for debugging

## Installation

### Prerequisites

- Go 1.21 or later

### Build from Source

```bash
# Install dependencies
make deps

# Build the binary
make build

# The binary will be available at bin/mcp-cli

Install directly
bashmake install
Usage
Basic Commands
bash# Check service health
mcp-cli health

# Ping the service
mcp-cli ping

# List all servers
mcp-cli get servers

# Get details for a specific server
mcp-cli get server <server-id>
Global Options

--url: Base URL of the MCP Registry Service (default: http://localhost:8080)
--verbose, -v: Enable verbose output

Examples
bash# Use a different service URL
mcp-cli --url https://registry.example.com health

# List servers with pagination
mcp-cli get servers --limit 10
mcp-cli get servers --cursor 123e4567-e89b-12d3-a456-426614174000

# Enable verbose output
mcp-cli -v get servers

# Get detailed server information
mcp-cli get server 123e4567-e89b-12d3-a456-426614174000