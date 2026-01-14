package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigFile = ".cdbm.json"
	defaultStoreFile  = "store.json"
)

type Config struct {
	StorePath string `json:"store_path"`
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	cfg := &Config{
		StorePath: filepath.Join(configPath, defaultStoreFile),
	}

	data, err := os.ReadFile(filepath.Join(configPath, defaultConfigFile))
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.StorePath = expandPath(cfg.StorePath)

	if cfg.StorePath == "" {
		return nil, fmt.Errorf("config field 'store_path' is empty or missing. Please set a valid store path in %s", filepath.Join(configPath, defaultConfigFile))
	}

	return cfg, nil
}

func getConfigPath() (string, error) {
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		return filepath.Join(configDir, "cdbm"), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "cdbm"), nil
}

func expandPath(path string) string {
	expanded := os.ExpandEnv(path)

	if len(expanded) > 0 && expanded[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			if len(expanded) == 1 || expanded[1] == '/' {
				expanded = filepath.Join(homeDir, expanded[1:])
			}
		}
	}

	return expanded
}
