package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "aumc-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testConfig := Config{
		MSMPath: "/opt/msm",
		BuildConfig: BuildConfig{
			BuildDirectory:   "/home/user/build",
			TempFolders:      []string{"BuildData", "Bukkit"},
			TempFiles:        []string{"BuildTools.log.txt"},
			MinecraftVersion: "1.20.1",
			DeleteSpigotJars: true,
			JarGitRepo:       "/home/user/jars",
		},
		WorldConfig: WorldConfig{
			ServerPropertiesTemplate: "/home/user/server.properties.template",
			WorldNames:               []string{"world1", "world2"},
		},
		OpUsernames: []string{"admin1", "admin2"},
	}

	configPath := filepath.Join(tmpDir, "config.json")
	configJSON, err := json.MarshalIndent(testConfig, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.MSMPath != testConfig.MSMPath {
		t.Errorf("Expected MSMPath %s, got %s", testConfig.MSMPath, cfg.MSMPath)
	}

	if cfg.BuildConfig.MinecraftVersion != testConfig.BuildConfig.MinecraftVersion {
		t.Errorf("Expected MinecraftVersion %s, got %s", testConfig.BuildConfig.MinecraftVersion, cfg.BuildConfig.MinecraftVersion)
	}

	if len(cfg.WorldConfig.WorldNames) != len(testConfig.WorldConfig.WorldNames) {
		t.Errorf("Expected %d world names, got %d", len(testConfig.WorldConfig.WorldNames), len(cfg.WorldConfig.WorldNames))
	}
}

func TestInitializeConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "aumc-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath, err := InitializeConfig(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config.json was not created at %s", configPath)
	}

	templatePath := filepath.Join(tmpDir, "server.properties.template")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Errorf("server.properties.template was not created at %s", templatePath)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("Failed to load initialized config: %v", err)
	}

	if len(cfg.BuildConfig.TempFolders) == 0 {
		t.Error("Expected default temp_folders to be populated")
	}

	if len(cfg.OpUsernames) == 0 {
		t.Error("Expected default op_usernames to be populated")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				MSMPath: "/opt/msm",
				BuildConfig: BuildConfig{
					BuildDirectory:   "/home/user/build",
					MinecraftVersion: "1.20.1",
					JarGitRepo:       "/home/user/jars",
				},
				WorldConfig: WorldConfig{
					ServerPropertiesTemplate: "/home/user/template",
				},
			},
			wantErr: false,
		},
		{
			name: "missing msm_path",
			config: Config{
				BuildConfig: BuildConfig{
					BuildDirectory:   "/home/user/build",
					MinecraftVersion: "1.20.1",
					JarGitRepo:       "/home/user/jars",
				},
				WorldConfig: WorldConfig{
					ServerPropertiesTemplate: "/home/user/template",
				},
			},
			wantErr: true,
			errMsg:  "msm_path is required",
		},
		{
			name: "missing build_directory",
			config: Config{
				MSMPath: "/opt/msm",
				BuildConfig: BuildConfig{
					MinecraftVersion: "1.20.1",
					JarGitRepo:       "/home/user/jars",
				},
				WorldConfig: WorldConfig{
					ServerPropertiesTemplate: "/home/user/template",
				},
			},
			wantErr: true,
			errMsg:  "build_config.build_directory is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDefaultConfigIsValidJSON(t *testing.T) {
	var cfg Config
	if err := json.Unmarshal([]byte(DefaultConfig), &cfg); err != nil {
		t.Fatalf("DefaultConfig is not valid JSON: %v", err)
	}

	if len(cfg.BuildConfig.TempFolders) == 0 {
		t.Error("Expected default temp_folders to be populated")
	}

	if len(cfg.OpUsernames) == 0 {
		t.Error("Expected default op_usernames to be populated")
	}
}
