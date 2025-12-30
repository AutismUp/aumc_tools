package minecraft

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/AutismUp/aumc_tools/internal/config"
	"github.com/AutismUp/aumc_tools/internal/mcprops"
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

// BuildNewJar cleans up previous jar build content and builds the desired version of Spigot
func (a *AuMc) BuildNewJar() error {
	buildConfig := &a.config.BuildConfig

	// Step 1: Remove temporary directories
	fmt.Printf("Removing temporary directories from %s\n", buildConfig.BuildDirectory)
	for _, dir := range buildConfig.TempFolders {
		dirPath := filepath.Join(buildConfig.BuildDirectory, dir)
		if dirExists(dirPath) {
			fmt.Printf("  Removing: %s\n", dirPath)
			if err := os.RemoveAll(dirPath); err != nil {
				fmt.Printf("  Error: %s, %v\n", dirPath, err)
				// Continue with other cleanup even if one fails
			}
		}
	}

	// Step 2: Remove temporary files
	fmt.Printf("Removing temporary files from %s\n", buildConfig.BuildDirectory)
	for _, file := range buildConfig.TempFiles {
		filePath := filepath.Join(buildConfig.BuildDirectory, file)
		if fileExists(filePath) {
			fmt.Printf("  Removing: %s\n", filePath)
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("  Error: %s, %v\n", filePath, err)
				// Continue with other cleanup even if one fails
			}
		}
	}

	// Step 3: Remove old spigot jars if configured
	if buildConfig.DeleteSpigotJars {
		fmt.Printf("Removing old spigot jars from %s\n", buildConfig.BuildDirectory)
		pattern := filepath.Join(buildConfig.BuildDirectory, "spigot*.jar")
		oldJars, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("failed to glob for old jars: %w", err)
		}

		for _, jarFile := range oldJars {
			fmt.Printf("  Removing: %s\n", jarFile)
			if err := os.Remove(jarFile); err != nil {
				fmt.Printf("  Error: %s, %v\n", jarFile, err)
				// Continue with other cleanup even if one fails
			}
		}
	}

	// Step 4: Run BuildTools
	fmt.Println("Running BuildTools")

	// Save current directory and change to build directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Printf("Warning: failed to restore original directory: %v\n", err)
		}
	}()

	if err := os.Chdir(buildConfig.BuildDirectory); err != nil {
		return fmt.Errorf("failed to change to build directory: %w", err)
	}

	// Execute BuildTools
	cmd := exec.Command("java", "-jar", "BuildTools.jar", "--rev", buildConfig.MinecraftVersion)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("BuildTools execution failed: %w", err)
	}

	// Step 5: Copy new jar files to git repo
	pattern := filepath.Join(buildConfig.BuildDirectory, "spigot*.jar")
	newJars, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob for new jars: %w", err)
	}

	if len(newJars) == 0 {
		return fmt.Errorf("no spigot jars found after build")
	}

	// Ensure destination directory exists
	destDir := filepath.Join(buildConfig.JarGitRepo, "jars")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create jars directory: %w", err)
	}

	for _, jarFile := range newJars {
		destPath := filepath.Join(destDir, filepath.Base(jarFile))
		fmt.Printf("Copying %s to %s\n", jarFile, destPath)
		if err := copyFile(jarFile, destPath); err != nil {
			return fmt.Errorf("failed to copy jar file: %w", err)
		}
	}

	fmt.Println("BuildTools complete")
	return nil
}

