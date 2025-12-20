# Python to Go Migration Plan

## Branch
Working branch: `python-to-go`

## Recent Updates
**December 1, 2025**: Merged `origin/main` into `python-to-go` branch
- Resolved merge conflicts in `go.mod`, `go.sum`, and `internal/config/config.go`
- Kept newer dependency versions (Go 1.25.4, Cobra v1.10.1, Viper v1.21.0)
- Added `SaveDefaultFiles()` function for CLI compatibility
- Fixed `cmd/aumc/main.go` to properly call root command
- All tests passing, binary builds and runs successfully

## Migration Strategy
Bottom-up conversion with incremental testing. Each phase should be completed and tested before moving to the next.

---

## Phase 1: Project Setup & Foundation

### Step 1.1: Initialize Go Module ✅
- [x] Create `go.mod` with module name `github.com/AutismUp/aumc_tools`
- [x] Create directory structure:
  ```
  cmd/aumc/          # CLI entry point
  internal/
    config/          # Configuration management (Viper)
    minecraft/       # Core business logic
    mcprops/         # server.properties parser
  ```
- [x] Add dependencies: `cobra`, `viper`

### Step 1.2: Setup Build Tooling ✅
- [x] Create `Makefile` with targets: `build`, `build-all` (cross-platform), `clean`
- [x] Add `.gitignore` entries for Go (`bin/`, `*.exe`, etc.)
- [x] Test: `go mod init` and `go mod tidy` work

---

## Phase 2: Core Infrastructure

### Step 2.1: Configuration Management (Viper) ✅
- [x] Port `config_templates.py` default templates as Go constants
- [x] Create Go structs for config.json structure:
  - `Config` (root)
  - `BuildConfig`
  - `WorldConfig`
- [x] Implement config loading with Viper (env var `AU_CONFIG_FILE`)
- [x] Implement config initialization prompt (create default files)
- [x] Test: Load sample config.json, verify struct marshaling

### Step 2.2: MCConfig Parser (server.properties) ✅
- [x] Port `MCConfig` class to Go struct/methods
- [x] Implement `LoadProperties(filepath)` - parse key=value format, skip comments
- [x] Implement `UpdateProperty(key, value)`
- [x] Implement `WriteProperties(filepath)` - write with timestamp header
- [x] Test: Read/modify/write server.properties file

### Step 2.3: Error Types ✅
- [x] Define custom error types as Go errors
- [x] Create error constructors/helpers
- [x] Test: Error creation and formatting

---

## Phase 3: Business Logic Layer

### Step 3.1: Core AuMc Struct
- [ ] Create `AuMc` struct with config field
- [ ] Implement `NewAuMc(configPath)` constructor
- [ ] Test: Initialize with valid config

### Step 3.2: Build New Jar
- [ ] Port `build_new_jar()` method
- [ ] Handle directory cleanup (temp folders, temp files, old jars)
- [ ] Execute Java subprocess for BuildTools
- [ ] Copy jar files to git repo
- [ ] Test: Dry run without actual BuildTools execution

### Step 3.3: World Management - Create
- [ ] Port `create_new_world(name, jargroup, version)` method
- [ ] Execute MSM commands via subprocess
- [ ] Generate eula.txt file
- [ ] Update server.properties via MCConfig
- [ ] Handle operator additions
- [ ] Set file ownership (chown/chgrp)
- [ ] Test: Create test world in Vagrant VM

### Step 3.4: World Management - Delete
- [ ] Port `delete_world(name)` method
- [ ] Backup world before deletion
- [ ] Copy latest backup to home directory
- [ ] Clean up MSM archives
- [ ] Test: Delete test world in Vagrant VM

### Step 3.5: Jar Publishing
- [ ] Port `publish_new_jar()` method
- [ ] Git operations (add, commit, push)
- [ ] MSM jargroup creation
- [ ] Test: Mock git operations

### Step 3.6: World Restoration (Optional)
- [ ] Port `restore_world(world_name, restore_to_date)` method
- [ ] Note: Not currently exposed in CLI
- [ ] Test: Manual testing if time permits

---

## Phase 4: CLI Layer (Cobra)

### Step 4.1: Root Command & Bootstrap
- [ ] Initialize Cobra app in `cmd/aumc/main.go`
- [ ] Implement root command with config file detection
- [ ] Implement config initialization prompt (AU_CONFIG_FILE not set)
- [ ] Generate default config.json and server.properties.template
- [ ] Test: Run binary without config, verify prompt

### Step 4.2: Config Commands
- [ ] Implement `aumc check-config` - print current configuration
- [ ] Implement `aumc reload-config` - reload from file
- [ ] Test: Verify config display and reload

### Step 4.3: Jar Commands
- [ ] Implement `aumc build-new-jar` - build Spigot jar
- [ ] Implement `aumc publish-new-jar --filename <name>` - publish to GitHub
- [ ] Test: Build command execution (may need mock mode)

### Step 4.4: World Commands
- [ ] Implement `aumc create-new-world --name <name> --jargroup <group> --version <ver>`
- [ ] Implement `aumc create-new-world --from-config` - batch create
- [ ] Implement `aumc delete-world --name <name>`
- [ ] Implement `aumc delete-world --from-config` - batch delete
- [ ] Test: Full world lifecycle in Vagrant VM

---

## Phase 5: Testing & Distribution

### Step 5.1: Integration Testing
- [ ] Test all commands against Vagrant VM
- [ ] Verify parity with Python version behavior
- [ ] Test cross-platform builds (Linux, macOS at minimum)

### Step 5.2: Build & Release
- [ ] Create release builds for: Linux (amd64), macOS (amd64, arm64), Windows (amd64)
- [ ] Document build process in README
- [ ] Create installation instructions for binary distribution

### Step 5.3: Documentation
- [ ] Update README.md with Go installation instructions
- [ ] Update CRUSH.md with Go build/test commands
- [ ] Add migration notes (deprecation timeline for Python version)
- [ ] Document binary usage vs pip install

---

## Key Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- Standard library: `os/exec`, `encoding/json`, `path/filepath`, `os`, `io`

## Testing Approach
- Unit tests for config parsing, property file handling
- Integration tests via Vagrant VM for MSM operations
- Manual smoke tests for each command before phase completion
- Keep Python version available for comparison testing

## Success Criteria
- Single binary runs on Linux, macOS, Windows
- All existing CLI commands work identically
- No external dependencies (Python, pip) required
- Binary size < 20MB
