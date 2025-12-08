package types

// Config represents the .codemap configuration file structure
type Config struct {
	Version string   `yaml:"version"`
	Sections []Section `yaml:"sections"`
}

// Section defines a named section of the codebase for mapping
type Section struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// Definition represents a parsed code definition (function, type, etc.)
type Definition struct {
	Type      string // "function", "type", etc.
	Name      string
	Line       int
	Signature string
	Comment   string
}

// FileMap contains the parsed definitions for a single file
type FileMap struct {
	Path        string
	Language    string
	Definitions []Definition
}