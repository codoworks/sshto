# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- GoReleaser configuration for cross-platform binary releases
- GitHub Actions release workflow triggered on version tags

## [0.1.0] - 2024-12-14

### Added

- Interactive fuzzy finder for quick server selection using Bubbletea TUI
- Server management commands: `add`, `edit`, `remove`, `list`
- Group support with color-coded organization (red, green, yellow, blue, magenta, cyan, white, gray)
- YAML-based configuration stored at `~/.config/sshto/config.yaml`
- Direct connection via `sshto <server>` command
- User override support with `-u` flag
- Group filtering with `-g` flag for list command
- Default settings for user, port, and SSH key
- Input validation for hosts, ports, and server names
- Beautiful terminal UI with Lipgloss styling

[Unreleased]: https://github.com/codoworks/sshto/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/codoworks/sshto/releases/tag/v0.1.0
