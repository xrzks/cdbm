package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateBookmarkName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid simple name",
			input:   "test",
			wantErr: false,
		},
		{
			name:    "valid name with dots",
			input:   "test.name",
			wantErr: false,
		},
		{
			name:    "valid name with underscores",
			input:   "test_name",
			wantErr: false,
		},
		{
			name:    "valid name with hyphens",
			input:   "test-name",
			wantErr: false,
		},
		{
			name:    "valid name with numbers",
			input:   "test123",
			wantErr: false,
		},
		{
			name:    "valid mixed characters",
			input:   "Test.Name_123-valid",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "name with spaces",
			input:   "test name",
			wantErr: true,
		},
		{
			name:    "name with special characters",
			input:   "test@name",
			wantErr: true,
		},
		{
			name:    "name with slash",
			input:   "test/name",
			wantErr: true,
		},
		{
			name:    "name too long",
			input:   string(make([]byte, 101)),
			wantErr: true,
		},
		{
			name:    "valid name at max length",
			input:   strings.Repeat("a", 100),
			wantErr: false,
		},
		{
			name:    "invalid chars at max length",
			input:   string(make([]byte, 100)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBookmarkName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBookmarkName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDirectory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "current directory",
			input:   ".",
			wantErr: false,
		},
		{
			name:    "parent directory",
			input:   "..",
			wantErr: false,
		},
		{
			name:    "tmp directory",
			input:   "/tmp",
			wantErr: true,
		},
		{
			name:    "home directory",
			input:   "~",
			wantErr: true,
		},
		{
			name:    "non-existent directory",
			input:   "/non/existent/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateDirectory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDirectorySymlink(t *testing.T) {
	tmpDir := t.TempDir()

	symlinkPath := filepath.Join(tmpDir, "symlink")
	targetDir := t.TempDir()

	err := os.Symlink(targetDir, symlinkPath)
	if err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	_, err = validateDirectory(symlinkPath)
	if err == nil {
		t.Error("validateDirectory() should reject symlinks")
	}
}

func TestValidateDirectoryFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile")

	err := os.WriteFile(filePath, []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = validateDirectory(filePath)
	if err == nil {
		t.Error("validateDirectory() should reject files")
	}
}

func TestBookmarkPretty(t *testing.T) {
	tests := []struct {
		name string
		bm   *Bookmark
	}{
		{
			name: "nil bookmark",
			bm:   nil,
		},
		{
			name: "normal bookmark",
			bm: &Bookmark{
				Name:      "test",
				Directory: "/tmp/test",
			},
		},
		{
			name: "bookmark with spaces in directory",
			bm: &Bookmark{
				Name:      "test",
				Directory: "/tmp/test dir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bm.Pretty()
			if result == "" {
				t.Error("Pretty() should return non-empty string")
			}
		})
	}
}
