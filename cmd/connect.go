package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/jbovet/mcp-cli/pkg/adapter"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

var (
	// Flags for connect command
	connectType     string
	connectURL      string
	connectCommand  string
	connectArgs     []string
	connectEnv      []string
	connectTimeout  time.Duration
	interactiveMode bool
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an MCP server using different transport methods",
	Long: `Connect to an MCP (Model Context Protocol) server using various transport methods.

Supports:
- stdio: Connect to a local server process via standard input/output
- http: Connect to an HTTP-based MCP server
- streamable: Connect to a streamable HTTP-based MCP server

The command can run in interactive mode to explore the server's capabilities
or execute specific operations.`,
	Example: `  # Connect to a stdio server
  mcp-cli connect --type stdio --command "python" --args "server.py"

  # Connect to an HTTP server
  mcp-cli connect --type http --url "http://localhost:8080/mcp"

  # Connect with custom environment variables
  mcp-cli connect --type stdio --command "node" --args "server.js" --env "DEBUG=1"

  # Connect in interactive mode
  mcp-cli connect --type stdio --command "python server.py" --interactive`,
	RunE: runConnectCommand,
}

func runConnectCommand(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	// Create adapter configuration
	config := adapter.Config{
		ServerURL: connectURL,
		Command:   connectCommand,
		Args:      connectArgs,
		Env:       connectEnv,
		Timeout:   connectTimeout,
		Verbose:   verbose,
	}

	// Parse command string if provided as a single argument
	if connectCommand != "" && len(connectArgs) == 0 {
		parts := strings.Fields(connectCommand)
		if len(parts) > 1 {
			config.Command = parts[0]
			config.Args = parts[1:]
		}
	}

	// Create the appropriate adapter
	adapterType := adapter.AdapterType(connectType)
	serverAdapter, err := adapter.NewAdapter(adapterType, config)
	if err != nil {
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	// Connect to the server
	if verbose {
		fmt.Printf("Connecting to MCP server using %s transport...\n", connectType)
	}

	if err := serverAdapter.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer serverAdapter.Disconnect()

	// Get server information
	serverInfo, err := serverAdapter.GetServerInfo()
	if err != nil {
		return fmt.Errorf("failed to get server info: %w", err)
	}

	fmt.Printf("✓ Connected to MCP server: %s (version %s)\n\n",
		serverInfo.Name, serverInfo.Version)

	if interactiveMode {
		return runInteractiveMode(ctx, serverAdapter)
	}

	// Default: show server capabilities
	return showServerCapabilities(ctx, serverAdapter)
}

func showServerCapabilities(ctx context.Context, adapter adapter.ServerAdapter) error {
	fmt.Println("Server Capabilities:")
	fmt.Println("===================")

	// Show tools
	tools, err := adapter.ListTools(ctx)
	if err != nil {
		fmt.Printf("Failed to list tools: %v\n", err)
	} else {
		fmt.Printf("\nTools (%d available):\n", len(tools))
		if len(tools) > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tDESCRIPTION")
			fmt.Fprintln(w, "----\t-----------")
			for _, tool := range tools {
				description := tool.Description
				if len(description) > 60 {
					description = description[:57] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\n", tool.Name, description)
			}
			w.Flush()
		} else {
			fmt.Println("  No tools available")
		}
	}

	// Show resources
	resources, err := adapter.ListResources(ctx)
	if err != nil {
		fmt.Printf("Failed to list resources: %v\n", err)
	} else {
		fmt.Printf("\nResources (%d available):\n", len(resources))
		if len(resources) > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "URI\tNAME\tDESCRIPTION")
			fmt.Fprintln(w, "---\t----\t-----------")
			for _, resource := range resources {
				description := resource.Description
				if len(description) > 50 {
					description = description[:47] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", resource.URI, resource.Name, description)
			}
			w.Flush()
		} else {
			fmt.Println("  No resources available")
		}
	}

	// Show prompts
	prompts, err := adapter.ListPrompts(ctx)
	if err != nil {
		fmt.Printf("Failed to list prompts: %v\n", err)
	} else {
		fmt.Printf("\nPrompts (%d available):\n", len(prompts))
		if len(prompts) > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tDESCRIPTION")
			fmt.Fprintln(w, "----\t-----------")
			for _, prompt := range prompts {
				description := prompt.Description
				if len(description) > 60 {
					description = description[:57] + "..."
				}
				fmt.Fprintf(w, "%s\t%s\n", prompt.Name, description)
			}
			w.Flush()
		} else {
			fmt.Println("  No prompts available")
		}
	}

	return nil
}

