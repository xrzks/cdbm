# cdbm - **C**hange **D**irectory to **B**ook**M**ark

CLI tool for managing directory bookmarks.

## Installation

Requires Go 1.25.5 or later.

```bash
go install github.com/xrzks/cdbm@latest
```

Add this line to your shell config (usually ~/.bashrc or ~/.zshrc):

```bash
eval "$(cdbm init zsh)"
# or for bash
eval "$(cdbm init bash)"
```

## Usage

```bash
# add a bookmark
cdbm add --name projects --directory ~/dev/projects

# list bookmarks
cdbm list

# navigate to a bookmark
cdbm projects
```

**Bookmark names**: Only letters, numbers, `.`, `_`, and `-` (max 100 chars).

## Configuration

Store file created at `~/.config/cdbm/store.json` on first `cdbm add`.

Config file: `~/.config/cdbm/.cdbm.json` (respects `$XDG_CONFIG_HOME`). Customize `store_path` there.
