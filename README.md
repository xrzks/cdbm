# cdbm

A CLI tool for managing directory bookmarks.

## Installation

```bash
go install github.com/xrzks/cdbm@latest
```

Add this line to your shell config (usually ~/.bashrc or ~/.zshrc)

```bash
eval "$(cdbm init zsh)"
# or for bash
eval "$(cdbm init bash)"
```

## Usage

```bash
# add a bookmark
cdbm add --name <name> --directory <path>

# list bookmarks
cdbm list

# navigate to a bookmark
cdbm <name>
```

## Configuration

cdbm stores configuration in `~/.config/cdbm/.cdbm.json` (respects `$XDG_CONFIG_HOME`).

The default configuration creates a store file at `~/.config/cdbm/store`. You can customize the store location by
setting `store_path` in the config file.
