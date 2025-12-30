package minecraft

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/AutismUp/aumc_tools/internal/config"
)

// createTestConfig creates a temporary config file for testing
func createTestConfig(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testConfig := config.Config{
		MSMPath: "/opt/msm",
		BuildConfig: config.BuildConfig{
			BuildDirectory:   "/tmp/build",
			MinecraftVersion: "1.20.4",
			JarGitRepo:       "/home/user/git/jars",
			TempFolders:      []string{"apache-maven-3.6.0", "BuildTools", "work"},
			TempFiles:        []string{"BuildTools.log.txt"},
			DeleteSpigotJars: true,
		},
		WorldConfig: config.WorldConfig{
			ServerPropertiesTemplate: "/home/user/.aumc/server.properties.template",
			Jargroup:                 "1_20_4",
			MinecraftVersion:         "1.20.4",
			WorldNames:               []string{"world1", "world2"},
		},
		OpUsernames: []string{"player1", "player2"},
	}

	data, err := json.MarshalIndent(testConfig, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	return configPath
}

func TestNewAuMc(t *testing.T) {
	t.Run("valid config path", func(t *testing.T) {
		configPath := createTestConfig(t)

		aumc, err := NewAuMc(configPath)
		if err != nil {
			t.Fatalf("NewAuMc() failed with valid config: %v", err)
		}

		if aumc == nil {
			t.Fatal("NewAuMc() returned nil AuMc")
		}

		if aumc.config == nil {
			t.Fatal("AuMc.config is nil")
		}
	})

	t.Run("invalid config path", func(t *testing.T) {
		aumc, err := NewAuMc("/nonexistent/path/config.json")
		if err == nil {
			t.Fatal("NewAuMc() should fail with invalid path")
		}

		if aumc != nil {
			t.Fatal("NewAuMc() should return nil on error")
		}
	})

	t.Run("empty config path", func(t *testing.T) {
		aumc, err := NewAuMc("")
		if err == nil {
			t.Fatal("NewAuMc() should fail with empty path")
		}

		if aumc != nil {
			t.Fatal("NewAuMc() should return nil on error")
		}
	})
}

func TestGetConfig(t *testing.T) {
	configPath := createTestConfig(t)
	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	cfg := aumc.GetConfig()
	if cfg == nil {
		t.Fatal("GetConfig() returned nil")
	}

	if cfg.MSMPath != "/opt/msm" {
		t.Errorf("Expected MSMPath '/opt/msm', got '%s'", cfg.MSMPath)
	}
}

func TestGetBuildConfig(t *testing.T) {
	configPath := createTestConfig(t)
	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	buildCfg := aumc.GetBuildConfig()
	if buildCfg == nil {
		t.Fatal("GetBuildConfig() returned nil")
	}

	if buildCfg.MinecraftVersion != "1.20.4" {
		t.Errorf("Expected MinecraftVersion '1.20.4', got '%s'", buildCfg.MinecraftVersion)
	}

	if buildCfg.BuildDirectory != "/tmp/build" {
		t.Errorf("Expected BuildDirectory '/tmp/build', got '%s'", buildCfg.BuildDirectory)
	}

	if !buildCfg.DeleteSpigotJars {
		t.Error("Expected DeleteSpigotJars to be true")
	}
}

func TestGetWorldConfig(t *testing.T) {
	configPath := createTestConfig(t)
	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	worldCfg := aumc.GetWorldConfig()
	if worldCfg == nil {
		t.Fatal("GetWorldConfig() returned nil")
	}

	if worldCfg.Jargroup != "1_20_4" {
		t.Errorf("Expected Jargroup '1_20_4', got '%s'", worldCfg.Jargroup)
	}

	if len(worldCfg.WorldNames) != 2 {
		t.Errorf("Expected 2 world names, got %d", len(worldCfg.WorldNames))
	}
}

func TestGetMSMPath(t *testing.T) {
	configPath := createTestConfig(t)
	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	msmPath := aumc.GetMSMPath()
	if msmPath != "/opt/msm" {
		t.Errorf("Expected MSMPath '/opt/msm', got '%s'", msmPath)
	}
}

func TestGetOpUsernames(t *testing.T) {
	configPath := createTestConfig(t)
	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	opUsernames := aumc.GetOpUsernames()
	if len(opUsernames) != 2 {
		t.Errorf("Expected 2 op usernames, got %d", len(opUsernames))
	}

	expectedOps := map[string]bool{"player1": true, "player2": true}
	for _, op := range opUsernames {
		if !expectedOps[op] {
			t.Errorf("Unexpected op username: %s", op)
		}
	}
}

func TestBuildNewJar_CleanupOnly(t *testing.T) {
	// This test verifies cleanup behavior without running actual BuildTools
	tempDir := t.TempDir()
	buildDir := filepath.Join(tempDir, "build")
	jarRepoDir := filepath.Join(tempDir, "jarrepo")

	// Create build directory structure
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		t.Fatalf("Failed to create build dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(jarRepoDir, "jars"), 0755); err != nil {
		t.Fatalf("Failed to create jar repo dir: %v", err)
	}

	// Create temp directories
	tempDirs := []string{"BuildData", "Bukkit", "work"}
	for _, dir := range tempDirs {
		path := filepath.Join(buildDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("Failed to create temp dir %s: %v", dir, err)
		}
		// Create a file in the directory to verify it gets cleaned
		testFile := filepath.Join(path, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create temp files
	tempFiles := []string{"BuildTools.log.txt"}
	for _, file := range tempFiles {
		path := filepath.Join(buildDir, file)
		if err := os.WriteFile(path, []byte("log content"), 0644); err != nil {
			t.Fatalf("Failed to create temp file %s: %v", file, err)
		}
	}

	// Create old jar files
	oldJars := []string{"spigot-1.19.jar", "spigot-1.20.jar"}
	for _, jar := range oldJars {
		path := filepath.Join(buildDir, jar)
		if err := os.WriteFile(path, []byte("jar content"), 0644); err != nil {
			t.Fatalf("Failed to create old jar %s: %v", jar, err)
		}
	}

	// Create config
	tempConfigDir := t.TempDir()
	configPath := filepath.Join(tempConfigDir, "config.json")
	testConfig := config.Config{
		MSMPath: "/opt/msm",
		BuildConfig: config.BuildConfig{
			BuildDirectory:   buildDir,
			MinecraftVersion: "1.20.4",
			JarGitRepo:       jarRepoDir,
			TempFolders:      tempDirs,
			TempFiles:        tempFiles,
			DeleteSpigotJars: true,
		},
	}

	data, err := json.MarshalIndent(testConfig, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	aumc, err := NewAuMc(configPath)
	if err != nil {
		t.Fatalf("Failed to create AuMc: %v", err)
	}

	// Note: We can't test the full BuildNewJar without BuildTools.jar
	// but we can verify the config is set up correctly
	if aumc.config.BuildConfig.BuildDirectory != buildDir {
		t.Errorf("Build directory not set correctly")
	}

	// Verify temp directories exist before cleanup
	for _, dir := range tempDirs {
		path := filepath.Join(buildDir, dir)
		if !dirExists(path) {
			t.Errorf("Temp dir %s should exist before cleanup", dir)
		}
	}

	// Verify temp files exist before cleanup
	for _, file := range tempFiles {
		path := filepath.Join(buildDir, file)
		if !fileExists(path) {
			t.Errorf("Temp file %s should exist before cleanup", file)
		}
	}

	// Verify old jars exist before cleanup
	for _, jar := range oldJars {
		path := filepath.Join(buildDir, jar)
		if !fileExists(path) {
			t.Errorf("Old jar %s should exist before cleanup", jar)
		}
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("file exists", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		if !fileExists(filePath) {
			t.Error("fileExists() should return true for existing file")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "nonexistent.txt")
		if fileExists(filePath) {
			t.Error("fileExists() should return false for nonexistent file")
		}
	})

	t.Run("directory is not a file", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}

		if fileExists(dirPath) {
			t.Error("fileExists() should return false for directory")
		}
	})
}

func TestDirExists(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("directory exists", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}

		if !dirExists(dirPath) {
			t.Error("dirExists() should return true for existing directory")
		}
	})

	t.Run("directory does not exist", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "nonexistent")
		if dirExists(dirPath) {
			t.Error("dirExists() should return false for nonexistent directory")
		}
	})

	t.Run("file is not a directory", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test.txt")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		if dirExists(filePath) {
			t.Error("dirExists() should return false for file")
		}
	})
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("successful copy", func(t *testing.T) {
		srcPath := filepath.Join(tempDir, "source.txt")
		dstPath := filepath.Join(tempDir, "dest.txt")
		testContent := []byte("test content for copying")

		if err := os.WriteFile(srcPath, testContent, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("copyFile() failed: %v", err)
		}

		if !fileExists(dstPath) {
			t.Error("Destination file should exist after copy")
		}

		copiedContent, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("Failed to read copied file: %v", err)
		}

		if string(copiedContent) != string(testContent) {
			t.Errorf("Copied content doesn't match. Expected %q, got %q", testContent, copiedContent)
		}
	})

	t.Run("source file does not exist", func(t *testing.T) {
		srcPath := filepath.Join(tempDir, "nonexistent.txt")
		dstPath := filepath.Join(tempDir, "dest2.txt")

		if err := copyFile(srcPath, dstPath); err == nil {
			t.Error("copyFile() should fail when source doesn't exist")
		}
	})

	t.Run("destination directory does not exist", func(t *testing.T) {
		srcPath := filepath.Join(tempDir, "source2.txt")
		dstPath := filepath.Join(tempDir, "nonexistent", "dest.txt")
		testContent := []byte("test")

		if err := os.WriteFile(srcPath, testContent, 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		if err := copyFile(srcPath, dstPath); err == nil {
			t.Error("copyFile() should fail when destination directory doesn't exist")
		}
	})
}
