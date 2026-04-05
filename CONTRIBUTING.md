# Contributing to gwiz

Thank you for your interest in contributing to gwiz! This document provides
guidelines and information for contributors.

## Reporting Bugs

- Check existing issues to avoid duplicates
- Include Go version, terminal emulator, and OS
- Provide a minimal reproducible example
- Include terminal screenshots if the issue is visual

## Suggesting Features

- Open an issue describing the feature and its use case
- For new step types, include a sketch of the expected API

## Pull Requests

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run all checks (see below)
5. Submit a pull request with a clear description

## Development Guidelines

### Code Style

- Run `gofmt` on all code before committing
- Run `go vet ./...` to catch common issues
- Follow existing conventions in the codebase:
  - Short receiver names (1-2 characters): `(w *Wizard)`, `(s *SelectStep)`, `(t Theme)`
  - Godoc comments on all exported types and functions
  - Godoc comments start with the name being documented

### Testing

- Add tests for new functionality
- Use table-driven tests where appropriate
- Run `go test ./...` before submitting
- Aim for high coverage on step types

### Commit Messages

- Use present tense ("Add feature" not "Added feature")
- Start with a verb ("Fix bug" not "Bug fix")
- Reference issues where applicable ("Fix #42")
- Keep the first line under 72 characters

## Project Structure

```
gwiz/
├── doc.go                # Package documentation
├── gwiz.go               # Wizard orchestrator
├── step.go               # Step interface & base types
├── state.go              # State key-value bag
├── nav.go                # Step navigation logic
├── chrome.go             # Border, header, and nav bar rendering
├── theme.go              # Color themes
├── option.go             # Option type for selections
├── input.go              # Text input step
├── select.go             # Single-select step
├── multiselect.go        # Multi-select step
├── form.go               # Multi-field form step
├── info.go               # Read-only info display
├── confirm.go            # Confirmation step
├── table.go              # Tabular data step
├── exec.go               # Long-running task step
├── *_test.go             # Tests
└── examples/             # Example applications
    ├── simple/            # Basic wizard
    └── installer/         # Multi-step installer
```
