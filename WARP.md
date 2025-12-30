# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

This is **aumc_tools** (Autism Up Minecraft Tool) - a CLI wrapper for managing Minecraft servers using MSM (Minecraft Server Manager). The project is currently in migration from Python to Go, with the Go implementation being the active development focus.

## Key Commands

### Building
```bash
make build              # Build binary for current platform (output: bin/aumc)
make build-all          # Build for Linux, macOS (amd64/arm64), Windows
make build-linux        # Build for Linux (required for Docker testing)
make clean              # Remove build artifacts
```

### Testing
```bash
make test               # Run all Go tests with verbose output
```

### Linting
```bash
golangci-lint run       # Run linters (uses .golangci.yml config)
```

### Running the Binary
```bash
bin/aumc --help         # Run built binary
bin/aumc init           # Initialize default config files
```

### Docker Test Environment
```bash
./docker-test.sh        # Build binary, build image, start container, run tests

# Access container shell
docker-compose exec minecraft-test bash

# Run aumc commands in container
docker-compose exec minecraft-test aumc --help
docker-compose exec minecraft-test aumc init

# Container management
docker-compose up -d    # Start container
docker-compose down     # Stop and remove container
docker-compose restart  # Restart container
docker-compose logs -f  # View container logs

# Rebuild everything
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

### Testing Code Changes with Docker
```bash
# 1. Make your code changes
# 2. Rebuild Linux binary
make build-linux

# 3. Copy new binary to container and test
docker-compose exec minecraft-test sudo cp /workspace/bin/aumc-linux-amd64 /usr/local/bin/aumc
docker-compose exec minecraft-test aumc --help

# 4. Run full test suite
docker-compose exec minecraft-test bash /workspace/test-aumc.sh
```

### Dependencies
```bash
make deps               # Download Go module dependencies
make tidy               # Tidy and verify Go module dependencies
```

## Architecture

### Project Structure
```
cmd/
  aumc/main.go         # Entry point - calls cmd.Execute()
  root.go              # Cobra root command setup, config initialization
  build.go             # Commands: build-new-jar, publish-new-jar
  world.go             # Commands: create-new-world, delete-world
  config.go            # Commands: check-config, reload-config
  init.go              # Command: init (creates default config files)

internal/
  config/              # Configuration management using Viper
    config.go          # Config structs (Config, BuildConfig, WorldConfig)
    templates.go       # Default config.json and server.properties templates
  
  minecraft/           # Core business logic
    aumc.go            # AuMc struct - main operations (BuildNewJar, etc.)
    minecraft.go       # Minecraft-specific utilities
    utils.go           # File system helpers
  
  mcprops/             # Minecraft server.properties parser
    mcconfig.go        # MCConfig struct - read/write server.properties
  
  errors/              # Custom error types for domain-specific errors

