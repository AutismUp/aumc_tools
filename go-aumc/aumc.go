package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	// create-world flags
	createWorldCmd := flag.NewFlagSet("create-world", flag.ExitOnError)
	createWorldName := createWorldCmd.String("name", "", "Name of the world")
	createWorldVersion := createWorldCmd.String("version", "", "Minecraft version to use")

	// delete-world flags
	deleteWorldCmd := flag.NewFlagSet("delete-world", flag.ExitOnError)
	deleteWorldName := deleteWorldCmd.String("Name", "", "Name of the world")

	if len(os.Args) < 2 || os.Args[1] == "-h" {
		helpCmd()
		os.Exit(0)
	}

	switch os.Args[1] {

	case "create-world":
		createWorldCmd.Parse(os.Args[2:])
		jarGroup := strings.ReplaceAll(*createWorldVersion, ".", "_")
		fmt.Println("Creating a new Minecraft world with the following parameters:")
		fmt.Println("  Name:", *createWorldName)
		fmt.Println("  Minecraft version:", *createWorldVersion)
		fmt.Println("  Jargroup name:", jarGroup)

	case "delete-world":
		deleteWorldCmd.Parse(os.Args[2:])
		fmt.Println("Deleting the following world")
		fmt.Println("  Name:", *deleteWorldName)
	default:
		fmt.Printf("%v is not a valid subcommand\n\n", os.Args[1])
		helpCmd()
		os.Exit(1)
	}
}

func helpCmd() {
	fmt.Print("\nAvailable subcommands:\n\n")
	fmt.Println("  create-world - Create a Minecraft world")
	fmt.Println("  delete-world - Delete a Minecraft world")
	fmt.Print("  get-config   - Check the aumc configuration\n\n")
	fmt.Print("Type '-h' after each subcommand for more help\n\n")

}
