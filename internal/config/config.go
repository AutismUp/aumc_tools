package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// BuildConfig represents the build configuration section
type BuildConfig struct {
	TempFolders      []string `json:"temp_folders" mapstructure:"temp_folders"`
	TempFiles        []string `json:"temp_files" mapstructure:"temp_files"`
	BuildDirectory   string   `json:"build_directory" mapstructure:"build_directory"`
	MinecraftVersion string   `json:"minecraft_version" mapstructure:"minecraft_version"`
	JarGitRepo       string   `json:"jar_git_repo" mapstructure:"jar_git_repo"`
	DeleteSpigotJars bool     `json:"delete_spigot_jars" mapstructure:"delete_spigot_jars"`
}

// WorldConfig represents the world configuration section
type WorldConfig struct {
	WorldNames               []string `json:"world_names" mapstructure:"world_names"`
	ServerPropertiesTemplate string   `json:"server_properties_template" mapstructure:"server_properties_template"`
}

// Config represents the root configuration structure
type Config struct {
	MSMPath     string      `json:"msm_path" mapstructure:"msm_path"`
	BuildConfig BuildConfig `json:"build_config" mapstructure:"build_config"`
	WorldConfig WorldConfig `json:"world_config" mapstructure:"world_config"`
	OpUsernames []string    `json:"op_usernames" mapstructure:"op_usernames"`
}

// Load reads and parses the configuration file using Viper
// It respects the AU_CONFIG_FILE environment variable
func Load() (*Config, error) {
	configPath := os.Getenv("AU_CONFIG_FILE")
	if configPath == "" {
		return nil, fmt.Errorf("AU_CONFIG_FILE environment variable not set")
	}

	return LoadFromPath(configPath)
}

// LoadFromPath reads and parses the configuration file from a specific path
func LoadFromPath(configPath string) (*Config, error) {
	// Expand home directory if present
	if len(configPath) >= 2 && configPath[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(home, configPath[2:])
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Configure Viper
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// InitializeConfig creates default configuration files if they don't exist
// Returns the path where config.json was created
func InitializeConfig(configDir string) (string, error) {
	// Expand home directory if present
	if len(configDir) >= 2 && configDir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(home, configDir[2:])
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create config.json
	configPath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Parse default config to ensure it's valid JSON
		var defaultCfg Config
		if err := json.Unmarshal([]byte(DefaultConfig), &defaultCfg); err != nil {
			return "", fmt.Errorf("invalid default config template: %w", err)
		}

		// Write formatted JSON
		prettyJSON, err := json.MarshalIndent(defaultCfg, "", "    ")
		if err != nil {
			return "", fmt.Errorf("failed to format default config: %w", err)
		}

		if err := os.WriteFile(configPath, prettyJSON, 0644); err != nil {
			return "", fmt.Errorf("failed to write config.json: %w", err)
		}
	}

	// Create server.properties.template
	templatePath := filepath.Join(configDir, "server.properties.template")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		if err := os.WriteFile(templatePath, []byte(DefaultServerProperties), 0644); err != nil {
			return "", fmt.Errorf("failed to write server.properties.template: %w", err)
		}
	}

	return configPath, nil
}

// Validate checks if the configuration has all required fields populated
func (c *Config) Validate() error {
	if c.MSMPath == "" {
		return fmt.Errorf("msm_path is required")
	}

	if c.BuildConfig.BuildDirectory == "" {
		return fmt.Errorf("build_config.build_directory is required")
	}

	if c.BuildConfig.MinecraftVersion == "" {
		return fmt.Errorf("build_config.minecraft_version is required")
	}

	if c.BuildConfig.JarGitRepo == "" {
		return fmt.Errorf("build_config.jar_git_repo is required")
	}

	if c.WorldConfig.ServerPropertiesTemplate == "" {
		return fmt.Errorf("world_config.server_properties_template is required")
	}

	return nil
}

// SaveDefaultFiles writes default `config.json` and `server.properties.template` to the current directory.
// It overwrites existing files.
func SaveDefaultFiles() error {
	var defaultCfg Config
	if err := json.Unmarshal([]byte(DefaultConfig), &defaultCfg); err != nil {
		return fmt.Errorf("invalid default config template: %w", err)
	}

	prettyJSON, err := json.MarshalIndent(defaultCfg, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to format default config: %w", err)
	}

	if err := os.WriteFile("config.json", prettyJSON, 0644); err != nil {
		return fmt.Errorf("failed to write config.json: %w", err)
	}

	if err := os.WriteFile("server.properties.template", []byte(DefaultServerProperties), 0644); err != nil {
		return fmt.Errorf("failed to write server.properties.template: %w", err)
	}

	return nil
}
