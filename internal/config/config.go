package config

import (
	"os"

	"gopkg.in/yaml.v2"

	"github.com/yourusername/codemap/internal/types"
)

// LoadConfig reads and parses the .codemap YAML configuration file
// It returns a Config struct or an error if the file cannot be read or parsed
func LoadConfig(path string) (*types.Config, error) {
	// TODO: Implement YAML file reading and parsing
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config types.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}