bin/                   # Build output directory (gitignored)
```

### Configuration System

The tool uses **Viper** for configuration management. Configuration is loaded via:
1. `--config` flag (highest priority)
2. `AU_CONFIG_FILE` environment variable
3. `./config.json` in current directory (fallback)

Configuration structure:
- `msm_path`: Path to MSM binary
- `build_config`: BuildTools settings, temp cleanup rules, jar git repo
- `world_config`: World creation defaults (jargroup, version, world names)
- `op_usernames`: List of operator usernames to apply to new worlds

### Command Framework

Uses **Cobra** for CLI. All commands are registered in `cmd/` package:
- Root command initializes config via Viper
- Subcommands call into `internal/minecraft.AuMc` for business logic
- Config loading happens automatically in root command's `initConfig()`

### Business Logic Flow

**Build Process** (`AuMc.BuildNewJar`):
1. Clean temp directories/files from build_config
2. Remove old spigot*.jar files if configured
3. Execute BuildTools.jar with specified Minecraft version
4. Copy resulting spigot*.jar to git repo jars/ directory

**World Management** (partially implemented):
- Create: Calls MSM commands, generates eula.txt, updates server.properties, sets operators, adjusts file ownership
- Delete: Creates backup, copies to home directory, cleans MSM archives

## Migration Status

This project is migrating from Python to Go. See `MIGRATION_PLAN.md` for detailed status. Current state:
- Phase 1-2 (Setup, Config, MCConfig parser): ✅ Complete
- Phase 3 (Build, World management): Partially complete
- Phase 4-5 (CLI commands, Testing): In progress

## Development Workflow

### Working on Migration Tasks
When working on migration tasks, follow this workflow:

1. **Review the task**: Check `MIGRATION_PLAN.md` for the current phase/step to understand what needs to be done.

2. **Update main branch**:
   ```bash
   git checkout main
   git pull
   ```

3. **Create feature branch**: Create a branch from `main`, naming it based on the step being worked on (e.g., `phase3-step3.3-world-create`):
   ```bash
   git checkout -b <branch-name>
   ```

4. **Push empty branch and create draft PR**:
   ```bash
   git push -u origin <branch-name>
   # Then create a draft PR on GitHub
   ```

5. **Complete the work**: Implement the changes, test thoroughly, and push to the feature branch:
   ```bash
   git add .
   git commit -m "Descriptive commit message"
   git push origin <branch-name>
   ```

6. **After merge**: Once the PR is approved and merged to main:
   ```bash
   git checkout main
   git pull
   git branch -D <branch-name>
   ```

## Testing Approach

- **Unit tests**: Use Go's standard testing package (`*_test.go` files)
- **Integration tests**: Require Docker container for MSM operations
- No mocking framework currently in use - tests use real file operations

To run specific test:
```bash
go test -v ./internal/config       # Test specific package
go test -v -run TestLoadConfig ./internal/config   # Test specific function
```

## Code Style

### Go Conventions
- Standard library imports first, then third-party (cobra, viper), then local
- Unexported (private) functions/types: start with lowercase
- Error handling: return errors, don't panic (except in main/init)
- Use `fmt.Errorf` with `%w` for error wrapping

### Naming
- Structs: PascalCase (`AuMc`, `BuildConfig`, `MCConfig`)
- Methods/Functions: PascalCase for exported, camelCase for unexported
- Package names: lowercase, single word when possible

### Error Handling
- Custom errors defined in `internal/errors`
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check errors immediately after operations

## Important Project Context

### MSM Dependency
This tool wraps **MSM (Minecraft Server Manager)**, which must be installed on the target system. MSM handles:
- Server lifecycle (start/stop/restart)
- Jargroup management
- World backups
- Server file organization

Commands executed via `os/exec` package must account for MSM's expected paths and behavior.

### BuildTools Integration
Spigot jar builds require **BuildTools.jar** in the build directory. The build process:
- Must run in the build_config.build_directory (working directory matters)
- Executes: `java -jar BuildTools.jar --rev <version>`
- Generates significant temp artifacts that must be cleaned up

### Server Properties Management
Minecraft's `server.properties` uses simple `key=value` format with comments. The `mcprops` package:
- Preserves comments and formatting when possible
- Adds timestamp header on writes
- Handles property updates without full file rewrites

### Docker Test Environment
The Docker container (`aumc-test-server`) provides:
- MSM installed at `/opt/msm`
- BuildTools.jar at `/opt/build_tools/BuildTools.jar`
- Users: `testuser` (main, has sudo), `minecraft` (MSM service user)
- Workspace synced to `/workspace` (maps to project root)
- Minecraft server port 25565 forwarded to host
- Persistent data volume: `minecraft-data`

## CI/CD

GitHub Actions workflow (`.github/workflows/test.yml`):
- Runs on pushes to `main` and `python-to-go` branches
- Two jobs: **test** (runs tests, builds all platforms) and **lint** (golangci-lint)
- Uses Go 1.23
- Validates: dependencies, tests pass, builds succeed for all platforms

## Configuration Files

When running `aumc init`, two files are created:
- `config.json`: Main configuration (paths, versions, world settings)
- `server.properties.template`: Default Minecraft server properties

These use templates from `internal/config/templates.go` and must be placed where `AU_CONFIG_FILE` environment variable points to the config.json location.
