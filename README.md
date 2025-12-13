# sshto

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/codoworks/sshto/actions/workflows/ci.yml/badge.svg)](https://github.com/codoworks/sshto/actions/workflows/ci.yml)

A fast, interactive SSH connection manager with fuzzy search.

## Features

- Interactive fuzzy finder for quick server selection
- Organize servers into color-coded groups
- Override connection parameters on the fly
- YAML-based configuration
- Beautiful terminal UI

## Installation

### Go Install

```bash
go install github.com/codoworks/sshto@latest
```

### Build from Source

```bash
git clone https://github.com/codoworks/sshto.git
cd sshto
make build
# Binary will be at bin/sshto
```

### Releases

Pre-built binaries for macOS, Linux, and Windows are available on the [Releases](https://github.com/codoworks/sshto/releases) page.

## Usage

```bash
sshto                     # Interactive fuzzy finder
sshto <server>            # Direct connect
sshto <server> -u root    # Connect with user override
sshto list                # List all servers
sshto list -g production  # Filter by group
sshto add                 # Interactive add form
sshto edit <server>       # Interactive edit form
sshto remove <server>     # Remove with confirmation
sshto groups              # List groups
sshto groups add <name>   # Add group
```

## Configuration

Configuration is stored at `~/.config/sshto/config.yaml`.

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

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
