package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	t.Run("command initialization", func(t *testing.T) {
		// Test basic command properties
		assert.Equal(t, "get", getCmd.Use)
		assert.Equal(t, "Get resources from the MCP Registry", getCmd.Short)

		// Test long description contains key information
		assert.True(t, strings.Contains(getCmd.Long, "Get various resources"))
		assert.True(t, strings.Contains(getCmd.Long, "servers: List all registered servers"))
		assert.True(t, strings.Contains(getCmd.Long, "server:  Get details of a specific server by ID"))
	})

	t.Run("command examples", func(t *testing.T) {
		// Test that examples contain correct usage patterns
		assert.True(t, strings.Contains(getCmd.Example, "mcp-cli get servers"))
		assert.True(t, strings.Contains(getCmd.Example, "mcp-cli get server"))
		assert.True(t, strings.Contains(getCmd.Example, "123e4567-e89b-12d3-a456-426614174000"))
	})

	t.Run("root command integration", func(t *testing.T) {
		// Verify get command is a subcommand of root
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd == getCmd {
				found = true
				break
			}
		}
		assert.True(t, found, "get command should be registered with root command")
	})

	t.Run("command execution", func(t *testing.T) {
		// Ensure command can be executed without error
		// Note: This just tests the base command without subcommands
		err := getCmd.Execute()
		assert.NoError(t, err)
	})
}
