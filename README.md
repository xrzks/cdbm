# cdbm

`cdbm` manages directory bookmarks. Save a directory under a name, then jump to
it from anywhere using your shell.

- JSON-backed store, XDG-aware paths
- Shell integration for Bash, Zsh, and Fish
- Symlink-safe, name-validated, path-normalized

## Install

```sh
go install github.com/xrzks/cdbm@latest
```

Requires Go 1.27 or newer.

### Shell setup

Add the line that matches your shell to its startup file (`~/.bashrc`,
`~/.zshrc`, or `~/.config/fish/config.fish`):

```sh
# bash
eval "$(cdbm init bash)"

# zsh
eval "$(cdbm init zsh)"

# fish
cdbm init fish | source
```

The wrapper defines a `cdbm` function. Subcommands (`cdbm add`, `cdbm list`, …)
pass through to the binary. Anything else (`cdbm <name>`) is resolved to a
bookmark and the function `cd`s you there.

Shell completion is also enabled.

## Usage

```text
cdbm [options] <name>
cdbm [options] <command> [command options]
```

```sh
cdbm add                       # bookmark the current directory (name = sanitized basename)
cdbm add --name projects       # bookmark the current directory under a custom name
cdbm add --name docs --directory ~/Documents/docs
cdbm list                      # show all bookmarks (colorized)
cdbm projects                  # jump to the "projects" bookmark
cdbm edit projects --name my-projects
cdbm remove projects           # or: cdbm rm projects
```

## Options

| Flag        | Alias | Description                                         |
| ----------- | ----- | --------------------------------------------------- |
| `--debug`   | `-D`  | append debug logs to ~/.local/state/cdbm/logs.jsonl |
| `--help`    | `-h`  | show help                                           |
| `--version` | `-v`  | print the version                                   |

## Commands

### `add`

Bookmark a directory.

| Flag                | Alias | Description                                           |
| ------------------- | ----- | ----------------------------------------------------- |
| `--name <name>`     | `-n`  | bookmark name (default: sanitized directory basename) |
| `--directory <dir>` | `-d`  | directory to bookmark (default: current directory)    |

If `--name` is omitted, the directory's basename is used with characters outside
`[a-zA-Z0-9._-]` stripped. The directory must exist, be a real directory (not a
symlink), and resolve to an absolute path.

```sh
cdbm add
cdbm add --name projects
cdbm add --name docs --directory ~/Documents/docs
```

### `list`

List bookmarks, sorted by name. Output is colorized when stdout is a TTY.

```sh
cdbm list
```

### `edit`

Rename or move a bookmark. At least one of `--name` or `--directory` is
required.

| Flag                | Alias | Description       |
| ------------------- | ----- | ----------------- |
| `--name <name>`     | `-n`  | new bookmark name |
| `--directory <dir>` | `-d`  | new directory     |

```sh
cdbm edit projects --name my-projects
cdbm edit docs --directory ~/Documents/docs-new
cdbm edit projects --name my-projects --directory ~/dev/projects-new
```

### `remove` (`rm`)

Remove a bookmark.

```sh
cdbm remove projects
cdbm rm projects
```

### `clean`

Remove bookmarks with missing directories or invalid names. By default both
checks run; use the flags to scope the run.

| Flag              | Alias | Description                                             |
| ----------------- | ----- | ------------------------------------------------------- |
| `--dry-run`       | `-n`  | show what would happen without making changes           |
| `--invalid-names` |       | only remove bookmarks with invalid names                |
| `--missing-dirs`  |       | only remove bookmarks whose directories no longer exist |

```sh
cdbm clean                       # remove invalid names + missing directories
cdbm clean --dry-run             # preview only
cdbm clean --missing-dirs        # remove only bookmarks pointing at missing dirs
cdbm clean --invalid-names --dry-run
```

### `init`

Print shell integration code for the named shell. The output is meant to be
`eval`'d (or `source`'d in Fish) from your shell startup file; see
[Shell setup](#shell-setup).

```sh
cdbm init bash
cdbm init zsh
cdbm init fish
```

### `<name>` (default action)

When the first argument is not a subcommand, `cdbm` looks it up as a bookmark
and prints `cd '<absolute path>'`. The shell wrapper `eval`s that line to change
directory.

```sh
cdbm projects
# -> cd '/Users/you/projects'
```

If the bookmarked directory has been deleted or is no longer a real directory,
the command fails with a clear error.

## Bookmark names

Names must match `^[a-zA-Z0-9._-]+$` and be at most 100 characters. Anything
else (including spaces and path separators) is rejected.

## Configuration

**Default locations:**

| File       | Path                             |
| ---------- | -------------------------------- |
| Store      | `~/.config/cdbm/store.json`      |
| Config     | `~/.config/cdbm/config.json`     |
| Debug logs | `~/.local/state/cdbm/logs.jsonl` |

**Environment variables:**

- `XDG_CONFIG_HOME` — overrides the config directory
- `XDG_STATE_HOME` — overrides the state directory used for debug logs

**Config file** (`~/.config/cdbm/config.json`) is optional. Supported keys:

```json
{
  "store_path": "/custom/path/to/store.json"
}
```

`store_path` accepts `~` and `$ENV` expansion. If the config file is missing or
unreadable, defaults are used.

## Debug mode

```sh
cdbm --debug add --name test
```

Each subcommand appends a JSON line to `~/.local/state/cdbm/logs.jsonl` with an
ISO-8601 UTC timestamp, the action name (`add`, `edit`, `remove`, `clean`), and
any relevant details (name, directory, old/new values, dry-run flag, etc.).

## Security

- Symlinks are rejected at every path boundary — bookmarks can only point at
  real directories.
- Bookmark names are validated against a strict regex; directories are validated
  and normalized to absolute paths via `filepath.Abs` + `filepath.Clean`.
- The store file and debug log are written with mode `0o600`.
- Shell output is single-quoted, with embedded `'` escaped, before being handed
  to `eval` by the shell wrapper.

## Development

```sh
go test ./...
go vet ./...
go build .
```

`cdbm --version` reports the version Go records at install time. To stamp a
version into a local build, use `go build -ldflags "-X main.version=v1.0.0" .`.

The root `main.go` loads the config and store and hands them to the command.
`internal/cli` defines the command and subcommands using `urfave/cli/v3`,
`internal/config` resolves XDG paths and `config.json`, `internal/store`
persists bookmarks, and `internal/logger` writes debug logs.

## License

MIT. See [LICENSE](LICENSE).
