# AGENTS.md

This document provides essential information for working effectively in the cdbm codebase.

## Project Overview

cdbm is a Go CLI tool for managing directory bookmarks. It allows users to save frequently used directories and quickly navigate to them via the command line with shell integration.

**Module**: `github.com/xrzks/cdbm`
**Go Version**: 1.25.5
**CLI Framework**: urfave/cli v3

## Essential Commands

### Build and Test
```bash
# Run all tests
go test ./...

# Build the binary
go build ./cmd/cdbm

# Install globally (for users)
go install github.com/xrzks/cdbm@latest
```

### Code Quality
```bash
# Format code (gofmt - should produce no output)
gofmt -l .

# Run vet
go vet ./...
```

## Project Structure

```
cdbm/
├── cmd/
│   └── cdbm/
│       └── main.go          # Application entry point
├── internal/
│   ├── cli/                 # CLI commands and logic
│   │   ├── add.go           # Add bookmark command
│   │   ├── init.go          # Shell initialization command
│   │   ├── list.go          # List bookmarks command
│   │   ├── root.go          # Root command (navigation)
│   │   └── root_test.go     # Tests for shell quoting
│   └── store/               # Data persistence
│       ├── bookmark.go      # Bookmark struct and pretty printing
│       ├── file.go          # File I/O operations
│       └── store.go         # Store logic and validation
├── main.go                  # Symlink to cmd/cdbm/main.go
├── store                    # JSON storage file (created at runtime)
├── go.mod
└── go.sum
```

## Code Organization

### Package Structure
- **`cmd/cdbm`**: Application entry point - initializes store and runs CLI
- **`internal/cli`**: Command definitions and handlers using urfave/cli/v3
- **`internal/store`**: Bookmark data persistence with validation

### Key Types
- `CLI` (cli package): Main CLI struct holding store reference
- `Store` (store package): In-memory bookmark map with file persistence
- `Bookmark` (store package): Data model with Name and Directory fields

### Command Pattern
Commands follow this pattern:
```go
func (c *CLI) New<Command>Command() *cli.Command {
    return &cli.Command{
        Name:   "command-name",
        Usage:  "description",
        Flags:  []cli.Flag{...},
        Action: c.Run<Command>Command,
    }
}

func (c *CLI) Run<Command>Command(ctx context.Context, cmd *cli.Command) error {
    // Implementation
}
```

## Naming Conventions

### Go Conventions
- **Exported**: PascalCase (e.g., `NewStore`, `GetOne`, `Bookmark`)
- **Private**: camelCase (e.g., `loadFile`, `writeFile`, `bookmarks`)
- **Constants**: PascalCase
- **Interfaces**: PascalCase with -er suffix (not currently used in codebase)

### File Naming
- Files match the primary type or functionality they contain
- Command files: `add.go`, `list.go`, `init.go`
- Test files: `<name>_test.go` in same package

## Code Style Patterns

### Error Handling
- Always return errors, never ignore them
- Use `fmt.Errorf` with `%w` verb for error wrapping
- Provide context-specific error messages
- Example:
  ```go
  return fmt.Errorf("failed to marshal bookmarks: %w", err)
  ```

### Security Considerations

**Bookmark Name Validation** (store.go:29-41):
- Regex: `^[a-zA-Z0-9._-]+$`
- Max length: 100 characters
- Cannot be empty

**Symlink Protection**:
- Uses `os.Lstat()` instead of `os.Stat()` to check for symlinks without following them
- Rejects symlinks in both `Add()` and `RunRootCommand()`
- Prevents TOCTOU (Time-of-Check-Time-of-Use) attacks

**Shell Output Quoting** (root.go:70-72):
- Single quotes with proper escaping: `"' + strings.ReplaceAll(s, "'", "'\\''") + "'"`
- Prevents shell injection attacks
- Comprehensive test coverage in `root_test.go`

**File Permissions**:
- Store file written with `0o600` (read/write owner only)
- Prevents information leakage

### File I/O Pattern
```go
// Reading
data, err := os.ReadFile(path)
if err != nil {
    if os.IsNotExist(err) {
        return nil, fmt.Errorf("specific message")
    }
    return nil, fmt.Errorf("generic error: %w", err)
}

// Writing
err = os.WriteFile(path, data, 0o600)
```

## Testing Approach

### Test Location
- Tests are co-located with source files in same package
- Example: `internal/cli/root_test.go`

### Current Test Coverage
- Only `internal/cli` package has tests (shell quoting safety)
- Security-focused tests: injection prevention, unicode handling, special characters

### Running Tests
```bash
# All tests
go test ./...

# Specific package
go test ./internal/cli

# With verbosity
go test ./... -v
```

### Test Patterns
- Table-driven tests for multiple scenarios
- Test names describe the scenario
- Test both happy path and edge cases
- Security tests verify invariants (e.g., "shell injection should be prevented")

## Important Gotchas

### Directory Structure
- The root `main.go` is a **symlink** to `cmd/cdbm/main.go`
- Don't modify the symlink - edit the actual file in `cmd/cdbm/`

### Storage Location
- Bookmarks stored in a file named `store` in the project root
- Format: JSON map of bookmark name to bookmark object
- Example structure:
  ```json
  {
    "name1": {"Name":"name1","Directory":"/absolute/path"},
    "name2": {"Name":"name2","Directory":"/another/path"}
  }
  ```

### Store Initialization
- Store file is created lazily on first `Add()` operation
- If store file doesn't exist, operations fail with message: "store file not found. Run 'cdbm init' to set up the application"
- `cdbm init` actually generates shell integration code, not the store file

### Shell Integration
- The `init` command outputs shell function code, not a file
- Users eval this output in their shell config: `eval "$(cdbm init <shell>)"`
- Shell function decides whether to execute commands or eval output based on first argument

### Path Handling
- All directory paths are converted to absolute paths with `filepath.Abs()`
- Paths are cleaned with `filepath.Clean()`
- Pretty printing shows only basename of directory (bookmark.go:19)

## Dependencies

### Runtime
- `github.com/urfave/cli/v3` v3.6.1 - CLI framework

### Standard Library
- `encoding/json` - Store file serialization
- `os` - File I/O, path operations
- `path/filepath` - Path manipulation
- `regexp` - Bookmark name validation
- `strings` - String manipulation
- `testing` - Unit tests
- `context` - Context propagation
- `fmt` - Formatted I/O
- `log` - Logging (fatal errors only)

## Development Workflow

1. Make changes to source files
2. Run tests: `go test ./...`
3. Run format check: `gofmt -l .`
4. Run vet: `go vet ./...`
5. Build binary to verify: `go build ./cmd/cdbm`

## Adding New Commands

To add a new CLI command:

1. Create `<command>.go` in `internal/cli/`
2. Implement `New<Command>Command()` and `Run<Command>Command()` methods
3. Register in `root.go` Commands list:
   ```go
   Commands: []*cli.Command{
       c.NewAddCommand(),
       c.NewListCommand(),
       c.NewInitCommand(),
       c.New<Command>Command(),  // Add here
   },
   ```
4. Add tests in `<command>_test.go` or `root_test.go`

## Adding New Store Operations

To add new store operations:

1. Implement method in `internal/store/store.go`
2. Validate bookmark names using `validateBookmarkName()`
3. Apply security checks (no symlinks, paths must exist)
4. Persist changes with `writeFile()` after modifications
5. Consider adding tests in `store_test.go`

## Module Information

- **Module Path**: `github.com/xrzks/cdbm`
- **Import Path**: Use `github.com/xrzks/cdbm/internal/...` for internal packages
- **Version**: Latest
