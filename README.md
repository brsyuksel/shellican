# Shellican

![CI](https://github.com/brsyuksel/shellican/actions/workflows/ci.yml/badge.svg)
![Release](https://github.com/brsyuksel/shellican/actions/workflows/release.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/brsyuksel/shellican)](https://goreportcard.com/report/github.com/brsyuksel/shellican)

Shellican (shell-i-can) is a powerful CLI for managing runnables (scripts/commands) in collections. It allows you to organize your shell scripts, define their environment, and run them easily.

## Features

- **Collections & Runnables**: Organize your scripts into collections.
- **YAML Configuration**: Define runnables and environments in `collection.yml` and `runnable.yml`.
- **Environment Management**: Inject environment variables defined in configuration.
- **Hooks**: Pre (`before`) and post (`after`) hooks for runnables.
- **Shell Helper**: Generate shell wrappers for easy access.
- **Import/Export**: Share collections easily.

## Installation

### Go Install

```bash
go install github.com/brsyuksel/shellican@latest
```

### Binary Releases

You can download the pre-compiled binary for your operating system (macOS, Linux) from the [Releases](https://github.com/brsyuksel/shellican/releases) page.

## Usage

### Structure
Default storage location: `~/.shellican` (or `$SHELLICAN_HOME/.shellican`).

```
~/.shellican/
  ├── my-collection/
  │   ├── collection.yml
  │   ├── script-a/
  │   │   └── runnable.yml
  │   └── script-b/
  │       ├── runnable.yml
  │       └── main.sh
```

### Commands

- **New Collection**: `shellican new collection <name>`
- **New Runnable**: `shellican new runnable <collection> <name>`
- **Run**: `shellican run <collection> <runnable> [args...]`
- **List**: `shellican list collections` or `shellican list runnables <collection>`
- **Show**: `shellican show collection <name>` or `shellican show runnable <collection> <name>`
- **Create Shell Shortcut**: `shellican create-shell <collection>` (creates `~/.local/bin/<collection>-shell`)
- **Import/Export**: `shellican import <source>` / `shellican export <collection>`

## Configuration

**collection.yml**
```yaml
summary: "My Scripts"
help: "A collection of useful scripts"
runnables:
  - script-a
environments:
  GLOBAL_VAR: "true"
```

**runnable.yml**
```yaml
summary: "Script A"
type: script # or inline
run: "./main.sh" # or "echo hello" if inline
environments:
  LOCAL_VAR: "123"
```

## Contributing

1. Fork the repo
2. Create feature branch
3. Commit changes
4. Push and create PR

## License

MIT
