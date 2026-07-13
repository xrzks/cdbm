# cdbm - cd to bookmark

CLI tool for managing directory bookmarks.

## Installation

Requires Go 1.26.4 or later.

```bash
go install github.com/xrzks/cdbm@latest
```

## Shell Setup

Add this line to your shell config (usually `~/.bashrc`, `~/.zshrc`, `~/.config/fish/config.fish`):

```bash
# bash
eval "$(cdbm init bash)"
# zsh
eval "$(cdbm init zsh)"
# fish
cdbm init fish | source
```

## Usage

```bash
# add a bookmark (current directory)
cdbm add

# add with custom name
cdbm add --name projects

# add with custom directory
cdbm add --name docs --directory ~/Documents/docs

# list all bookmarks
cdbm list

# navigate to a bookmark
cdbm projects

# edit a bookmark (rename)
cdbm edit projects --name my-projects

# edit a bookmark (move)
cdbm edit docs --directory ~/Documents/docs-new

# edit a bookmark (rename and move)
cdbm edit projects --name my-projects --directory ~/dev/projects-new

# remove a bookmark
cdbm remove projects
```

## Bookmark Names

Only letters, numbers, `.`, `_`, and `-` (max 100 chars).

## Configuration

**Default locations:**

- Store file: `~/.config/cdbm/store.json`
- Config file: `~/.config/cdbm/.cdbm.json`
- Debug logs: `~/.local/state/cdbm/logs.jsonl`

**Environment variables:**

- `XDG_CONFIG_HOME`: Override config directory
- `XDG_STATE_HOME`: Override state directory (for debug logs)

**Config options** (in `.cdbm.json`):

```json
{
  "store_path": "/path/to/store.json"
}
```

## Debug Mode

Enable debug logging to troubleshoot issues:

```bash
cdbm --debug add --name test
```

Logs are written to `~/.local/state/cdbm/logs.jsonl`.

## Security

- Symlinks are rejected for security
- All directory paths are validated and normalized
- Sensitive files are written with restricted permissions (0o600)

## Roadmap

- [x] Add, list, edit, remove, and navigate to bookmarks
- [x] Clean invalid bookmarks and missing directories, with dry-run support
- [x] Support Bash, Zsh, and Fish shell integration
- [x] Support XDG paths, custom store locations, and debug logging
- [ ] Add shell completion for bookmark names
- [ ] Add bookmark search and fuzzy matching
- [ ] Add import, export, and backup commands
