# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o sshto .      # Build binary
go run .                 # Run directly
go mod tidy              # Update dependencies
```

## Architecture

```
cmd/           # Cobra CLI commands (thin layer, orchestrates app)
internal/
  app/         # Application orchestration, dependency injection
  config/      # YAML config load/save, Server/Group models
  ssh/         # SSH command execution
  ui/          # Bubbletea TUI components (list, form, styles)
```

**Flow**: `cmd` → `app.App` → `config.Config` + `ssh.Client` + `ui.*`

## Key Patterns

- **Config location**: `~/.config/sshto/config.yaml`
- **Models**: `config.Server` and `config.Group` with YAML tags
- **App struct**: Holds Config and SSHClient, resolves defaults/overrides
- **TUI**: Bubbletea models in `internal/ui/` with Lipgloss styling

## CLI Usage

```bash
sshto                     # Interactive fuzzy finder
sshto <server>            # Direct connect
sshto <server> -u root    # Connect with user override
sshto list -g production  # Filter by group
sshto add                 # Interactive add form
sshto edit <server>       # Interactive edit form
sshto remove <server>     # Remove with confirmation
sshto groups              # List groups
sshto groups add <name>   # Add group
```

## Config Schema

```yaml
groups:
  - name: production
    color: red           # red, green, yellow, blue, magenta, cyan, white, gray

servers:
  - name: web-prod
    host: 192.168.1.10
    user: deploy
    port: 22
    key: ~/.ssh/id_rsa
    group: production

defaults:
  user: ""
  port: 22
  key: ""
```
