package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCommand(t *testing.T) {
	t.Run("command initialization", func(t *testing.T) {
		assert.Equal(t, "server <id-or-name>", serverCmd.Use)
		assert.Contains(t, serverCmd.Short, "Get detailed information about a specific server")
		assert.Contains(t, serverCmd.Long, "You can specify either")
		assert.Contains(t, serverCmd.Long, "server ID (UUID format)")
		assert.Contains(t, serverCmd.Long, "server name")
	})

	t.Run("flags configuration", func(t *testing.T) {
		nameFlag := serverCmd.Flags().Lookup("name")
		assert.NotNil(t, nameFlag)
		assert.Equal(t, "false", nameFlag.DefValue)

		matchesFlag := serverCmd.Flags().Lookup("show-matches")
		assert.NotNil(t, matchesFlag)
		assert.Equal(t, "false", matchesFlag.DefValue)
	})

	t.Run("examples contain both ID and name usage", func(t *testing.T) {
		assert.Contains(t, serverCmd.Example, "123e4567-e89b-12d3-a456-426614174000")
		assert.Contains(t, serverCmd.Example, "io.github.owner/server-name")
		assert.Contains(t, serverCmd.Example, "--name")
		assert.Contains(t, serverCmd.Example, "--show-matches")
	})
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"123E4567-E89B-12D3-A456-426614174000", true},
		{"io.github.owner/server-name", false},
		{"redis-server", false},
		{"not-a-uuid", false},
		{"123e4567-e89b-12d3-a456", false}, // too short
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := isUUID(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
