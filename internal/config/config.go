// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package config provides configuration loading for the poweradmin CLI.
// Configuration is read from a YAML file and can be overridden by
// environment variables or CLI flags.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the CLI configuration loaded from the config file.
// Fields map directly to YAML keys in the config file.
type Config struct {
	// URL is the base URL of the Poweradmin instance (e.g. https://dns.example.com).
	URL string `yaml:"url"`
	// APIKey is the API key used to authenticate against the Poweradmin API.
	APIKey string `yaml:"api_key"`
}

// DefaultPath returns the platform-appropriate default path for the config file.
// On Linux and macOS this resolves to ~/.config/poweradmin/config.yaml.
func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "poweradmin", "config.yaml")
}

// Load reads and parses the YAML config file at the given path.
// If the file does not exist, an empty Config is returned without error —
// missing config is not treated as a failure since credentials can be
// supplied via environment variables or CLI flags instead.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
