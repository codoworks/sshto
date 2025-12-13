# Contributing to sshto

Thank you for your interest in contributing to sshto! This document provides guidelines and instructions for contributing.

## Prerequisites

- Go 1.24 or later
- Git

## Development Setup

1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/sshto.git
   cd sshto
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Code Style

- Run `make fmt` before committing
- Run `make vet` to check for common issues
- Run `make lint` for comprehensive linting (requires golangci-lint)
- Follow standard Go conventions and idioms

## Making Changes

1. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit with clear, descriptive messages

3. Ensure all tests pass:
   ```bash
   make test
   ```

4. Push to your fork and open a pull request

## Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Ensure CI checks pass
- Keep changes focused and atomic

## Reporting Issues

When reporting issues, please include:

- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Go version and operating system
- Relevant configuration (with sensitive data removed)

## Questions?

Feel free to open an issue for questions or discussion.
