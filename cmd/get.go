package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command group
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from the MCP Registry",
	Long: `Get various resources from the MCP Registry Service.

Available subcommands:
- servers: List all registered servers
- server:  Get details of a specific server by ID`,
	Example: `  # List all servers
  mcp-cli get servers

  # Get server details
  mcp-cli get server 123e4567-e89b-12d3-a456-426614174000`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
