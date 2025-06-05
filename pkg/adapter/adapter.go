package adapter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ServerAdapter provides an abstraction layer for connecting to MCP servers
// using different transport mechanisms (stdio, HTTP, etc.)
type ServerAdapter interface {
	// Connect establishes a connection to the MCP server
	Connect(ctx context.Context) error

	// Disconnect closes the connection to the MCP server
	Disconnect() error

	// GetServerInfo returns information about the connected server
	GetServerInfo() (*mcp.Implementation, error)

	// ListTools returns available tools from the server
	ListTools(ctx context.Context) ([]mcp.Tool, error)

	// CallTool executes a tool on the server
	CallTool(ctx context.Context, name string, arguments map[string]any) (*mcp.CallToolResult, error)

	// ListResources returns available resources from the server
	ListResources(ctx context.Context) ([]mcp.Resource, error)

	// ReadResource reads a specific resource
	ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error)

	// ListPrompts returns available prompts from the server
	ListPrompts(ctx context.Context) ([]mcp.Prompt, error)

	// GetPrompt retrieves a specific prompt
	GetPrompt(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error)

	// IsConnected returns whether the adapter is currently connected
	IsConnected() bool
}

// Config holds configuration for server adapters
type Config struct {
	// ServerURL for HTTP-based connections
	ServerURL string

	// Command and arguments for stdio-based connections
	Command string
	Args    []string
	Env     []string

	// Connection timeout
	Timeout time.Duration

	// Verbose logging
	Verbose bool
}

// AdapterType represents the type of adapter
type AdapterType string

const (
	AdapterTypeStdio      AdapterType = "stdio"
	AdapterTypeHTTP       AdapterType = "http"
	AdapterTypeStreamable AdapterType = "streamable"
)

// NewAdapter creates a new server adapter based on the configuration
func NewAdapter(adapterType AdapterType, config Config) (ServerAdapter, error) {
	switch adapterType {
	case AdapterTypeStdio:
		return NewStdioAdapter(config)
	case AdapterTypeHTTP, AdapterTypeStreamable:
		return NewHTTPAdapter(config)
	default:
		return nil, fmt.Errorf("unsupported adapter type: %s", adapterType)
	}
}

// BaseAdapter provides common functionality for all adapters
type BaseAdapter struct {
	config     Config
	connected  bool
	serverInfo *mcp.Implementation
}

func (b *BaseAdapter) IsConnected() bool {
	return b.connected
}

func (b *BaseAdapter) GetServerInfo() (*mcp.Implementation, error) {
	if !b.connected {
		return nil, fmt.Errorf("not connected to server")
	}
	return b.serverInfo, nil
}

func (b *BaseAdapter) setConnected(connected bool) {
	b.connected = connected
}

func (b *BaseAdapter) setServerInfo(info *mcp.Implementation) {
	b.serverInfo = info
}

func (b *BaseAdapter) logf(format string, args ...any) {
	if b.config.Verbose {
		log.Printf(format, args...)
	}
}
