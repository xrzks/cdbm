package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.StorePath == "" {
		t.Error("Load() should set default store path")
	}
}

func TestLoadWithExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	configDir := filepath.Join(tmpDir, "cdbm")
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, ".cdbm.json")
	configContent := `{"store_path": "/custom/store/path.json"}`
	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := "/custom/store/path.json"
	if cfg.StorePath != expectedPath {
		t.Errorf("Load() store_path = %v, want %v", cfg.StorePath, expectedPath)
	}
}

func TestLoadWithEnvVar(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	configDir := filepath.Join(tmpDir, "cdbm")
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, ".cdbm.json")
	configContent := `{"store_path": "$CUSTOM_DIR/store.json"}`
	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	customDir := filepath.Join(tmpDir, "custom")
	os.Setenv("CUSTOM_DIR", customDir)
	defer os.Unsetenv("CUSTOM_DIR")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(customDir, "store.json")
	if cfg.StorePath != expectedPath {
		t.Errorf("Load() store_path = %v, want %v", cfg.StorePath, expectedPath)
	}
}

func TestLoadWithTilde(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	oldHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHomeDir)

	configDir := filepath.Join(tmpDir, "cdbm")
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, ".cdbm.json")
	configContent := `{"store_path": "~/custom/store.json"}`
	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "custom", "store.json")
	if cfg.StorePath != expectedPath {
		t.Errorf("Load() store_path = %v, want %v", cfg.StorePath, expectedPath)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	configDir := filepath.Join(tmpDir, "cdbm")
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, ".cdbm.json")
	configContent := `invalid json`
	err = os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Error("Load() should return error for invalid JSON")
	}
}

func TestGetConfigPath(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "cdbm")
	if path != expectedPath {
		t.Errorf("getConfigPath() = %v, want %v", path, expectedPath)
	}
}

func TestGetConfigPathWithHome(t *testing.T) {
	tmpDir := t.TempDir()

	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	oldHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHomeDir)

	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}

	expectedPath := filepath.Join(tmpDir, ".config", "cdbm")
	if path != expectedPath {
		t.Errorf("getConfigPath() = %v, want %v", path, expectedPath)
	}
}

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		homeDir  string
		expected string
	}{
		{
			name:     "simple path",
			input:    "/simple/path",
			expected: "/simple/path",
		},
		{
			name:     "env var",
			input:    "$TEST_VAR/path",
			envVars:  map[string]string{"TEST_VAR": "/test"},
			expected: "/test/path",
		},
		{
			name:     "env var with braces",
			input:    "${TEST_VAR}/path",
			envVars:  map[string]string{"TEST_VAR": "/test"},
			expected: "/test/path",
		},
		{
			name:     "tilde alone",
			input:    "~",
			homeDir:  "/home/user",
			expected: "/home/user",
		},
		{
			name:     "tilde with path",
			input:    "~/path",
			homeDir:  "/home/user",
			expected: "/home/user/path",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no expansion",
			input:    "no/special/chars",
			expected: "no/special/chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			oldHomeDir := os.Getenv("HOME")
			if tt.homeDir != "" {
				os.Setenv("HOME", tt.homeDir)
			}
			defer os.Setenv("HOME", oldHomeDir)

			result := expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
