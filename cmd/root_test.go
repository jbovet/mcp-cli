package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	t.Run("command initialization", func(t *testing.T) {
		assert.Equal(t, "mcp-cli", rootCmd.Use)
		assert.Equal(t, "A CLI for interacting with the MCP Registry Service", rootCmd.Short)

		// Test long description contains key features
		assert.True(t, strings.Contains(rootCmd.Long, "List registered MCP servers"))
		assert.True(t, strings.Contains(rootCmd.Long, "Get detailed information"))
		assert.True(t, strings.Contains(rootCmd.Long, "Check service health"))

		// Test examples
		assert.True(t, strings.Contains(rootCmd.Example, "mcp-cli get servers"))
		assert.True(t, strings.Contains(rootCmd.Example, "mcp-cli get server 123e4567"))
		assert.True(t, strings.Contains(rootCmd.Example, "mcp-cli health"))
	})

	t.Run("persistent flags", func(t *testing.T) {
		urlFlag := rootCmd.PersistentFlags().Lookup("url")
		assert.NotNil(t, urlFlag)
		assert.Equal(t, "http://localhost:8080", urlFlag.DefValue)

		verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
		assert.NotNil(t, verboseFlag)
		assert.Equal(t, "false", verboseFlag.DefValue)
		assert.Equal(t, "v", verboseFlag.Shorthand)
	})

	t.Run("execute", func(t *testing.T) {
		// Reset flags to default values before test
		baseURL = "http://localhost:8080"
		verbose = false

		err := Execute()
		assert.NoError(t, err)
	})
}

func TestFlagValues(t *testing.T) {
	t.Run("custom url flag", func(t *testing.T) {
		rootCmd.PersistentFlags().Set("url", "http://custom:9090")
		assert.Equal(t, "http://custom:9090", baseURL)
	})

	t.Run("verbose flag", func(t *testing.T) {
		rootCmd.PersistentFlags().Set("verbose", "true")
		assert.True(t, verbose)
	})
}
