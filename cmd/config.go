package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkConfigCmd represents the check-config command
var checkConfigCmd = &cobra.Command{
	Use:   "check-config",
	Short: "Display current configuration",
	Long:  `Prints the current configuration to the screen for verification.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("msm_path") {
			fmt.Fprintln(os.Stderr, "No configuration loaded. Run 'aumc init' to create default config files.")
			os.Exit(1)
		}

		// Get all config as a map and pretty print it
		allSettings := viper.AllSettings()
		configJSON, err := json.MarshalIndent(allSettings, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Current Configuration:")
		fmt.Println(string(configJSON))
	},
}

// reloadConfigCmd represents the reload-config command
var reloadConfigCmd = &cobra.Command{
	Use:   "reload-config",
	Short: "Reload configuration file",
	Long:  `Reloads the configuration file from disk.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Re-read the config file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration reloaded successfully from:", viper.ConfigFileUsed())
	},
}

func init() {
	rootCmd.AddCommand(checkConfigCmd)
	rootCmd.AddCommand(reloadConfigCmd)
}