package adapter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// StdioAdapter implements ServerAdapter for stdio-based MCP servers
type StdioAdapter struct {
	BaseAdapter
	client mcpclient.MCPClient
}

// NewStdioAdapter creates a new stdio adapter
func NewStdioAdapter(config Config) (*StdioAdapter, error) {
	if config.Command == "" {
		return nil, fmt.Errorf("command is required for stdio adapter")
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &StdioAdapter{
		BaseAdapter: BaseAdapter{
			config: config,
		},
	}, nil
}

// Connect establishes a stdio connection to the MCP server
func (s *StdioAdapter) Connect(ctx context.Context) error {
	if s.connected {
		return fmt.Errorf("already connected")
	}

	s.logf("Connecting to MCP server via stdio: %s %v", s.config.Command, s.config.Args)
	client, err := mcpclient.NewStdioMCPClient(
		s.config.Command,
		s.config.Env,
		s.config.Args...,
	)
	if err != nil {
		return fmt.Errorf("failed to create stdio client: %w", err)
	}

	// Log client for debugging
	s.logf("Created stdio client: %+v", client)

	s.client = client

	// Wait and check if process is still alive
	if err := s.waitForProcessReady(ctx); err != nil {
		if closeErr := s.client.Close(); closeErr != nil {
			s.logf("Warning: failed to close stdio client during cleanup: %v", closeErr)
		}
		return err
	}

	// Initialize the connection
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-cli-adapter",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{
		Roots: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: true,
		},
	}

	s.logf("Sending initialize request with timeout: %v", s.config.Timeout)
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()

	result, err := s.client.Initialize(ctx, initRequest)
	if err != nil {
		s.logf("Initialize failed: %v", err)
		if err := s.client.Close(); err != nil {
			// Log the error but don't return it since this is likely in a cleanup context
			fmt.Fprintf(os.Stderr, "Warning: failed to close stdio client: %v\n", err)
		}
		return fmt.Errorf("failed to initialize: %w", err)
	}

	s.setConnected(true)
	s.setServerInfo(&result.ServerInfo)
	s.logf("Successfully connected to server: %s %s", result.ServerInfo.Name, result.ServerInfo.Version)

	return nil
}

// Disconnect closes the stdio connection
func (s *StdioAdapter) Disconnect() error {
	if !s.connected {
		return nil
	}

	s.logf("Disconnecting from MCP server")
	err := s.client.Close()
	s.setConnected(false)
	s.setServerInfo(nil)
	return err
}

// ListTools returns available tools from the server
func (s *StdioAdapter) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListToolsRequest{}
	result, err := s.client.ListTools(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	return result.Tools, nil
}

// CallTool executes a tool on the server
func (s *StdioAdapter) CallTool(ctx context.Context, name string, arguments map[string]any) (*mcp.CallToolResult, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = name
	request.Params.Arguments = arguments

	result, err := s.client.CallTool(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", name, err)
	}

	return result, nil
}

// ListResources returns available resources from the server
func (s *StdioAdapter) ListResources(ctx context.Context) ([]mcp.Resource, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListResourcesRequest{}
	result, err := s.client.ListResources(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return result.Resources, nil
}

// ReadResource reads a specific resource
func (s *StdioAdapter) ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ReadResourceRequest{}
	request.Params.URI = uri

	result, err := s.client.ReadResource(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource %s: %w", uri, err)
	}

	return result, nil
}

// ListPrompts returns available prompts from the server
func (s *StdioAdapter) ListPrompts(ctx context.Context) ([]mcp.Prompt, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.ListPromptsRequest{}
	result, err := s.client.ListPrompts(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	return result.Prompts, nil
}

// GetPrompt retrieves a specific prompt
func (s *StdioAdapter) GetPrompt(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	request := mcp.GetPromptRequest{}
	request.Params.Name = name
	request.Params.Arguments = arguments

	result, err := s.client.GetPrompt(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt %s: %w", name, err)
	}

	return result, nil
}

func (s *StdioAdapter) waitForProcessReady(ctx context.Context) error {
	s.logf("Waiting for process to be ready...")

	// Check multiple times with increasing delays
	delays := []time.Duration{100 * time.Millisecond, 500 * time.Millisecond, 1 * time.Second}

	for i, delay := range delays {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		s.logf("Process readiness check %d/%d", i+1, len(delays))

		// Try a quick ping to see if process responds
		pingCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		err := s.pingProcess(pingCtx)
		cancel()

		if err == nil {
			s.logf("Process is ready!")
			return nil
		}

		if isProcessExitError(err) {
			return fmt.Errorf("process exited unexpectedly - check command '%s %v'",
				s.config.Command, s.config.Args)
		}

		s.logf("Process not ready yet (attempt %d): %v", i+1, err)
	}

	return fmt.Errorf("process did not become ready within expected time")
}

func (s *StdioAdapter) pingProcess(ctx context.Context) error {
	// Send a minimal request to test if the process is alive and responding
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "ping",
		Version: "1.0.0",
	}

	_, err := s.client.Initialize(ctx, initRequest)
	return err
}
func isProcessExitError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "process") ||
		strings.Contains(errStr, "exit") ||
		strings.Contains(errStr, "pipe") ||
		strings.Contains(errStr, "broken") ||
		strings.Contains(errStr, "EOF")
}