func runInteractiveMode(ctx context.Context, adapter adapter.ServerAdapter) error {
	fmt.Println("Interactive Mode - Type 'help' for available commands")
	fmt.Println("====================================================")

	// Use bufio.Scanner for better input handling
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "help":
			showInteractiveHelp()
		case "tools":
			listToolsInteractive(ctx, adapter)
		case "resources":
			listResourcesInteractive(ctx, adapter)
		case "prompts":
			listPromptsInteractive(ctx, adapter)
		case "call":
			if len(parts) < 2 {
				fmt.Println("Usage: call <tool-name> [arguments...]")
				continue
			}
			callToolInteractive(ctx, adapter, parts[1], parts[2:])
		case "read":
			if len(parts) < 2 {
				fmt.Println("Usage: read <resource-uri>")
				continue
			}
			readResourceInteractive(ctx, adapter, parts[1])
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return nil
		default:
			fmt.Printf("Unknown command: %s (type 'help' for available commands)\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}

func showInteractiveHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  help                                    - Show this help message")
	fmt.Println("  tools                                   - List available tools")
	fmt.Println("  resources                               - List available resources")
	fmt.Println("  prompts                                 - List available prompts")
	fmt.Println("  call <tool-name> [arguments...]         - Call a tool with arguments")
	fmt.Println("  read <uri>                              - Read a resource")
	fmt.Println("  quit, exit                              - Exit interactive mode")
	fmt.Println()
	fmt.Println("Examples for API Linter tools:")
	fmt.Println("  call list_rules                         - List all API Linter rules")
	fmt.Println("  call get_rule_categories                - Get rule categories")
	fmt.Println("  call get_api_docs                       - Get API design guide")
	fmt.Println("  call validate_api_specification <spec>  - Validate OpenAPI spec")
	fmt.Println("  call calculate_api_quality_score <spec> - Calculate quality score")
	fmt.Println()
}

func listToolsInteractive(ctx context.Context, adapter adapter.ServerAdapter) {
	tools, err := adapter.ListTools(ctx)
	if err != nil {
		fmt.Printf("Error listing tools: %v\n", err)
		return
	}

	if len(tools) == 0 {
		fmt.Println("No tools available")
		return
	}

	fmt.Printf("Available tools (%d):\n", len(tools))
	for i, tool := range tools {
		fmt.Printf("%d. %s - %s\n", i+1, tool.Name, tool.Description)
	}
}

func listResourcesInteractive(ctx context.Context, adapter adapter.ServerAdapter) {
	resources, err := adapter.ListResources(ctx)
	if err != nil {
		fmt.Printf("Error listing resources: %v\n", err)
		return
	}

	if len(resources) == 0 {
		fmt.Println("No resources available")
		return
	}

	fmt.Printf("Available resources (%d):\n", len(resources))
	for i, resource := range resources {
		fmt.Printf("%d. %s (%s) - %s\n", i+1, resource.URI, resource.Name, resource.Description)
	}
}

func listPromptsInteractive(ctx context.Context, adapter adapter.ServerAdapter) {
	prompts, err := adapter.ListPrompts(ctx)
	if err != nil {
		fmt.Printf("Error listing prompts: %v\n", err)
		return
	}

	if len(prompts) == 0 {
		fmt.Println("No prompts available")
		return
	}

	fmt.Printf("Available prompts (%d):\n", len(prompts))
	for i, prompt := range prompts {
		fmt.Printf("%d. %s - %s\n", i+1, prompt.Name, prompt.Description)
	}
}

func callToolInteractive(ctx context.Context, adapter adapter.ServerAdapter, toolName string, args []string) {
	// Handle different argument patterns based on the tool
	arguments := make(map[string]any)

	switch toolName {
	case "list_rules", "get_rule_categories", "get_api_docs", "get_api_linter_link":
		// These tools typically don't need arguments
		// Keep arguments empty

	case "validate_api_specification":
		if len(args) > 0 {
			arguments["specification"] = args[0]
		} else {
			fmt.Println("Usage: call validate_api_specification <specification_content_or_url>")
			return
		}
		if len(args) > 1 {
			arguments["format"] = args[1] // e.g., "openapi", "swagger"
		}

	case "calculate_api_quality_score":
		if len(args) > 0 {
			arguments["specification"] = args[0]
		} else {
			fmt.Println("Usage: call calculate_api_quality_score <specification_content_or_url>")
			return
		}

	default:
		// For other tools, use generic argument parsing
		if len(args) == 1 && (strings.HasPrefix(args[0], "{") || strings.Contains(args[0], ":")) {
			// Try to parse as JSON-like input
			arguments["input"] = args[0]
		} else {
			// Simple key=value or positional arguments
			for i, arg := range args {
				if strings.Contains(arg, "=") {
					parts := strings.SplitN(arg, "=", 2)
					arguments[parts[0]] = parts[1]
				} else {
					arguments[fmt.Sprintf("arg%d", i)] = arg
				}
			}
		}
	}

	if verbose {
		fmt.Printf("Calling tool '%s' with arguments: %+v\n", toolName, arguments)
	}

	result, err := adapter.CallTool(ctx, toolName, arguments)
	if err != nil {
		fmt.Printf("Error calling tool %s: %v\n", toolName, err)
		return
	}

	fmt.Printf("Tool '%s' result:\n", toolName)
	if result.IsError {
		fmt.Println("❌ Tool execution resulted in an error:")
	} else {
		fmt.Println("✅ Tool executed successfully:")
	}

	for i, content := range result.Content {
		switch c := content.(type) {
		case mcp.TextContent:
			if len(result.Content) > 1 {
				fmt.Printf("  Content %d (Text):\n", i+1)
			}
			// Format the text content nicely
			text := strings.TrimSpace(c.Text)
			if strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
				// Try to format JSON
				fmt.Printf("%s\n", text)
			} else {
				// Regular text, add some formatting
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						fmt.Printf("    %s\n", line)
					}
				}
			}
		case mcp.ImageContent:
			fmt.Printf("  Image (%s): [base64 data - %d bytes]\n", c.MIMEType, len(c.Data))
		case mcp.AudioContent:
			fmt.Printf("  Audio (%s): [base64 data - %d bytes]\n", c.MIMEType, len(c.Data))
		default:
			fmt.Printf("  Content %d: %+v\n", i+1, c)
		}
	}
	fmt.Println()
}

