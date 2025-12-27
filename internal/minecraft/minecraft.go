// Package minecraft provides business logic for Minecraft server management.
// Use internal/mcprops for server.properties file handling.
package minecraft

import (
	"fmt"
	"os"
	"path/filepath"
)

// ServerManager handles Minecraft server operations
type ServerManager struct {
	MSMPath string
}

// NewServerManager creates a new ServerManager instance
func NewServerManager(msmPath string) *ServerManager {
	return &ServerManager{
		MSMPath: msmPath,
	}
}

// WorldInfo represents information about a Minecraft world
type WorldInfo struct {
	Name             string
	JarGroup         string
	MinecraftVersion string
	ServerProperties string
	OpUsernames      []string
}

// CreateWorld creates a new Minecraft world
func (sm *ServerManager) CreateWorld(info *WorldInfo) error {
	// TODO: Implement world creation using MSM
	fmt.Printf("Creating world: %s with jargroup: %s, version: %s\n",
		info.Name, info.JarGroup, info.MinecraftVersion)
	return fmt.Errorf("world creation not yet implemented")
}

// DeleteWorld deletes a Minecraft world after creating a backup
func (sm *ServerManager) DeleteWorld(worldName string) error {
	// TODO: Implement world deletion with backup
	fmt.Printf("Deleting world: %s (with backup)\n", worldName)
	return fmt.Errorf("world deletion not yet implemented")
}

// RestoreWorld restores a world from backup
func (sm *ServerManager) RestoreWorld(worldName, backupPath string) error {
	// TODO: Implement world restoration
	fmt.Printf("Restoring world: %s from backup: %s\n", worldName, backupPath)
	return fmt.Errorf("world restoration not yet implemented")
}

// BuildJar builds a Spigot jar using BuildTools
func (sm *ServerManager) BuildJar(version, tempFolder, buildToolsJar, outputFolder string) error {
	// TODO: Implement jar building
	fmt.Printf("Building Spigot jar for version: %s\n", version)
	fmt.Printf("  Temp folder: %s\n", tempFolder)
	fmt.Printf("  BuildTools jar: %s\n", buildToolsJar)
	fmt.Printf("  Output folder: %s\n", outputFolder)
	return fmt.Errorf("jar building not yet implemented")
}

// DirectoryManager provides utilities for managing directories
type DirectoryManager struct {
	originalDir string
}

// NewDirectoryManager creates a new DirectoryManager
func NewDirectoryManager() *DirectoryManager {
	cwd, err := os.Getwd()
	if err != nil {
		// If we can't get current directory, use empty string
		fmt.Fprintf(os.Stderr, "Warning: failed to get current directory: %v\n", err)
		cwd = ""
	}
	return &DirectoryManager{originalDir: cwd}
}

// EnterDir changes to a directory and returns a function to restore the original directory
func (dm *DirectoryManager) EnterDir(path string) (func(), error) {
	expandedPath := filepath.Clean(os.ExpandEnv(path))

	if err := os.Chdir(expandedPath); err != nil {
		return nil, fmt.Errorf("failed to change directory to %s: %w", expandedPath, err)
	}

	return func() {
		if err := os.Chdir(dm.originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore directory: %v\n", err)
		}
	}, nil
}
