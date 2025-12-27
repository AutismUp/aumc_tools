package minecraft

import (
	"github.com/AutismUp/aumc_tools/internal/config"
)

// AuMc is the main business logic handler for Autism Up Minecraft server management
type AuMc struct {
	config *config.Config
}

// NewAuMc creates a new AuMc instance by loading configuration from the specified path
func NewAuMc(configPath string) (*AuMc, error) {
	cfg, err := config.LoadFromPath(configPath)
	if err != nil {
		return nil, err
	}

	return &AuMc{
		config: cfg,
	}, nil
}

// GetConfig returns the full configuration
func (a *AuMc) GetConfig() *config.Config {
	return a.config
}

// GetBuildConfig returns the build configuration section
func (a *AuMc) GetBuildConfig() *config.BuildConfig {
	return &a.config.BuildConfig
}

// GetWorldConfig returns the world configuration section
func (a *AuMc) GetWorldConfig() *config.WorldConfig {
	return &a.config.WorldConfig
}

// GetMSMPath returns the MSM path from configuration
func (a *AuMc) GetMSMPath() string {
	return a.config.MSMPath
}

// GetOpUsernames returns the list of operator usernames
func (a *AuMc) GetOpUsernames() []string {
	return a.config.OpUsernames
}
