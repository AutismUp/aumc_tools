# CI/CD Setup Documentation

## Overview

This repository is now configured with GitHub Actions for automated testing and code quality checks. Tests will run automatically on every push and pull request to the `main` and `python-to-go` branches.

## Workflow Configuration

### File: `.github/workflows/test.yml`

The workflow includes two jobs:

#### 1. Test Job
- **Runs on**: Ubuntu Latest
- **Steps**:
  - Checkout code
  - Set up Go 1.25.4 with caching
  - Download and verify dependencies
  - Run all tests (`make test`)
  - Build binary for current platform
  - Build binaries for all platforms (Linux, macOS, Windows)

#### 2. Lint Job
- **Runs on**: Ubuntu Latest
- **Steps**:
  - Checkout code
  - Set up Go 1.25.4 with caching
  - Run golangci-lint for code quality checks

## Linter Configuration

### File: `.golangci.yml`

Enabled linters:
- `errcheck` - Check for unchecked errors
- `gosimple` - Simplify code
- `govet` - Vet examines Go source code
- `ineffassign` - Detect ineffectual assignments
- `staticcheck` - Static analysis checks
- `unused` - Check for unused code
- `gofmt` - Check formatting
- `goimports` - Check imports formatting

Configuration highlights:
- 5-minute timeout
- Tests are included in linting
- Some linters are relaxed for test files (e.g., errcheck)

## Status Badge

The README now includes a status badge showing the current test status:

[![Tests](https://github.com/AutismUp/aumc_tools/actions/workflows/test.yml/badge.svg)](https://github.com/AutismUp/aumc_tools/actions/workflows/test.yml)

## Viewing Test Results

1. Go to the repository on GitHub
2. Click on the "Actions" tab
3. Select the "Tests" workflow
4. View individual workflow runs and their results

## Running Tests Locally

Before pushing, you can run the same checks locally:

```bash
# Run tests
make test

# Build binary
make build

# Build for all platforms
make build-all

# Run linter (requires golangci-lint installed)
golangci-lint run
```

### Installing golangci-lint locally

```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Or using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Troubleshooting

### Workflow fails on dependency download
- Check that `go.mod` and `go.sum` are committed
- Run `go mod tidy` locally and commit changes

### Linter fails
- Run `golangci-lint run` locally to see issues
- Fix issues or update `.golangci.yml` if needed

### Build fails for specific platform
- Check the build logs in GitHub Actions
- Test cross-compilation locally: `GOOS=linux GOARCH=amd64 go build ./cmd/aumc`

## Future Enhancements

Potential improvements to consider:
- Add code coverage reporting
- Add release automation
- Add Docker image building
- Add integration tests with test environment
- Add security scanning
