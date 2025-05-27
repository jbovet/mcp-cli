package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/jbovet/mcp-cli/pkg/client"
	"github.com/jbovet/mcp-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// Flags for server command
	serverByName bool
	showMatches  bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server <id-or-name>",
	Short: "Get detailed information about a specific server",
	Long: `Get detailed information about a specific MCP server by its ID or name.

You can specify either:
- A server ID (UUID format): 123e4567-e89b-12d3-a456-426614174000  
- A server name: io.github.owner/server-name
- A name pattern with --name flag for fuzzy matching

The command will automatically detect whether you provided an ID or name.`,
	Example: `  # Get server details by ID
  mcp-cli get server 123e4567-e89b-12d3-a456-426614174000

  # Get server details by exact name
  mcp-cli get server io.github.owner/server-name

  # Search for servers by name pattern
  mcp-cli get server redis --name --show-matches

  # Find servers containing "github" in the name
  mcp-cli get server github --name`,
	Args: cobra.ExactArgs(1),
	RunE: runServerCommand,
}

func runServerCommand(cmd *cobra.Command, args []string) error {
	identifier := args[0]

	// Create API client
	apiClient := client.NewClient(baseURL)

	var server *client.ServerDetail
	var err error

	// Determine if identifier is UUID or name
	if isUUID(identifier) && !serverByName {
		// It's a UUID, get directly by ID
		if verbose {
			fmt.Fprintf(os.Stderr, "Looking up server by ID: %s\n", identifier)
		}
		server, err = apiClient.GetServer(identifier)
		if err != nil {
			return fmt.Errorf("failed to fetch server by ID: %w", err)
		}
	} else {
		// It's a name or forced name lookup
		if showMatches {
			// Show all matches for pattern
			return showServerMatches(apiClient, identifier)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Looking up server by name: %s\n", identifier)
		}

		// Try exact name match first
		server, err = apiClient.GetServerByName(identifier)
		if err != nil {
			// If exact match fails, suggest pattern search
			if strings.Contains(err.Error(), "not found") {
				matches, searchErr := apiClient.FindServersByNamePattern(identifier)
				if searchErr != nil {
					return fmt.Errorf("failed to search for servers: %w", searchErr)
				}

				if len(matches) == 0 {
					return fmt.Errorf("no servers found matching '%s'", identifier)
				}

				if len(matches) == 1 {
					// Only one match, get its details
					if verbose {
						fmt.Fprintf(os.Stderr, "Found single match: %s\n", matches[0].Name)
					}
					server, err = apiClient.GetServer(matches[0].ID)
					if err != nil {
						return fmt.Errorf("failed to fetch server details: %w", err)
					}
				} else {
					// Multiple matches, show them
					fmt.Printf("Multiple servers found matching '%s':\n\n", identifier)
					return displayServerList(matches, "Use exact name or ID to get details")
				}
			} else {
				return fmt.Errorf("failed to fetch server by name: %w", err)
			}
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Successfully fetched server details\n")
	}

	// Display server details
	if err := displayServerDetails(server); err != nil {
		return fmt.Errorf("failed to display server details: %w", err)
	}

	return nil
}

// isUUID checks if a string is in UUID format
func isUUID(s string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(s)
}

// showServerMatches displays all servers that match a pattern
func showServerMatches(apiClient *client.Client, pattern string) error {
	if verbose {
		fmt.Fprintf(os.Stderr, "Searching for servers matching pattern: %s\n", pattern)
	}

	matches, err := apiClient.FindServersByNamePattern(pattern)
	if err != nil {
		return fmt.Errorf("failed to search servers: %w", err)
	}

	if len(matches) == 0 {
		fmt.Printf("No servers found matching pattern '%s'\n", pattern)
		return nil
	}

	fmt.Printf("Found %d server(s) matching '%s':\n\n", len(matches), pattern)
	return displayServerList(matches, "Use exact name or ID to get details")
}

// displayServerList shows a list of servers in table format
func displayServerList(servers []models.Server, footer string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		_ = w.Flush()
	}()

	// Print header
	_, _ = fmt.Fprintln(w, "ID\tNAME\tVERSION\tDESCRIPTION")
	_, _ = fmt.Fprintln(w, "---\t----\t-------\t-----------")

	// Print each server
	for _, server := range servers {
		description := server.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			server.ID,
			server.Name,
			server.VersionDetail.Version,
			description,
		)
	}

	if footer != "" {
		_, _ = fmt.Fprintf(w, "\n%s\n", footer)
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
				_, _ = fmt.Fprintln(w, "     TYPE\tNAME\tREQUIRED\tDEFAULT\tDESCRIPTION")
				for _, arg := range pkg.PackageArguments {
					argType := string(arg.Type)
					required := "No"
					if arg.IsRequired {
						required = "Yes"
					}
					_, _ = fmt.Fprintf(w, "     %s\t%s\t%s\t%s\t%s\n",
						argType, arg.Name, required, arg.Default, arg.Description)
				}
				_ = w.Flush()
			}

			// Runtime arguments
			if len(pkg.RuntimeArguments) > 0 {
				fmt.Printf("   Runtime Arguments:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				_, _ = fmt.Fprintln(w, "     TYPE\tNAME\tREQUIRED\tDEFAULT\tDESCRIPTION")
				for _, arg := range pkg.RuntimeArguments {
					argType := string(arg.Type)
					required := "No"
					if arg.IsRequired {
						required = "Yes"
					}
					_, _ = fmt.Fprintf(w, "     %s\t%s\t%s\t%s\t%s\n",
						argType, arg.Name, required, arg.Default, arg.Description)
				}
				_ = w.Flush()
			}

			// Environment variables
			if len(pkg.EnvironmentVariables) > 0 {
				fmt.Printf("   Environment Variables:\n")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				_, _ = fmt.Fprintln(w, "     NAME\tREQUIRED\tDESCRIPTION")
				for _, env := range pkg.EnvironmentVariables {
					required := "No"
					if env.IsRequired {
						required = "Yes"
					}
					_, _ = fmt.Fprintf(w, "     %s\t%s\t%s\n", env.Name, required, env.Description)
				}
				_ = w.Flush()
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
	serverCmd.Flags().BoolVar(&serverByName, "name", false, "Force name-based lookup (enables pattern matching)")
	serverCmd.Flags().BoolVar(&showMatches, "show-matches", false, "Show all servers matching the pattern instead of details")
}
