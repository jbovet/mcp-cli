package adapter

import (
	"context"
	"fmt"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// HTTPAdapter implements ServerAdapter for HTTP-based MCP servers
type HTTPAdapter struct {
	BaseAdapter
	client mcpclient.MCPClient
}

// NewHTTPAdapter creates a new HTTP adapter
func NewHTTPAdapter(config Config) (*HTTPAdapter, error) {
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server URL is required for HTTP adapter")
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &HTTPAdapter{
		BaseAdapter: BaseAdapter{
			config: config,
		},
	}, nil
}

// Connect establishes an HTTP connection to the MCP server
func (h *HTTPAdapter) Connect(ctx context.Context) error {
	if h.connected {
		return fmt.Errorf("already connected")
	}

	h.logf("Connecting to MCP server via HTTP: %s", h.config.ServerURL)

	// Create streamable HTTP client
	client, err := mcpclient.NewStreamableHttpClient(h.config.ServerURL)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	h.client = client

	ctx, cancel := context.WithTimeout(ctx, h.config.Timeout)
	defer cancel()

	// Initialize the connection
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-cli-adapter",
		Version: "1.0.0",
	}

	result, err := h.client.Initialize(ctx, initRequest)
	if err != nil {
		h.client.Close()
		return fmt.Errorf("failed to initialize: %w", err)
	}

	h.setConnected(true)
	h.setServerInfo(&result.ServerInfo)
	h.logf("Successfully connected to server: %s %s", result.ServerInfo.Name, result.ServerInfo.Version)

	return nil
}

// Disconnect closes the HTTP connection
func (h *HTTPAdapter) Disconnect() error {
	if !h.connected {
		return nil
	}

	h.logf("Disconnecting from MCP server")
	err := h.client.Close()
	h.setConnected(false)
	h.setServerInfo(nil)
	return err
}

// ListTools returns available tools from the server
func (h *HTTPAdapter) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListToolsRequest{}
	result, err := h.client.ListTools(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return result.Tools, nil
}

// CallTool executes a tool on the server
func (h *HTTPAdapter) CallTool(ctx context.Context, name string, arguments map[string]any) (*mcp.CallToolResult, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = name
	request.Params.Arguments = arguments

	result, err := h.client.CallTool(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", name, err)
	}

	return result, nil
}

// ListResources returns available resources from the server
func (h *HTTPAdapter) ListResources(ctx context.Context) ([]mcp.Resource, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListResourcesRequest{}
	result, err := h.client.ListResources(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return result.Resources, nil
}

// ReadResource reads a specific resource
func (h *HTTPAdapter) ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ReadResourceRequest{}
	request.Params.URI = uri

	result, err := h.client.ReadResource(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource %s: %w", uri, err)
	}

	return result, nil
}

// ListPrompts returns available prompts from the server
func (h *HTTPAdapter) ListPrompts(ctx context.Context) ([]mcp.Prompt, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListPromptsRequest{}
	result, err := h.client.ListPrompts(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	return result.Prompts, nil
}

// GetPrompt retrieves a specific prompt
func (h *HTTPAdapter) GetPrompt(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error) {
	if !h.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.GetPromptRequest{}
	request.Params.Name = name
	request.Params.Arguments = arguments

	result, err := h.client.GetPrompt(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt %s: %w", name, err)
	}

	return result, nil
}
