package cmd

import (
	"fmt"

	"github.com/jbovet/mcp-cli/pkg/client"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Check the service status and version",
	Long: `Ping the MCP Registry Service to check its status and get version information.

This command sends a ping request to the service and displays the response,
including the service version and operational status.`,
	Example: `  # Ping the service
  mcp-cli ping`,
	RunE: runPingCommand,
}

func runPingCommand(cmd *cobra.Command, args []string) error {
	// Create API client
	apiClient := client.NewClient(baseURL)

	// Ping the service
	ping, err := apiClient.Ping()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	if verbose {
		fmt.Printf("Ping completed successfully\n")
	}

	// Display ping response
	fmt.Printf("Service Ping Response\n")
	fmt.Printf("====================\n\n")
	fmt.Printf("Status:  %s\n", ping.Status)
	fmt.Printf("Version: %s\n", ping.Version)

	if ping.Status == "ok" {
		fmt.Printf("\n✓ Service is responding normally\n")
	} else {
		fmt.Printf("\n✗ Service ping indicates issues\n")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
