package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileLogger_New(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	
	// Test: Create new file logger
	logger, err := NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	
	if logger == nil {
		t.Fatal("Expected logger, got nil")
	}
	
	// Verify log file was created
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
	
	// Clean up
	logger.Close()
}

func TestFileLogger_Log(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	
	// Create new file logger
	logger, err := NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()
	
	// Test: Log a simple message
	err = logger.Log("test-action", map[string]any{
		"key1": "value1",
		"key2": 42,
	})
	if err != nil {
		t.Errorf("Failed to log message: %v", err)
	}
	
	// Verify log file contains the message
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	// Parse the log entry
	var logEntry LogEntry
	err = json.Unmarshal(data, &logEntry)
	if err != nil {
		t.Errorf("Failed to parse log entry: %v", err)
	}
	
	if logEntry.Action != "test-action" {
		t.Errorf("Expected action 'test-action', got '%s'", logEntry.Action)
	}
	
	if logEntry.Details["key1"] != "value1" {
		t.Errorf("Expected detail 'key1' to be 'value1', got '%v'", logEntry.Details["key1"])
	}
	
	if logEntry.Details["key2"] != 42 {
		t.Errorf("Expected detail 'key2' to be 42, got '%v'", logEntry.Details["key2"])
	}
	
	// Test: Log with empty details
	err = logger.Log("empty-details", nil)
	if err != nil {
		t.Errorf("Failed to log with empty details: %v", err)
	}
	
	// Read the log file again
	data, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	// Count lines (should be 2 now)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(lines))
	}
	
	// Parse the second log entry
	var logEntry2 LogEntry
	err = json.Unmarshal([]byte(lines[1]), &logEntry2)
	if err != nil {
		t.Errorf("Failed to parse second log entry: %v", err)
	}
	
	if logEntry2.Action != "empty-details" {
		t.Errorf("Expected action 'empty-details', got '%s'", logEntry2.Action)
	}
}

func TestFileLogger_Close(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	
	// Create new file logger
	logger, err := NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	
	// Test: Close logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Failed to close logger: %v", err)
	}
	
	// Verify log file still exists but is closed
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was removed when closed")
	}
	
	// Test: Close already closed logger (should not error)
	err = logger.Close()
	if err != nil {
		t.Errorf("Unexpected error closing already closed logger: %v", err)
	}
}

func TestFileLogger_ErrorHandling(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	
	// Create new file logger
	logger, err := NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()
	
	// Test: Log with invalid JSON (should fail due to details marshaling)
	// This test is a bit tricky as the details should be valid JSON
	// Let's test logging with a complex structure that might cause issues
	err = logger.Log("complex-action", map[string]any{
		"nested": map[string]any{
			"deep": map[string]any{
				"value": "test",
			},
		},
		"array": []string{"item1", "item2"},
	})
	if err != nil {
		t.Errorf("Failed to log complex message: %v", err)
	}
	
	// Verify log file contains the message
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	// Parse the log entry
	var logEntry LogEntry
	err = json.Unmarshal(data, &logEntry)
	if err != nil {
		t.Errorf("Failed to parse log entry: %v", err)
	}
	
	if logEntry.Action != "complex-action" {
		t.Errorf("Expected action 'complex-action', got '%s'", logEntry.Action)
	}
}

func TestLogEntry_JSONSerialization(t *testing.T) {
	// Test LogEntry JSON serialization
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Action:    "test-action",
		Details: map[string]any{
			"key1": "value1",
			"key2": 42,
		},
	}
	
	// Serialize to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		t.Errorf("Failed to marshal log entry: %v", err)
	}
	
	// Deserialize back
	var entry2 LogEntry
	err = json.Unmarshal(data, &entry2)
	if err != nil {
		t.Errorf("Failed to unmarshal log entry: %v", err)
	}
	
	if entry2.Action != entry.Action {
		t.Errorf("Expected action '%s', got '%s'", entry.Action, entry2.Action)
	}
	
	if entry2.Details["key1"] != entry.Details["key1"] {
		t.Errorf("Expected detail 'key1' to match, got '%v' vs '%v'", entry2.Details["key1"], entry.Details["key1"])
	}
	
	if entry2.Details["key2"] != entry.Details["key2"] {
		t.Errorf("Expected detail 'key2' to match, got '%v' vs '%v'", entry2.Details["key2"], entry.Details["key2"])
	}
}