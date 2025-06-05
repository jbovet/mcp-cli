package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapterFactory(t *testing.T) {
	factory := NewAdapterFactory()

	t.Run("CreateStdioAdapter", func(t *testing.T) {
		config := map[string]interface{}{
			"type":    "stdio",
			"command": "echo",
			"args":    []string{"hello"},
			"verbose": true,
		}

		adapter, err := factory.CreateFromConfig(config)
		require.NoError(t, err)
		assert.NotNil(t, adapter)

		stdioAdapter, ok := adapter.(*StdioAdapter)
		assert.True(t, ok)
		assert.Equal(t, "echo", stdioAdapter.config.Command)
		assert.Equal(t, []string{"hello"}, stdioAdapter.config.Args)
		assert.True(t, stdioAdapter.config.Verbose)
	})

	t.Run("CreateHTTPAdapter", func(t *testing.T) {
		config := map[string]interface{}{
			"type": "http",
			"url":  "http://localhost:8080/mcp",
		}

		adapter, err := factory.CreateFromConfig(config)
		require.NoError(t, err)
		assert.NotNil(t, adapter)

		httpAdapter, ok := adapter.(*HTTPAdapter)
		assert.True(t, ok)
		assert.Equal(t, "http://localhost:8080/mcp", httpAdapter.config.ServerURL)
	})

	t.Run("CreateFromURL_HTTP", func(t *testing.T) {
		adapter, err := factory.CreateFromURL("http://localhost:8080/mcp", false)
		require.NoError(t, err)
		assert.NotNil(t, adapter)

		httpAdapter, ok := adapter.(*HTTPAdapter)
		assert.True(t, ok)
		assert.Equal(t, "http://localhost:8080/mcp", httpAdapter.config.ServerURL)
	})

	t.Run("CreateFromURL_Command", func(t *testing.T) {
		adapter, err := factory.CreateFromURL("python server.py --port 8080", true)
		require.NoError(t, err)
		assert.NotNil(t, adapter)

		stdioAdapter, ok := adapter.(*StdioAdapter)
		assert.True(t, ok)
		assert.Equal(t, "python", stdioAdapter.config.Command)
		assert.Equal(t, []string{"server.py", "--port", "8080"}, stdioAdapter.config.Args)
		assert.True(t, stdioAdapter.config.Verbose)
	})

	t.Run("ValidateConfig", func(t *testing.T) {
		// Valid stdio config
		config := Config{
			Command: "python",
			Args:    []string{"server.py"},
			Timeout: 30 * time.Second,
		}
		err := factory.ValidateConfig(AdapterTypeStdio, config)
		assert.NoError(t, err)

		// Invalid stdio config - missing command
		config.Command = ""
		err = factory.ValidateConfig(AdapterTypeStdio, config)
		assert.Error(t, err)

		// Valid HTTP config
		config = Config{
			ServerURL: "http://localhost:8080",
			Timeout:   30 * time.Second,
		}
		err = factory.ValidateConfig(AdapterTypeHTTP, config)
		assert.NoError(t, err)

		// Invalid HTTP config - missing URL
		config.ServerURL = ""
		err = factory.ValidateConfig(AdapterTypeHTTP, config)
		assert.Error(t, err)

		// Invalid HTTP config - bad URL
		config.ServerURL = "not-a-url"
		err = factory.ValidateConfig(AdapterTypeHTTP, config)
		assert.Error(t, err)
	})

	t.Run("GetSupportedTypes", func(t *testing.T) {
		types := factory.GetSupportedTypes()
		assert.Len(t, types, 3)
		assert.Contains(t, types, AdapterTypeStdio)
		assert.Contains(t, types, AdapterTypeHTTP)
		assert.Contains(t, types, AdapterTypeStreamable)
	})
}

func TestBaseAdapter(t *testing.T) {
	base := BaseAdapter{
		config: Config{Verbose: true},
	}

	t.Run("InitialState", func(t *testing.T) {
		assert.False(t, base.IsConnected())
		_, err := base.GetServerInfo()
		assert.Error(t, err)
	})

	t.Run("SetConnected", func(t *testing.T) {
		base.setConnected(true)
		assert.True(t, base.IsConnected())

		base.setConnected(false)
		assert.False(t, base.IsConnected())
	})
}

func TestStdioAdapter(t *testing.T) {
	t.Run("NewStdioAdapter", func(t *testing.T) {
		config := Config{
			Command: "echo",
			Args:    []string{"hello"},
			Timeout: 30 * time.Second,
		}

		adapter, err := NewStdioAdapter(config)
		require.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, "echo", adapter.config.Command)
		assert.False(t, adapter.IsConnected())
	})

	t.Run("NewStdioAdapter_MissingCommand", func(t *testing.T) {
		config := Config{
			Timeout: 30 * time.Second,
		}

		_, err := NewStdioAdapter(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "command is required")
	})

	t.Run("Connect_NotConnected", func(t *testing.T) {
		config := Config{
			Command: "nonexistent-command",
			Timeout: 1 * time.Second,
		}

		adapter, err := NewStdioAdapter(config)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = adapter.Connect(ctx)
		assert.Error(t, err) // Should fail because command doesn't exist
		assert.False(t, adapter.IsConnected())
	})
}

func TestHTTPAdapter(t *testing.T) {
	t.Run("NewHTTPAdapter", func(t *testing.T) {
		config := Config{
			ServerURL: "http://localhost:8080",
			Timeout:   30 * time.Second,
		}

		adapter, err := NewHTTPAdapter(config)
		require.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, "http://localhost:8080", adapter.config.ServerURL)
		assert.False(t, adapter.IsConnected())
	})

	t.Run("NewHTTPAdapter_MissingURL", func(t *testing.T) {
		config := Config{
			Timeout: 30 * time.Second,
		}

		_, err := NewHTTPAdapter(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server URL is required")
	})

	t.Run("Connect_InvalidURL", func(t *testing.T) {
		config := Config{
			ServerURL: "http://localhost:9999", // Non-existent server
			Timeout:   1 * time.Second,
		}

		adapter, err := NewHTTPAdapter(config)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err = adapter.Connect(ctx)
		assert.Error(t, err) // Should fail because server doesn't exist
		assert.False(t, adapter.IsConnected())
	})
}

// Integration test helpers
func TestAdapterIntegration(t *testing.T) {
	// Skip integration tests in CI unless specifically enabled
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	t.Run("StdioEcho", func(t *testing.T) {
		// Test with a simple echo command that should work on most systems
		config := Config{
			Command: "echo",
			Args:    []string{"test"},
			Timeout: 5 * time.Second,
		}

		adapter, err := NewStdioAdapter(config)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Note: This will likely fail because echo is not an MCP server
		// but it tests the basic connection mechanism
		err = adapter.Connect(ctx)
		if err != nil {
			t.Logf("Expected failure for echo command: %v", err)
		}

		// Always disconnect to clean up
		adapter.Disconnect()
	})
}
