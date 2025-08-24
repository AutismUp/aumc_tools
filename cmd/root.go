package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aumc",
	Short: "Autism Up Minecraft Tool - CLI wrapper for managing Minecraft servers",
	Long: `AUMC is a CLI tool for managing Minecraft servers specifically for AU's configuration.
The tool handles server creation, jar building, world management, and configuration 
for Minecraft servers using MSM (Minecraft Server Manager).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global persistent flag for config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $AU_CONFIG_FILE or ./config.json)")
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Check for AU_CONFIG_FILE environment variable
		if configPath := os.Getenv("AU_CONFIG_FILE"); configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			// Find config in current directory with name "config" (without extension).
			viper.AddConfigPath(".")
			viper.SetConfigType("json")
			viper.SetConfigName("config")
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Offer to create default config if none found
		fmt.Fprintln(os.Stderr, "WARNING: No config file found")
		fmt.Fprintln(os.Stderr, "Set AU_CONFIG_FILE environment variable or use --config flag")
		fmt.Fprintln(os.Stderr, "Run 'aumc init' to create default configuration files")
	}
}