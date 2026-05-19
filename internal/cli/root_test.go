package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/logger"
	"github.com/xrzks/cdbm/internal/store"
)

func TestCLI_New(t *testing.T) {
	// Create a temporary store path
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	cmd := New(s)
	
	// Verify command structure
	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}
	
	if cmd.Name != "cdbm" {
		t.Errorf("Expected command name 'cdbm', got '%s'", cmd.Name)
	}
	
	if cmd.Usage != "Manage directory bookmarks" {
		t.Errorf("Expected usage 'Manage directory bookmarks', got '%s'", cmd.Usage)
	}
	
	// Check that commands are registered
	expectedCommands := []string{"add", "list", "init", "edit", "remove"}
	cmdMap := make(map[string]bool)
	for _, subCmd := range cmd.Commands {
		cmdMap[subCmd.Name] = true
	}
	
	for _, expectedCmd := range expectedCommands {
		if !cmdMap[expectedCmd] {
			t.Errorf("Expected command '%s' not found", expectedCmd)
		}
	}
}

func TestCLI_setupLogger(t *testing.T) {
	// Create a temporary store path
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	c := &CLI{store: s}
	
	// Test without debug flag
	ctx := context.Background()
	cmd := &cli.Command{}
	// Use reflection or direct function call to check flags in setupLogger
	
	newCtx, err := c.setupLogger(ctx, cmd)
	if err != nil {
		t.Errorf("Unexpected error in setupLogger without debug: %v", err)
	}
	
	if newCtx != ctx {
		t.Error("Expected same context returned when debug is false")
	}
	
	// Test with debug flag - skip this test as it requires mocking config.GetStatePath
	t.Skip("Skipping debug flag test as it requires mocking config.GetStatePath")
}

func TestCLI_closeLogger(t *testing.T) {
	// Create a temporary store path
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create CLI
	c := &CLI{store: s}
	
	// Test closing nil logger
	err = c.closeLogger(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("Unexpected error closing nil logger: %v", err)
	}
	
	// Test with a real logger
	logPath := filepath.Join(tmpDir, "test.log")
	fileLogger, err := logger.NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	
	c.logger = fileLogger
	err = c.closeLogger(context.Background(), &cli.Command{})
	if err != nil {
		t.Errorf("Unexpected error closing real logger: %v", err)
	}
}

func TestCLI_logDebug(t *testing.T) {
	// Create a temporary store path
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Test logging without logger (should not crash)
	c := &CLI{store: s}
	c.logDebug("test", map[string]any{"key": "value"})
	
	// Test with logger
	logPath := filepath.Join(tmpDir, "test.log")
	fileLogger, err := logger.NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	
	c.logger = fileLogger
	
	// This should not crash
	c.logDebug("test", map[string]any{"key": "value"})
	
	// Close logger
	fileLogger.Close()
}

func TestCLI_logDebugWithError(t *testing.T) {
	// Create a temporary store path
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "store.json")
	
	// Create a new store
	s, err := store.NewStore(storePath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create a mock logger that always returns an error
	mockLogger := &mockErrorLogger{}
	c := &CLI{store: s, logger: mockLogger}
	
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	
	defer func() {
		w.Close()
		os.Stderr = oldStderr
	}()
	
	// This should not crash even though the logger returns an error
	c.logDebug("test", map[string]any{"key": "value"})
	
	// Restore stderr
	w.Close()
	
	// Read captured stderr
	output := make([]byte, 1024)
	n, err := r.Read(output)
	if err != nil && err != io.EOF {
		t.Errorf("Error reading from pipe: %v", err)
	}
	strOutput := string(output[:n])
	
	// Verify that the warning message was written to stderr
	if !strings.Contains(strOutput, "Warning: failed to write debug log entry") {
		t.Errorf("Expected warning message in stderr, got: %s", strOutput)
	}
}

// Mock logger that always returns an error
type mockErrorLogger struct{}

func (m *mockErrorLogger) Log(action string, details map[string]any) error {
	return errors.New("mock logger error")
}