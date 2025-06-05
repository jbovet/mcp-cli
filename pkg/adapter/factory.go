package adapter

import (
	"fmt"
	"strings"
	"time"
)

// AdapterFactory provides methods to create server adapters
type AdapterFactory struct{}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{}
}

// CreateFromConfig creates an adapter based on configuration
func (f *AdapterFactory) CreateFromConfig(config map[string]interface{}) (ServerAdapter, error) {
	adapterType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("adapter type is required")
	}

	adapterConfig := Config{
		Verbose: getBool(config, "verbose", false),
		Timeout: getDuration(config, "timeout", 30*time.Second),
	}

	switch AdapterType(adapterType) {
	case AdapterTypeStdio:
		command, ok := config["command"].(string)
		if !ok {
			return nil, fmt.Errorf("command is required for stdio adapter")
		}
		adapterConfig.Command = command

		if args, ok := config["args"].([]string); ok {
			adapterConfig.Args = args
		}

		if env, ok := config["env"].([]string); ok {
			adapterConfig.Env = env
		}

		return NewStdioAdapter(adapterConfig)

	case AdapterTypeHTTP, AdapterTypeStreamable:
		url, ok := config["url"].(string)
		if !ok {
			return nil, fmt.Errorf("URL is required for HTTP adapter")
		}
		adapterConfig.ServerURL = url

		return NewHTTPAdapter(adapterConfig)

	default:
		return nil, fmt.Errorf("unsupported adapter type: %s", adapterType)
	}
}

// CreateFromURL creates an adapter from a URL string
func (f *AdapterFactory) CreateFromURL(url string, verbose bool) (ServerAdapter, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		config := Config{
			ServerURL: url,
			Verbose:   verbose,
			Timeout:   30 * time.Second,
		}
		return NewHTTPAdapter(config)
	}

	// Assume it's a command for stdio
	parts := strings.Fields(url)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command: %s", url)
	}

	config := Config{
		Command: parts[0],
		Args:    parts[1:],
		Verbose: verbose,
		Timeout: 30 * time.Second,
	}

	return NewStdioAdapter(config)
}

// ValidateConfig validates adapter configuration
func (f *AdapterFactory) ValidateConfig(adapterType AdapterType, config Config) error {
	switch adapterType {
	case AdapterTypeStdio:
		if config.Command == "" {
			return fmt.Errorf("command is required for stdio adapter")
		}
	case AdapterTypeHTTP, AdapterTypeStreamable:
		if config.ServerURL == "" {
			return fmt.Errorf("server URL is required for HTTP adapter")
		}
		if !strings.HasPrefix(config.ServerURL, "http://") &&
			!strings.HasPrefix(config.ServerURL, "https://") {
			return fmt.Errorf("server URL must start with http:// or https://")
		}
	default:
		return fmt.Errorf("unsupported adapter type: %s", adapterType)
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	return nil
}

// GetSupportedTypes returns the list of supported adapter types
func (f *AdapterFactory) GetSupportedTypes() []AdapterType {
	return []AdapterType{
		AdapterTypeStdio,
		AdapterTypeHTTP,
		AdapterTypeStreamable,
	}
}

// Helper functions
func getBool(config map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := config[key].(bool); ok {
		return val
	}
	return defaultValue
}

func getString(config map[string]interface{}, key string, defaultValue string) string {
	if val, ok := config[key].(string); ok {
		return val
	}
	return defaultValue
}

func getStringSlice(config map[string]interface{}, key string) []string {
	if val, ok := config[key].([]string); ok {
		return val
	}
	if val, ok := config[key].([]interface{}); ok {
		result := make([]string, len(val))
		for i, v := range val {
			if s, ok := v.(string); ok {
				result[i] = s
			}
		}
		return result
	}
	return nil
}

func getDuration(config map[string]interface{}, key string, defaultValue time.Duration) time.Duration {
	if val, ok := config[key].(time.Duration); ok {
		return val
	}
	if val, ok := config[key].(string); ok {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}
