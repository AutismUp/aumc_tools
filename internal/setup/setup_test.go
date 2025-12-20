package setup

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGoModuleExists verifies that go.mod exists and has the correct module name
func TestGoModuleExists(t *testing.T) {
	// Get the project root (two levels up from internal/setup)
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	goModPath := filepath.Join(projectRoot, "go.mod")
	
	// Check if go.mod exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Fatal("go.mod does not exist")
	}

	// Read go.mod content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	// Verify module name
	expectedModule := "module github.com/AutismUp/aumc_tools"
	if !strings.Contains(string(content), expectedModule) {
		t.Errorf("go.mod does not contain expected module name: %s", expectedModule)
	}
}

// TestRequiredDirectoriesExist verifies that all required directories exist
func TestRequiredDirectoriesExist(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	requiredDirs := []string{
		"cmd/aumc",
		"internal/config",
		"internal/minecraft",
		"internal/mcprops",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(projectRoot, dir)
		if info, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Required directory does not exist: %s", dir)
		} else if !info.IsDir() {
			t.Errorf("Path exists but is not a directory: %s", dir)
		}
	}
}

// TestRequiredDependencies verifies that required dependencies are in go.mod
func TestRequiredDependencies(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	goModPath := filepath.Join(projectRoot, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	requiredDeps := []string{
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
	}

	for _, dep := range requiredDeps {
		if !strings.Contains(string(content), dep) {
			t.Errorf("Required dependency not found in go.mod: %s", dep)
		}
	}
}

// TestMakefileExists verifies that Makefile exists
func TestMakefileExists(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	makefilePath := filepath.Join(projectRoot, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Fatal("Makefile does not exist")
	}
}

// TestMakefileTargets verifies that required Makefile targets exist
func TestMakefileTargets(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	makefilePath := filepath.Join(projectRoot, "Makefile")
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}

	requiredTargets := []string{
		".PHONY: build",
		".PHONY: build-all",
		".PHONY: build-linux",
		".PHONY: build-darwin",
		".PHONY: build-windows",
		".PHONY: clean",
		".PHONY: test",
	}

	for _, target := range requiredTargets {
		if !strings.Contains(string(content), target) {
			t.Errorf("Required Makefile target not found: %s", target)
		}
	}
}

// TestGitignoreHasGoEntries verifies that .gitignore has Go-specific entries
func TestGitignoreHasGoEntries(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	gitignorePath := filepath.Join(projectRoot, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	requiredEntries := []string{
		"bin/",
		"*.exe",
	}

	for _, entry := range requiredEntries {
		if !strings.Contains(string(content), entry) {
			t.Errorf("Required .gitignore entry not found: %s", entry)
		}
	}
}

// TestGoModTidy verifies that go mod tidy works without errors
func TestGoModTidy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping go mod tidy test in short mode")
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod tidy failed: %v\nOutput: %s", err, string(output))
	}
}

// TestBuildProducesBinary verifies that make build produces a binary
func TestBuildProducesBinary(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping build test in short mode")
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Clean first
	cleanCmd := exec.Command("make", "clean")
	cleanCmd.Dir = projectRoot
	if output, err := cleanCmd.CombinedOutput(); err != nil {
		t.Logf("make clean output: %s", string(output))
	}

	// Build
	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make build failed: %v\nOutput: %s", err, string(output))
	}

	// Verify binary exists
	binaryPath := filepath.Join(projectRoot, "bin", "aumc")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Errorf("Binary was not created at %s", binaryPath)
	}

	// Verify binary is executable
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("Failed to stat binary: %v", err)
	}
	
	mode := info.Mode()
	if mode&0111 == 0 {
		t.Error("Binary is not executable")
	}
}

// TestCleanRemovesBuildArtifacts verifies that make clean removes build artifacts
func TestCleanRemovesBuildArtifacts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping clean test in short mode")
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Build first
	buildCmd := exec.Command("make", "build")
	buildCmd.Dir = projectRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("make build failed: %v\nOutput: %s", err, string(output))
	}

	// Verify binary exists
	binaryPath := filepath.Join(projectRoot, "bin", "aumc")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatal("Binary was not created before clean test")
	}

	// Clean
	cleanCmd := exec.Command("make", "clean")
	cleanCmd.Dir = projectRoot
	output, err := cleanCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make clean failed: %v\nOutput: %s", err, string(output))
	}

	// Verify bin directory is removed
	binPath := filepath.Join(projectRoot, "bin")
	if _, err := os.Stat(binPath); !os.IsNotExist(err) {
		t.Error("bin/ directory still exists after make clean")
	}
}

// TestCrossCompilationTargets verifies that cross-compilation targets work
func TestCrossCompilationTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cross-compilation test in short mode")
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Clean first
	cleanCmd := exec.Command("make", "clean")
	cleanCmd.Dir = projectRoot
	cleanCmd.Run()

	tests := []struct {
		target       string
		expectedFile string
	}{
		{"build-linux", "bin/aumc-linux-amd64"},
		{"build-darwin", "bin/aumc-darwin-amd64"},
		{"build-darwin", "bin/aumc-darwin-arm64"},
		{"build-windows", "bin/aumc-windows-amd64.exe"},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			cmd := exec.Command("make", tt.target)
			cmd.Dir = projectRoot
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("make %s failed: %v\nOutput: %s", tt.target, err, string(output))
			}

			binaryPath := filepath.Join(projectRoot, tt.expectedFile)
			if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
				t.Errorf("Expected binary not created: %s", tt.expectedFile)
			}
		})
	}

	// Clean up
	cleanCmd = exec.Command("make", "clean")
	cleanCmd.Dir = projectRoot
	cleanCmd.Run()
}

// TestMainPackageExists verifies that cmd/aumc/main.go exists
func TestMainPackageExists(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	mainPath := filepath.Join(projectRoot, "cmd", "aumc", "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		t.Fatal("cmd/aumc/main.go does not exist")
	}

	// Verify it's a valid Go file
	content, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	if !strings.Contains(string(content), "package main") {
		t.Error("main.go does not contain 'package main'")
	}

	if !strings.Contains(string(content), "func main()") {
		t.Error("main.go does not contain 'func main()'")
	}
}

// TestRootCommandExists verifies that cmd/root.go exists and has basic structure
func TestRootCommandExists(t *testing.T) {
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	rootPath := filepath.Join(projectRoot, "cmd", "root.go")
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		t.Fatal("cmd/root.go does not exist")
	}

	content, err := os.ReadFile(rootPath)
	if err != nil {
		t.Fatalf("Failed to read root.go: %v", err)
	}

	requiredElements := []string{
		"package cmd",
		"var rootCmd",
		"func Execute()",
		"cobra.Command",
	}

	for _, element := range requiredElements {
		if !strings.Contains(string(content), element) {
			t.Errorf("root.go does not contain required element: %s", element)
		}
	}
}