// CreateNewWorld creates a new Minecraft world with the specified name, jargroup, and version
func (a *AuMc) CreateNewWorld(name, jargroup, version string) error {
	msmServerPath := filepath.Join(a.config.MSMPath, "servers")
	timeStamp := time.Now().Format("Mon Jan 02 15:04:05 MST 2006")

	// Step 1: Create the world using MSM
	fmt.Printf("Creating world '%s'...\n", name)
	cmd := exec.Command("sudo", "msm", "server", "create", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create world: %w", err)
	}

	// Step 2: Set the jargroup for the world
	fmt.Printf("Setting jargroup '%s' for world '%s'...\n", jargroup, name)
	cmd = exec.Command("sudo", "msm", name, "jar", jargroup)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set jargroup: %w", err)
	}

	// Step 3: Create the eula.txt file
	fmt.Println("Creating eula.txt...")
	eulaPath := filepath.Join(msmServerPath, name, "eula.txt")
	eulaContent := fmt.Sprintf("#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).\n#%s\neula=true\n", timeStamp)
	if err := os.WriteFile(eulaPath, []byte(eulaContent), 0644); err != nil {
		return fmt.Errorf("failed to create eula.txt: %w", err)
	}

	// Step 4: Update server.properties template and copy to the server folder
	fmt.Println("Configuring server.properties...")
	serverProps, err := mcprops.LoadProperties(a.config.WorldConfig.ServerPropertiesTemplate)
	if err != nil {
		return fmt.Errorf("failed to load server.properties template: %w", err)
	}

	serverProps.UpdateProperty("msm-version", fmt.Sprintf("minecraft/%s", version))
	serverProps.UpdateProperty("motd", fmt.Sprintf("Autism Up - %s", name))

	serverPropsPath := filepath.Join(msmServerPath, name, "server.properties")
	if err := serverProps.WriteProperties(serverPropsPath); err != nil {
		return fmt.Errorf("failed to write server.properties: %w", err)
	}

	// Step 5: Start the server
	fmt.Println("Starting server...")
	cmd = exec.Command("sudo", "msm", name, "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Step 6: Add operators
	for _, operator := range a.config.OpUsernames {
		fmt.Printf("Adding operator '%s'...\n", operator)
		cmd = exec.Command("sudo", "msm", name, "op", "add", operator)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Warning: failed to add operator %s: %v\n", operator, err)
			// Continue with other operators even if one fails
		}
	}

	// Step 7: Stop the server
	fmt.Println("Stopping server...")
	cmd = exec.Command("sudo", "msm", name, "stop", "now")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	// Step 8: Configure world RAM settings
	fmt.Println("Configuring world RAM settings...")
	cmd = exec.Command("sudo", "msm", name, "worlds", "ram", "world")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to configure world RAM: %v\n", err)
		// Non-critical, continue
	}

	// Step 9: Set file ownership to minecraft user
	fmt.Println("Setting file ownership...")
	worldPath := filepath.Join(msmServerPath, name)
	cmd = exec.Command("sudo", "chown", "-R", "minecraft", worldPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to chown: %v\n", err)
	}

	cmd = exec.Command("sudo", "chgrp", "-R", "minecraft", worldPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to chgrp: %v\n", err)
	}

	fmt.Printf("World named \"%s\" created\n", name)
	return nil
}

// DeleteWorld backs up a world, copies the backup to the home directory, and deletes it from MSM
func (a *AuMc) DeleteWorld(name string) error {
	msmArchivePath := filepath.Join(a.config.MSMPath, "archives")
	backupPath := filepath.Join(msmArchivePath, "backups", name)

	// Step 1: Backup the world before deletion
	fmt.Printf("Backing up world '%s'...\n", name)
	cmd := exec.Command("sudo", "msm", name, "backup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to backup world: %w", err)
	}

	// Step 2: Find the latest backup file
	fmt.Println("Finding latest backup file...")
	backupFiles, err := filepath.Glob(filepath.Join(backupPath, "*"))
	if err != nil {
		return fmt.Errorf("failed to find backup files: %w", err)
	}

	if len(backupFiles) == 0 {
		return fmt.Errorf("no backup files found in %s", backupPath)
	}

	// Find the most recently created backup file
	var latestFile string
	var latestTime time.Time
	for _, file := range backupFiles {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latestFile = file
		}
	}

	if latestFile == "" {
		return fmt.Errorf("could not determine latest backup file")
	}

	// Step 3: Copy the latest backup to the home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	destPath := filepath.Join(homedir, filepath.Base(latestFile))
	fmt.Printf("Copying backup to %s...\n", destPath)
	if err := copyFile(latestFile, destPath); err != nil {
		return fmt.Errorf("failed to copy backup file: %w", err)
	}

	// Step 4: Delete the server
	fmt.Printf("Deleting world '%s'...\n", name)
	cmd = exec.Command("sudo", "msm", "server", "delete", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	// Step 5: Clean up MSM archives
	fmt.Println("Cleaning up MSM archives...")

	// Remove backups directory
	cmd = exec.Command("sudo", "rm", "-rf", filepath.Join(msmArchivePath, "backups", name))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to remove backups directory: %v\n", err)
	}

	// Remove logs directory
	cmd = exec.Command("sudo", "rm", "-rf", filepath.Join(msmArchivePath, "logs", name))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to remove logs directory: %v\n", err)
	}

	// Remove worlds directory
	cmd = exec.Command("sudo", "rm", "-rf", filepath.Join(msmArchivePath, "worlds", name))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to remove worlds directory: %v\n", err)
	}

	fmt.Printf("World '%s' deleted successfully. Backup saved to %s\n", name, destPath)
	return nil
}
