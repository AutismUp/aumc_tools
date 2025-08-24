package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AutismUp/aumc_tools/internal/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration files",
	Long: `Creates default config.json and server.properties.template files in the current directory.
These files should be customized and the AU_CONFIG_FILE environment variable should be set
to point to the config.json file location.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to create default configuration files? (y/n): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			if err := config.SaveDefaultFiles(); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating default files: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Default configuration files created successfully:")
			fmt.Println("- config.json")
			fmt.Println("- server.properties.template")
			fmt.Println()
			fmt.Println("Please:")
			fmt.Println("1. Update these files with your desired configurations")
			fmt.Println("2. Place them in your preferred location")
			fmt.Println("3. Set the AU_CONFIG_FILE environment variable to the path of config.json")
			fmt.Println("   Example: export AU_CONFIG_FILE=/path/to/your/config.json")
		} else {
			fmt.Println("Configuration file creation cancelled.")
			fmt.Println("Set the AU_CONFIG_FILE environment variable to point to your config file.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}