func readResourceInteractive(ctx context.Context, adapter adapter.ServerAdapter, uri string) {
	result, err := adapter.ReadResource(ctx, uri)
	if err != nil {
		fmt.Printf("Error reading resource %s: %v\n", uri, err)
		return
	}

	fmt.Printf("Resource content:\n")
	for _, content := range result.Contents {
		switch c := content.(type) {
		case mcp.TextResourceContents:
			fmt.Printf("  Text (%s): %s\n", c.MIMEType, c.Text)
		case mcp.BlobResourceContents:
			fmt.Printf("  Blob (%s): %s\n", c.MIMEType, "base64 data")
		default:
			fmt.Printf("  Content: %+v\n", c)
		}
	}
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringVar(&connectType, "type", "stdio", "Transport type (stdio, http, streamable)")
	connectCmd.Flags().StringVar(&connectURL, "url", "", "Server URL for HTTP-based connections")
	connectCmd.Flags().StringVar(&connectCommand, "command", "", "Command to execute for stdio connections")
	connectCmd.Flags().StringArrayVar(&connectArgs, "args", nil, "Arguments for the command")
	connectCmd.Flags().StringArrayVar(&connectEnv, "env", nil, "Environment variables for the command")
	connectCmd.Flags().DurationVar(&connectTimeout, "timeout", 60*time.Second, "Connection timeout")
	connectCmd.Flags().BoolVar(&interactiveMode, "interactive", false, "Run in interactive mode")
}
