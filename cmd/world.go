package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createNewWorldCmd represents the create-new-world command
var createNewWorldCmd = &cobra.Command{
	Use:   "create-new-world",
	Short: "Create a new Minecraft world",
	Long: `Create a new world using Autism Up default configurations.
Can create a single world by name or all worlds listed in the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fromConfig, _ := cmd.Flags().GetBool("from-config")
		jargroup, _ := cmd.Flags().GetString("jargroup")
		version, _ := cmd.Flags().GetString("version")

		if fromConfig {
			fmt.Println("Creating worlds defined in the configuration file.")
			
			// Get config values
			configJargroup := viper.GetString("world_config.jargroup")
			configVersion := viper.GetString("world_config.minecraft_version")
			worldNames := viper.GetStringSlice("world_config.world_names")

			if len(worldNames) == 0 {
				fmt.Fprintln(os.Stderr, "No world names found in configuration")
				os.Exit(1)
			}

			for _, world := range worldNames {
				fmt.Printf("Creating world: %s\n", world)
				// TODO: Implement world creation
				fmt.Printf("  Jargroup: %s, Version: %s\n", configJargroup, configVersion)
			}
		} else {
			if name == "" {
				fmt.Fprintln(os.Stderr, "World name is required when not using --from-config")
				os.Exit(1)
			}
			
			fmt.Printf("Creating individual world named: %s\n", name)
			// TODO: Implement world creation
			fmt.Printf("  Jargroup: %s, Version: %s\n", jargroup, version)
		}
		
		fmt.Println("World creation functionality not yet implemented")
	},
}

// deleteWorldCmd represents the delete-world command
var deleteWorldCmd = &cobra.Command{
	Use:   "delete-world",
	Short: "Delete a Minecraft world",
	Long: `Deletes a world using Autism Up configurations.
Creates a backup before deletion. Can delete a single world by name or all worlds listed in the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		fromConfig, _ := cmd.Flags().GetBool("from-config")

		if fromConfig {
			fmt.Println("Deleting all the worlds from the config file.")
			
			worldNames := viper.GetStringSlice("world_config.world_names")
			if len(worldNames) == 0 {
				fmt.Fprintln(os.Stderr, "No world names found in configuration")
				os.Exit(1)
			}

			for _, world := range worldNames {
				fmt.Printf("Deleting world: %s\n", world)
				// TODO: Implement world deletion
			}
		} else {
			if name == "" {
				fmt.Fprintln(os.Stderr, "World name is required when not using --from-config")
				os.Exit(1)
			}
			
			fmt.Printf("Deleting world: %s\n", name)
			// TODO: Implement world deletion
		}
		
		fmt.Println("World deletion functionality not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(createNewWorldCmd)
	rootCmd.AddCommand(deleteWorldCmd)

	// Add flags for create-new-world command
	createNewWorldCmd.Flags().StringP("name", "n", "", "Name of the Minecraft server to create")
	createNewWorldCmd.Flags().Bool("from-config", false, "Create all worlds listed in the configuration file")
	createNewWorldCmd.Flags().StringP("jargroup", "j", "", "Jargroup to use for the server")
	createNewWorldCmd.Flags().StringP("version", "v", "", "Version of Minecraft that will be used")

	// Add flags for delete-world command
	deleteWorldCmd.Flags().StringP("name", "n", "", "Name of the Minecraft server to delete")
	deleteWorldCmd.Flags().Bool("from-config", false, "Delete all worlds listed in the configuration file")
}