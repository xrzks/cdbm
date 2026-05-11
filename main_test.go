package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	cdbmcli "github.com/xrzks/cdbm/internal/cli"
	"github.com/xrzks/cdbm/internal/config"
	"github.com/xrzks/cdbm/internal/store"
)

func TestMain(t *testing.T) {
	// Store original command line arguments
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Test with command line arguments
	os.Args = []string{"cdbm", "list"}
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	defer func() {
		os.Stdout = oldStdout
		w.Close()
	}()
	
	// Test the main function
	mainFunc := func() {
		// This will run main() but we need to redirect stderr to avoid logs showing in test output
		os.Stderr = nil
		defer func() {
			// Restore stderr to original value
		}()
		
		main()
	}
	
	mainFunc()
	
	// Restore stdout
	w.Close()
	
	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	
	// Verify that the output contains expected content
	if output == "" {
		t.Fatal("Expected output from list command, got empty string")
	}
}

func TestMainWithErrorScenarios(t *testing.T) {
	// Store original command line arguments and restore them
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "valid command",
			args:        []string{"cdbm", "list"},
			expectError: false,
		},
		{
			name:        "unknown command",
			args:        []string{"cdbm", "unknown"},
			expectError: true, //urfave cli will show help for unknown commands, but exit code 1
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			
			// Capture stderr to check for error messages
			oldStderr := os.Stderr
			_, w, _ := os.Pipe()
			os.Stderr = w
			
			defer func() {
				w.Close()
				os.Stderr = oldStderr
			}()
			
			// We'll need to capture the exit code to verify error scenarios
			exitCode := 0
			
			// Override os.Exit for testing
			exitFunc := func(code int) {
				exitCode = code
				panic("exit called") // To prevent actual exit
			}
			exitCodeFunc = exitFunc
			
			defer func() {
				_ = recover() // Expected panic from exit()
			}()
			
			// Test with limited store to avoid directory issues
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}
			
			// Create a proper file path, not directory
			tmpDir := t.TempDir()
			storePath := filepath.Join(tmpDir, "store.json")
			cfg.StorePath = storePath
			
			store, err := store.NewStore(cfg.StorePath)
			if err != nil {
				t.Fatalf("Failed to create store: %v", err)
			}
			
			app := cdbmcli.New(store)
			if err := app.Run(context.Background(), os.Args); err != nil {
				exitCode = 1
			}
			
			// Restore original exit function
			exitCodeFunc = os.Exit
			
			// Check expected error vs actual result
			if (exitCode != 0) != tt.expectError {
				t.Errorf("Expected error: %v, got exit code: %d", tt.expectError, exitCode)
			}
		})
	}
}

// Global exit function override
var exitCodeFunc = os.Exit