package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config holds the configuration settings
type Config struct {
	settings map[string]string
}

// LoadConfig loads the configuration from a file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &Config{settings: make(map[string]string)}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by "=" to get key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid configuration line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config.settings[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}

// Get retrieves a configuration value by key, returning the value and a boolean indicating if it exists
func (c *Config) Get(key string) (string, bool) {
	value, exists := c.settings[key]
	return value, exists
}

// GetOrDefault retrieves a configuration value by key or returns a default value if the key is not found
func (c *Config) GetOrDefault(key, defaultValue string) string {
	if value, exists := c.settings[key]; exists {
		return value
	}
	return defaultValue
}
