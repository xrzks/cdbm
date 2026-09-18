# cdbm

A small CLI for managing directory bookmarks. Save a directory under a name, then jump to it from anywhere using your shell.

- Single static binary, written in Go
- JSON-backed store, XDG-aware paths
- Shell integration for Bash, Zsh, and Fish
- Symlink-safe, name-validated, path-normalized

## Installation

Requires Go 1.26.7 or later.

```bash
go install github.com/xrzks/cdbm@latest
```

The binary is installed as `cdbm`.

## Shell setup

Add the line that matches your shell to its startup file (`~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`):

```bash
# bash
eval "$(cdbm init bash)"

# zsh
eval "$(cdbm init zsh)"

# fish
cdbm init fish | source
```

The wrapper defines a `cdbm` function. Subcommands (`cdbm add`, `cdbm list`, …) pass through to the binary. Anything else (`cdbm <name>`) is resolved to a bookmark and the function `cd`s you there.

Shell completion is also enabled.

## Quick start

```bash
cdbm add                       # bookmark the current directory (name = sanitized basename)
cdbm add --name projects       # bookmark the current directory under a custom name
cdbm add --name docs --directory ~/Documents/docs
cdbm list                      # show all bookmarks (colorized)
cdbm projects                  # jump to the "projects" bookmark
cdbm edit projects --name my-projects
cdbm remove projects           # or: cdbm rm projects
```

## Commands

### `add`

Bookmark a directory.

| Flag | Alias | Description |
| --- | --- | --- |
| `--name` | `-n` | Bookmark name. Defaults to the sanitized basename of the directory. |
| `--directory` | `-d` | Directory to bookmark. Defaults to the current working directory. |

If `--name` is omitted, the directory's basename is used with characters outside `[a-zA-Z0-9._-]` stripped. The directory must exist, be a real directory (not a symlink), and resolve to an absolute path.

```bash
cdbm add
cdbm add --name projects
cdbm add --name docs --directory ~/Documents/docs
```

### `list`

Print all bookmarks, sorted by name. Output is colorized when stdout is a TTY.

```bash
cdbm list
```

### `edit`

Rename or move an existing bookmark. At least one of `--name` or `--directory` is required.

| Flag | Alias | Description |
| --- | --- | --- |
| `--name` | `-n` | New bookmark name. |
| `--directory` | `-d` | New directory. |

```bash
cdbm edit projects --name my-projects
cdbm edit docs --directory ~/Documents/docs-new
cdbm edit projects --name my-projects --directory ~/dev/projects-new
```

### `remove` (`rm`)

Delete a bookmark.

```bash
cdbm remove projects
cdbm rm projects
```

### `clean`

Remove bookmarks that are no longer useful: those pointing at missing directories or with invalid names.

By default both checks run. Use the flags to scope the run.

| Flag | Alias | Description |
| --- | --- | --- |
| `--dry-run` | `-n` | Print what would be removed, but change nothing. |
| `--invalid-names` | | Only remove bookmarks whose names fail validation. |
| `--missing-dirs` | | Only remove bookmarks whose directories no longer exist. |

```bash
cdbm clean                       # remove invalid names + missing directories
cdbm clean --dry-run             # preview only
cdbm clean --missing-dirs        # remove only bookmarks pointing at missing dirs
cdbm clean --invalid-names --dry-run
```

### `init`

Emit shell integration code for the named shell. The output is meant to be `eval`'d (or `source`'d in Fish) from your shell startup file.

```bash
cdbm init bash
cdbm init zsh
cdbm init fish
```

### `<name>` (default action)

When the first argument is not a subcommand, `cdbm` looks it up as a bookmark and prints `cd '<absolute path>'`. The shell wrapper `eval`s that line to change directory.

```bash
cdbm projects
# -> cd '/Users/you/projects'
```

If the bookmarked directory has been deleted or is no longer a real directory, the command fails with a clear error.

## Global flags

| Flag | Alias | Description |
| --- | --- | --- |
| `--debug` | `-D` | Append structured debug logs to `~/.local/state/cdbm/logs.jsonl` for this invocation. |

## Bookmark names

Names must match `^[a-zA-Z0-9._-]+$` and be at most 100 characters. Anything else (including spaces and path separators) is rejected.

## Configuration

**Default locations:**

| File | Path |
| --- | --- |
| Store | `~/.config/cdbm/store.json` |
| Config | `~/.config/cdbm/.cdbm.json` |
| Debug logs | `~/.local/state/cdbm/logs.jsonl` |

**Environment variables:**

- `XDG_CONFIG_HOME` — overrides the config directory
- `XDG_STATE_HOME` — overrides the state directory used for debug logs

**Config file** (`~/.config/cdbm/.cdbm.json`) is optional. Supported keys:

```json
{
  "store_path": "/custom/path/to/store.json"
}
```

`store_path` accepts `~` and `$ENV` expansion. If the config file is missing or unreadable, defaults are used.

## Debug mode

```bash
cdbm --debug add --name test
```

Each subcommand appends a JSON line to `~/.local/state/cdbm/logs.jsonl` with an ISO-8601 UTC timestamp, the action name (`add`, `edit`, `remove`, `clean`), and any relevant details (name, directory, old/new values, dry-run flag, etc.).

## Security

- Symlinks are rejected at every path boundary — bookmarks can only point at real directories.
- Bookmark names are validated against a strict regex; directories are validated and normalized to absolute paths via `filepath.Abs` + `filepath.Clean`.
- The store file and debug log are written with mode `0o600`.
- Shell output is single-quoted, with embedded `'` escaped, before being handed to `eval` by the shell wrapper.

## License

MIT — see [LICENSE](./LICENSE).