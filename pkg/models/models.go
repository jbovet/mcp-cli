package models

// Server represents basic server information
type Server struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Repository    Repository    `json:"repository"`
	VersionDetail VersionDetail `json:"version_detail"`
}

// ServerDetail represents detailed server information
type ServerDetail struct {
	Server   `json:",inline"`
	Packages []Package `json:"packages,omitempty"`
	Remotes  []Remote  `json:"remotes,omitempty"`
}

// Repository represents a source code repository
type Repository struct {
	URL    string `json:"url"`
	Source string `json:"source"`
	ID     string `json:"id"`
}

// VersionDetail represents the version details of a server
type VersionDetail struct {
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
	IsLatest    bool   `json:"is_latest"`
}

// Package represents a package configuration
type Package struct {
	RegistryName         string          `json:"registry_name"`
	Name                 string          `json:"name"`
	Version              string          `json:"version"`
	RuntimeHint          string          `json:"runtime_hint,omitempty"`
	RuntimeArguments     []Argument      `json:"runtime_arguments,omitempty"`
	PackageArguments     []Argument      `json:"package_arguments,omitempty"`
	EnvironmentVariables []KeyValueInput `json:"environment_variables,omitempty"`
}

// Remote represents a remote connection endpoint
type Remote struct {
	TransportType string  `json:"transport_type"`
	URL           string  `json:"url"`
	Headers       []Input `json:"headers,omitempty"`
}

// Format represents input format types
type Format string

const (
	FormatString   Format = "string"
	FormatNumber   Format = "number"
	FormatBoolean  Format = "boolean"
	FormatFilePath Format = "file_path"
)

// ArgumentType represents argument types
type ArgumentType string

const (
	ArgumentTypePositional ArgumentType = "positional"
	ArgumentTypeNamed      ArgumentType = "named"
)

// Input represents a user input configuration
type Input struct {
	Description string           `json:"description,omitempty"`
	IsRequired  bool             `json:"is_required,omitempty"`
	Format      Format           `json:"format,omitempty"`
	Value       string           `json:"value,omitempty"`
	IsSecret    bool             `json:"is_secret,omitempty"`
	Default     string           `json:"default,omitempty"`
	Choices     []string         `json:"choices,omitempty"`
	Template    string           `json:"template,omitempty"`
	Properties  map[string]Input `json:"properties,omitempty"`
}

// InputWithVariables extends Input with variables
type InputWithVariables struct {
	Input     `json:",inline"`
	Variables map[string]Input `json:"variables,omitempty"`
}

// KeyValueInput represents a key-value input pair
type KeyValueInput struct {
	InputWithVariables `json:",inline"`
	Name               string `json:"name"`
}

// Argument represents a runtime or package argument
type Argument struct {
	InputWithVariables `json:",inline"`
	Type               ArgumentType `json:"type"`
	Name               string       `json:"name,omitempty"`
	IsRepeated         bool         `json:"is_repeated,omitempty"`
	ValueHint          string       `json:"value_hint,omitempty"`
}
