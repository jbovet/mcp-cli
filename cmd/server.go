package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jbovet/mcp-cli/pkg/client"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server <id>",
	Short: "Get detailed information about a specific server",
	Long: `Get detailed information about a specific MCP server by its ID.

This command fetches and displays comprehensive information about a server,
including its packages, environment variables, and configuration details.`,
	Example: `  # Get server details by ID
  mcp-cli get server 123e4567-e89b-12d3-a456-426614174000`,
	Args: cobra.ExactArgs(1),
	RunE: runServerCommand,
}

func runServerCommand(cmd *cobra.Command, args []string) error {
	serverID := args[0]

	// Create API client
	apiClient := client.NewClient(baseURL)

	// Fetch server details from API
	server, err := apiClient.GetServer(serverID)
	if err != nil {
		return fmt.Errorf("failed to fetch server details: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Fetched server details for ID: %s\n", serverID)
	}

	// Display server details
	if err := displayServerDetails(server); err != nil {
		return fmt.Errorf("failed to display server details: %w", err)
	}

	return nil
}

func displayServerDetails(server *client.ServerDetail) error {
	fmt.Printf("Server Details\n")
	fmt.Printf("==============\n\n")

	// Basic information
	fmt.Printf("ID:          %s\n", server.ID)
	fmt.Printf("Name:        %s\n", server.Name)
	fmt.Printf("Description: %s\n", server.Description)
	fmt.Printf("Version:     %s\n", server.VersionDetail.Version)
	fmt.Printf("Release Date: %s\n", server.VersionDetail.ReleaseDate)
	fmt.Printf("Is Latest:   %t\n", server.VersionDetail.IsLatest)

	// Repository information
	fmt.Printf("\nRepository\n")
	fmt.Printf("----------\n")
	fmt.Printf("URL:    %s\n", server.Repository.URL)
	fmt.Printf("Source: %s\n", server.Repository.Source)
	if server.Repository.ID != "" {
		fmt.Printf("ID:     %s\n", server.Repository.ID)
	}

	// Packages
	if len(server.Packages) > 0 {
		fmt.Printf("\nPackages\n")
		fmt.Printf("--------\n")
		for i, pkg := range server.Packages {
			fmt.Printf("%d. %s (%s)\n", i+1, pkg.Name, pkg.RegistryName)
			fmt.Printf("   Version: %s\n", pkg.Version)
			if pkg.RuntimeHint != "" {
				fmt.Printf("   Runtime Hint: %s\n", pkg.RuntimeHint)
			}

			// Package arguments
			if len(pkg.PackageArguments) > 0 {
				fmt.Printf("   Package Arguments:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "     TYPE\tNAME\tREQUIRED\tDEFAULT\tDESCRIPTION")
				for _, arg := range pkg.PackageArguments {
					argType := string(arg.Type)
					required := "No"
					if arg.IsRequired {
						required = "Yes"
					}
					fmt.Fprintf(w, "     %s\t%s\t%s\t%s\t%s\n",
						argType, arg.Name, required, arg.Default, arg.Description)
				}
				w.Flush()
			}

			// Runtime arguments
			if len(pkg.RuntimeArguments) > 0 {
				fmt.Printf("   Runtime Arguments:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "     TYPE\tNAME\tREQUIRED\tDEFAULT\tDESCRIPTION")
				for _, arg := range pkg.RuntimeArguments {
					argType := string(arg.Type)
					required := "No"
					if arg.IsRequired {
						required = "Yes"
					}
					fmt.Fprintf(w, "     %s\t%s\t%s\t%s\t%s\n",
						argType, arg.Name, required, arg.Default, arg.Description)
				}
				w.Flush()
			}

			// Environment variables
			if len(pkg.EnvironmentVariables) > 0 {
				fmt.Printf("   Environment Variables:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "     NAME\tREQUIRED\tDESCRIPTION")
				for _, env := range pkg.EnvironmentVariables {
					required := "No"
					if env.IsRequired {
						required = "Yes"
					}
					fmt.Fprintf(w, "     %s\t%s\t%s\n", env.Name, required, env.Description)
				}
				w.Flush()
			}

			if i < len(server.Packages)-1 {
				fmt.Println()
			}
		}
	}

	// Remotes
	if len(server.Remotes) > 0 {
		fmt.Printf("\nRemote Connections\n")
		fmt.Printf("------------------\n")
		for i, remote := range server.Remotes {
			fmt.Printf("%d. Transport: %s\n", i+1, remote.TransportType)
			fmt.Printf("   URL: %s\n", remote.URL)
			if len(remote.Headers) > 0 {
				fmt.Printf("   Headers:\n")
				for _, header := range remote.Headers {
					fmt.Printf("     %s: %s\n", header.Value, header.Description)
				}
			}
			if i < len(server.Remotes)-1 {
				fmt.Println()
			}
		}
	}

	return nil
}

func init() {
	getCmd.AddCommand(serverCmd)
}
