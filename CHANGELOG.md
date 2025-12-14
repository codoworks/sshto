# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2025-12-14

### Fixed

- Fix Homebrew cask configuration (use default `Casks` directory instead of `Formula`)
- Add `binaries` field to properly install CLI binary
- Add caveats with usage instructions

### Changed

- Expand `homebrew_casks` config with commit author, message template, and skip_upload

## [0.3.0] - 2025-12-14

### Changed

- Change Homebrew directory from `Casks` to `Formula` (CLI tools use Formula, not Casks)

## [0.2.0] - 2025-12-14

### Changed

- Migrate from `brews` to `homebrew_casks` in GoReleaser (deprecation fix)
- Update `archives.format` to `archives.formats` list syntax (deprecation fix)

## [0.1.0] - 2025-12-14

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
- GoReleaser configuration for cross-platform binary releases
- GitHub Actions release workflow triggered on version tags
- Homebrew tap support via `brew tap codoworks/tap && brew install sshto`

[Unreleased]: https://github.com/codoworks/sshto/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/codoworks/sshto/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/codoworks/sshto/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/codoworks/sshto/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/codoworks/sshto/releases/tag/v0.1.0
