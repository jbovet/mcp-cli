package cmd

import (
	"fmt"
	"os"

	"github.com/jbovet/mcp-cli/pkg/client"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health status of the MCP Registry Service",
	Long: `Check the health status of the MCP Registry Service.

This command performs a health check against the service and displays
the current status along with any additional health information.`,
	Example: `  # Check service health
  mcp-cli health`,
	RunE: runHealthCommand,
}

func runHealthCommand(cmd *cobra.Command, args []string) error {
	// Create API client
	apiClient := client.NewClient(baseURL)

	// Perform health check
	health, err := apiClient.GetHealth()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Health check completed successfully\n")
	}

	// Display health status
	fmt.Printf("Service Health Status\n")
	fmt.Printf("====================\n\n")
	fmt.Printf("Status: %s\n", health.Status)

	if health.GitHubClientID != "" {
		fmt.Printf("GitHub Client ID: %s\n", health.GitHubClientID)
	}

	if health.Status == "ok" {
		fmt.Printf("\n✓ Service is healthy and operational\n")
	} else {
		fmt.Printf("\n✗ Service health check indicates issues\n")
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
