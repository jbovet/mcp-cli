package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Global flags
	baseURL string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcp-cli",
	Short: "A CLI for interacting with the MCP Registry Service",
	Long: `MCP Registry CLI is a command-line interface for fetching and displaying
information from the Model Context Protocol Registry Service.

This CLI allows you to:
- List registered MCP servers
- Get detailed information about specific servers
- Check service health and status`,
	Example: `  # List all servers
  mcp-cli get servers

  # Get details for a specific server
  mcp-cli get server 123e4567-e89b-12d3-a456-426614174000

  # Check service health
  mcp-cli health`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "http://localhost:8080", "Base URL of the MCP Registry Service")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
