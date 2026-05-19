package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger interface {
	Log(action string, details map[string]any) error
}

type FileLogger struct {
	path string
	file *os.File
}

type LogEntry struct {
	Timestamp string         `json:"timestamp"`
	Action    string         `json:"action"`
	Details   map[string]any `json:"details,omitempty"`
}

func NewFileLogger(path string) (*FileLogger, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileLogger{
		path: path,
		file: file,
	}, nil
}

func (l *FileLogger) Log(action string, details map[string]any) error {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Action:    action,
		Details:   details,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to serialize log entry: %w", err)
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry to file: %w", err)
	}

	return nil
}

func (l *FileLogger) Close() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
		l.file = nil
	}
	return nil
}
