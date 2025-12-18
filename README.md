# shellican

<div align="center">
  <img src="assets/shellican-logo.png" alt="shellican logo" width="300" />
</div>

**shellican** (*shell-i-can*) is a powerful CLI for managing runnables (scripts/commands) in collections. It allows you to organize your shell scripts, define their environment, and run them easily.

![CI](https://github.com/brsyuksel/shellican/actions/workflows/ci.yml/badge.svg)
![Release](https://github.com/brsyuksel/shellican/actions/workflows/release.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/brsyuksel/shellican)](https://goreportcard.com/report/github.com/brsyuksel/shellican)

<div align="center">
  <img src="assets/shellican.gif" alt="shellican demo" width="800" />
</div>

## Why shellican?

**shellican** bridges the gap between simple shell aliases and complex Makefile setups.

- **Beyond One-Liners & Aliases**: Instead of cluttering your `.zshrc` or `.bashrc` with forgotten aliases, turn them into organized **runnables**. Keep your logic structured, named, and easy to find.
- **Shareable Workflows**: Stop saying "copy this script and change line 5". With **import / export**, you can package your collections and share them with friends or colleagues. They get a ready-to-run environment without manual editing.
- **Self-Documenting**: Every runnable supports a `help` description and a full `README`. You (and your team) will know exactly what a script does, even months later.

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
name: "My Scripts"
help: "A collection of useful scripts"
readme: "README.md"
runnables:
  - script-a
environments:
  GLOBAL_VAR: "true"
```

**script-a/runnable.yml**
```yaml
name: "Script A"
help: "This script does something awesome"
readme: "README.md"
run: "./main.sh" # or "echo hello"
before: "echo 'Starting...'"
after: "echo 'Finished!'"
environments:
  LOCAL_VAR: "123"
```

## Examples

- [dirty-vm](https://github.com/brsyuksel/dirty-vm) - A collection for creating and managing virtual machines with QEMU, cloud-init, and networking support.

## Contributing

1. Fork the repo
2. Create feature branch
3. Commit changes
4. Push and create PR

## License

MIT
