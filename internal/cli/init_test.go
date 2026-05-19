package cli

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestInstallShellIntegration(t *testing.T) {
	tests := []struct {
		shell     string
		expectErr bool
	}{
		{"zsh", false},
		{"bash", false},
		{"fish", false},
		{"", true},
		{"unsupported", true},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			oldStdout := stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := installShellIntegration(tt.shell)

			w.Close()
			os.Stdout = oldStdout

			output := make([]byte, 1024)
			n, _ := r.Read(output)
			strOutput := string(output[:n])

			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got error: %v", tt.expectErr, err)
			}

			if !tt.expectErr {
				if len(strOutput) == 0 {
					t.Error("Expected output for supported shell, but got empty string")
				}
			} else {
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

			_ = strOutput
			_ = io.EOF
		})
	}
}

var stdout = os.Stdout

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
