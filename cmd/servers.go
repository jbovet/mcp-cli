package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jbovet/mcp-cli/pkg/client"
	"github.com/spf13/cobra"
)

var (
	// Flags for servers command
	limit  int
	cursor string
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "List all registered MCP servers",
	Long: `List all registered MCP servers from the registry with pagination support.

This command fetches and displays a list of all MCP servers registered in the service,
showing their basic information including ID, name, description, and version.`,
	Example: `  # List servers with default limit (30)
  mcp-cli get servers

  # List servers with custom limit
  mcp-cli get servers --limit 10

  # Continue from a specific cursor
  mcp-cli get servers --cursor 123e4567-e89b-12d3-a456-426614174000`,
	RunE: runServersCommand,
}

func runServersCommand(cmd *cobra.Command, args []string) error {
	// Create API client
	apiClient := client.NewClient(baseURL)

	// Fetch servers from API
	response, err := apiClient.GetServers(cursor, limit)
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Fetched %d servers\n", len(response.Servers))
	}

	// Display results in table format
	if err := displayServersTable(response); err != nil {
		return fmt.Errorf("failed to display servers: %w", err)
	}

	return nil
}

func displayServersTable(response *client.ServersResponse) error {
	if len(response.Servers) == 0 {
		fmt.Println("No servers found.")
		return nil
	}

	// Create table writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintln(w, "ID\tNAME\tVERSION\tDESCRIPTION")
	fmt.Fprintln(w, "---\t----\t-------\t-----------")

	// Print each server
	for _, server := range response.Servers {
		description := server.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			server.ID,
			server.Name,
			server.VersionDetail.Version,
			description,
		)
	}

	// Print pagination info if available
	if response.Metadata.NextCursor != "" {
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Next cursor: %s\n", response.Metadata.NextCursor)
		fmt.Fprintf(w, "Use --cursor flag to continue pagination\n")
	}

	return nil
}

func init() {
	getCmd.AddCommand(serversCmd)
	serversCmd.Flags().IntVar(&limit, "limit", 30, "Maximum number of servers to return (1-100)")
	serversCmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor (UUID)")
}
