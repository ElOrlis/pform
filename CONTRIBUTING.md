# Contributing to pform

Thank you for your interest in contributing to pform! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Code Style](#code-style)
- [LLM Usage Policy](#llm-usage-policy)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/pform.git
   cd pform
   ```
3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/ElOrlis/pform.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Verify Your Setup

Run the tests to ensure everything is working:

```bash
go test ./...
```

## Making Changes

1. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the [Code Style](#code-style) guidelines

3. Write or update tests as needed

4. Commit your changes with clear, descriptive commit messages:
   ```bash
   git commit -m "Add feature: brief description of changes"
   ```

## Testing

All contributions must include appropriate tests. Run the full test suite before submitting:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Requirements

- All new functionality must have corresponding tests
- All tests must pass before submitting a PR
- Aim for high test coverage, especially for edge cases
- Use table-driven tests where appropriate

## Submitting a Pull Request

1. Ensure all tests pass locally

2. Run the linters:
   ```bash
   go vet ./...
   staticcheck ./...
   ```

3. Format your code:
   ```bash
   gofmt -w .
   ```

4. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

5. Open a Pull Request against the `main` branch

6. Fill out the PR template with:
   - A clear description of the changes
   - The motivation for the changes
   - Any relevant issue numbers

### PR Review Process

- All PRs require at least one review before merging
- CI checks must pass (tests, linting, code validation)
- Address any feedback from reviewers promptly
- Keep PRs focused and reasonably sized

## Code Style

### General Guidelines

- Follow standard Go conventions and idioms
- Use `gofmt` to format all code
- Write clear, self-documenting code
- Add comments only where the logic isn't self-evident
- Keep functions focused and concise

### Naming Conventions

- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Use descriptive names that convey purpose
- Avoid abbreviations unless widely understood

### Error Handling

- Always handle errors explicitly
- Use the custom error types (`ParseError`, `RequiredFieldError`) where appropriate
- Wrap errors with context when propagating

### Documentation

- All exported functions, types, and constants must have GoDoc comments
- Comments should be complete sentences
- Start comments with the name of the element being documented

## LLM Usage Policy

This project has specific guidelines for LLM (Large Language Model) usage:

### Allowed

LLM assistance is permitted for:

- **Testing** - Writing and improving test cases
- **CI/CD** - Creating and maintaining GitHub Actions workflows
- **Documentation** - Writing and updating documentation (including this file)

### Not Allowed

LLM assistance should **not** be used for:

- Core library code (`decode.go`, `parsers.go`, `error.go`)
- Any modifications to the main decoding logic

The core library is intentionally human-written to ensure quality and understanding of the codebase.

## Questions?

If you have questions about contributing, feel free to open an issue for discussion.

---

Thank you for contributing to pform!
