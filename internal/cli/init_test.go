package cli

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestInstallShellIntegration(t *testing.T) {
	// Test cases for supported shells
	tests := []struct {
		shell     string
		expectErr bool
	}{
		{"zsh", false},
		{"bash", false},
		{"fish", false},
		{"", true}, // No shell specified
		{"unsupported", true}, // Unsupported shell
	}
	
	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			// Capture stdout
			oldStdout := stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			
			defer func() {
				os.Stdout = oldStdout
				w.Close()
			}()
			
			err := installShellIntegration(tt.shell)
			
			// Restore stdout
			w.Close()
			
			// Read captured stdout
			output := make([]byte, 1024)
			n, err := r.Read(output)
			if err != nil && err != io.EOF {
				t.Errorf("Error reading from pipe: %v", err)
			}
			strOutput := string(output[:n])
			
			// Check error expectations
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got error: %v", tt.expectErr, err)
			}
			
			// Check output expectations
			if !tt.expectErr {
				if len(strOutput) == 0 {
					t.Error("Expected output for supported shell, but got empty string")
				}
			} else {
				// For error cases, verify error message contains expected info
				if tt.shell == "" {
					if !strings.Contains(err.Error(), "no shell specified") {
						t.Errorf("Expected 'no shell specified' error message, got: %v", err)
					}
				} else {
					if !strings.Contains(err.Error(), "unsupported shell") {
						t.Errorf("Expected 'unsupported shell' error message, got: %v", err)
					}
				}
			}
		})
	}
}

// Override stdout for testing
var stdout = os.Stdout

// Mock command for testing init command
type mockCommand struct {
	args []string
}

func (m *mockCommand) GetArg(index int) string {
	if index < len(m.args) {
		return m.args[index]
	}
	return ""
}

func (m *mockCommand) Args() []string {
	return m.args
}