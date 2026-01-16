# AGENTS.md

## Project Overview

cdbm is a Go CLI tool for managing directory bookmarks using urfave/cli/v3.

**Module**: `github.com/xrzks/cdbm`
**Go Version**: 1.25.5

## Essential Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/cli

# Run a single test
go test ./internal/cli -run TestFunctionName

# Run tests with verbosity
go test ./... -v

# Build the binary
go build ./cmd/cdbm

# Format code check (should produce no output)
gofmt -l .

# Run vet
go vet ./...

# Install globally
go install github.com/xrzks/cdbm@latest
```

## Code Style Guidelines

### Import Ordering
Imports are grouped with blank lines between groups:
1. Standard library (alphabetical order)
2. Third-party packages (alphabetical order)
3. Internal packages (github.com/xrzks/cdbm/...)

Example:
```go
package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/xrzks/cdbm/internal/config"
)
```

### Naming Conventions
- **Exported**: PascalCase (e.g., `NewStore`, `GetOne`, `Bookmark`)
- **Private**: camelCase (e.g., `loadFile`, `bookmarks`)
- **Constants**: PascalCase (e.g., `MaxNameLength`)
- **Interfaces**: PascalCase with -er suffix
- **File names**: Match primary functionality, test files as `<name>_test.go`

### Types and Structs
- Exported struct types use PascalCase
- Private struct fields use camelCase
- Constructor functions: `New<TypeName()`
- Receiver names are short abbreviations (s, c, bm)

### Error Handling
- Always return errors, never ignore them
- Use `fmt.Errorf` with `%w` verb for error wrapping
- Provide context-specific error messages
- Check `os.IsNotExist()` for missing files, return defaults when appropriate

```go
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}
```

### Security Requirements
- Use `os.Lstat()` (not `os.Stat()`) to detect symlinks
- Validate bookmark names with regex `^[a-zA-Z0-9._-]+$` (max 100 chars)
- Shell quoting: `' + strings.ReplaceAll(s, "'", "'\\''") + "'`
- Write sensitive files with `0o600` permissions
- All directory paths must be absolute (use `filepath.Abs()`)
- Paths must be cleaned with `filepath.Clean()`

### File I/O Pattern
```go
// Reading
data, err := os.ReadFile(path)
if err != nil {
    if os.IsNotExist(err) {
        return defaultConfig, nil  // or specific error
    }
    return nil, fmt.Errorf("failed to read: %w", err)
}

// Writing
err = os.WriteFile(path, data, 0o600)
```

### Command Pattern (CLI)
```go
func (c *CLI) New<Command>Command() *cli.Command {
    return &cli.Command{
        Name:   "command-name",
        Usage:  "description",
        Action: c.Run<Command>Command,
    }
}

func (c *CLI) Run<Command>Command(ctx context.Context, cmd *cli.Command) error {
    // Implementation
}
```

## Project Structure

```
cdbm/
├── cmd/cdbm/main.go         # Entry point
├── internal/
│   ├── cli/                 # CLI commands
│   ├── config/              # Config loading
│   └── store/               # Data persistence
└── go.mod
```

## Important Notes

- Store returns empty map `{}` if file doesn't exist (not an error)
- Config is optional - defaults work without config file
- Entry point: `cmd/cdbm/main.go`
- Default store: `~/.config/cdbm/store.json`
- Module import: `github.com/xrzks/cdbm/internal/...`

## Development Workflow

1. Make changes
2. Run tests: `go test ./...`
3. Format check: `gofmt -l .`
4. Run vet: `go vet ./...`
5. Build: `go build ./cmd/cdbm`
