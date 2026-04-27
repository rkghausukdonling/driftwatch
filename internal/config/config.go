package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Provider   string            `yaml:"provider"`
	Region     string            `yaml:"region"`
	Statefile  string            `yaml:"statefile"`
	OutputFmt  string            `yaml:"output_format"`
	Ignore     []string          `yaml:"ignore"`
	Tags       map[string]string `yaml:"tags"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Provider:  "aws",
		Region:    "us-east-1",
		OutputFmt: "table",
	}
}

// Load reads a YAML config file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// Validate checks that required fields are present and valid.
func (c *Config) Validate() error {
	validProviders := map[string]bool{"aws": true, "gcp": true, "azure": true}
	if !validProviders[c.Provider] {
		return fmt.Errorf("unsupported provider %q (must be aws, gcp, or azure)", c.Provider)
	}

	validFormats := map[string]bool{"table": true, "json": true, "yaml": true}
	if !validFormats[c.OutputFmt] {
		return fmt.Errorf("unsupported output_format %q (must be table, json, or yaml)", c.OutputFmt)
	}

	return nil
}
