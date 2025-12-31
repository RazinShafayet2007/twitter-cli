package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	CurrentUser string `json:"current_user,omitempty"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.twitter-cli-config.json"
	}
	return filepath.Join(home, ".twitter-cli", "config.json")
}

// LoadConfig loads the config from disk
func LoadConfig() (*Config, error) {
	path := GetConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, return empty config
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the config to disk
func SaveConfig(config *Config) error {
	path := GetConfigPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetCurrentUser returns the currently logged-in user
func GetCurrentUser() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	if config.CurrentUser == "" {
		return "", errors.New("not logged in")
	}

	return config.CurrentUser, nil
}

// SetCurrentUser sets the currently logged-in user
func SetCurrentUser(username string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.CurrentUser = username
	return SaveConfig(config)
}

// ClearCurrentUser logs out the current user
func ClearCurrentUser() error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.CurrentUser = ""
	return SaveConfig(config)
}
