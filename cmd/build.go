package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// buildNewJarCmd represents the build-new-jar command
var buildNewJarCmd = &cobra.Command{
	Use:   "build-new-jar",
	Short: "Build the latest version of Spigot Minecraft",
	Long: `Builds the latest version of Spigot Minecraft using BuildTools
and copies it to the Git repo for publication.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building new Spigot jar...")
		// TODO: Implement build functionality
		fmt.Println("Build functionality not yet implemented")
	},
}

// publishNewJarCmd represents the publish-new-jar command
var publishNewJarCmd = &cobra.Command{
	Use:   "publish-new-jar",
	Short: "Publish jarfile to GitHub and create MSM jargroup",
	Long: `Push the specified jarfile of Minecraft to GitHub and create a new JarGroup in MSM.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		
		fmt.Printf("Publishing jar file: %s\n", filename)
		fmt.Println("This publishes the new jar to GitHub")
		// TODO: Implement publish functionality
		fmt.Println("Publish functionality not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(buildNewJarCmd)
	rootCmd.AddCommand(publishNewJarCmd)

	// Add flags for publish command
	publishNewJarCmd.Flags().StringP("filename", "f", "", "Name of the jarfile to publish")
}