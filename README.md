# Glocate

**A simple cli tool for indexing/searching files in your filesystem**

`glocate` is an alternative to the locate/updatedb commands written in Go.

## Installation

```bash
go install github.com/fedeztk/glocate@latest
```

## Usage

Create the index database
```bash
glocate --index
```

Search for a pattern (regex are supported out of the box), by default it is case sensitive.
```bash
glocate "pattern"
```

```bash
glocate --smartcase "pattern" # case insensitive if the pattern is all lowercase
```

For a full list of options and shortcuts see the help page
```bash
glocate --help
```

### Configuration

The configuration can be done via environment variables, flags, and a config file.

The config file is in yaml format. It will be created automatically if it does not exist under `$HOME/.config/glocate/glocate.yaml`. Default values are shown below.
```yaml
directories: # directories to index
  - "$HOME"

ignoredPatterns: # patterns to ignore
  - "$HOME/.cache"

ignoreSymlinks: true # do not follow symlinks
ignoreHidden: false # ignore hidden files
```

### Acknowledgements

A special thanks to the creator of the walk implementation used to walk the filesystem ([see here](https://github.com/opencoff/go-walk)).

Acknowledgements also go to the creators of all the other libraries used in this project, see [go.mod](go.mod) for a full list.