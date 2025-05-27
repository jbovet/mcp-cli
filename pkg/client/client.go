package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jbovet/mcp-cli/pkg/models"
)

// Client represents the MCP Registry API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ServersResponse represents the response from the servers endpoint
type ServersResponse struct {
	Servers  []models.Server `json:"servers"`
	Metadata Metadata        `json:"metadata,omitempty"`
}

// Metadata contains pagination metadata
type Metadata struct {
	NextCursor string `json:"next_cursor,omitempty"`
	Count      int    `json:"count,omitempty"`
	Total      int    `json:"total,omitempty"`
}

// HealthResponse represents the response from the health endpoint
type HealthResponse struct {
	Status         string `json:"status"`
	GitHubClientID string `json:"github_client_id"`
}

// PingResponse represents the response from the ping endpoint
type PingResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ServerDetail represents detailed server information
type ServerDetail = models.ServerDetail

// GetServers fetches a list of servers from the API
func (c *Client) GetServers(cursor string, limit int) (*ServersResponse, error) {
	// Build query parameters
	params := url.Values{}
	if cursor != "" {
		params.Add("cursor", cursor)
	}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}

	// Build URL
	endpoint := fmt.Sprintf("%s/v0/servers", c.baseURL)
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	// Make HTTP request
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response ServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetServer fetches detailed information about a specific server by ID
func (c *Client) GetServer(id string) (*ServerDetail, error) {
	// Build URL
	endpoint := fmt.Sprintf("%s/v0/servers/%s", c.baseURL, id)

	// Make HTTP request
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("server with ID '%s' not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var server ServerDetail
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &server, nil
}

// GetServerByName fetches detailed information about a server by name
// This method searches through all servers to find a match by name
func (c *Client) GetServerByName(name string) (*ServerDetail, error) {
	var cursor string

	// Search through all pages to find the server
	for {
		response, err := c.GetServers(cursor, 100) // Use max limit for efficiency
		if err != nil {
			return nil, fmt.Errorf("failed to search servers: %w", err)
		}

		// Check each server in this page
		for _, server := range response.Servers {
			if server.Name == name {
				// Found the server, now get full details
				return c.GetServer(server.ID)
			}
		}

		// If no more pages, stop searching
		if response.Metadata.NextCursor == "" {
			break
		}
		cursor = response.Metadata.NextCursor
	}

	return nil, fmt.Errorf("server with name '%s' not found", name)
}

// FindServersByNamePattern finds servers that match a name pattern (case-insensitive substring match)
func (c *Client) FindServersByNamePattern(pattern string) ([]models.Server, error) {
	var matchingServers []models.Server
	var cursor string
	pattern = strings.ToLower(pattern)

	// Search through all pages
	for {
		response, err := c.GetServers(cursor, 100)
		if err != nil {
			return nil, fmt.Errorf("failed to search servers: %w", err)
		}

		// Check each server in this page
		for _, server := range response.Servers {
			if strings.Contains(strings.ToLower(server.Name), pattern) {
				matchingServers = append(matchingServers, server)
			}
		}

		// If no more pages, stop searching
		if response.Metadata.NextCursor == "" {
			break
		}
		cursor = response.Metadata.NextCursor
	}

	return matchingServers, nil
}

// GetHealth performs a health check against the service
func (c *Client) GetHealth() (*HealthResponse, error) {
	// Build URL
	endpoint := fmt.Sprintf("%s/v0/health", c.baseURL)

	// Make HTTP request
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &health, nil
}

// Ping sends a ping request to the service
func (c *Client) Ping() (*PingResponse, error) {
	// Build URL
	endpoint := fmt.Sprintf("%s/v0/ping", c.baseURL)

	// Make HTTP request
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var ping PingResponse
	if err := json.NewDecoder(resp.Body).Decode(&ping); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ping, nil
}
