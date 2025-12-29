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
