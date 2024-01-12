package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	// create-world flags
	createWorldCmd := flag.NewFlagSet("create-new-world", flag.ExitOnError)
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

	case "create-new-world":
		createWorldCmd.Parse(os.Args[2:])
		createNewWorld(*createWorldName, *createWorldVersion, "/opt/msm/servers")

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
	fmt.Println("  create-new-world - Create a Minecraft world")
	fmt.Println("  delete-world - Delete a Minecraft world")
	fmt.Print("  get-config   - Check the aumc configuration\n\n")
	fmt.Print("Type '-h' after each subcommand for more help\n\n")

}

func createNewWorld(name string, version string, msmServerPath string) {
	fmt.Printf("Creating individual world named %s\n", name)
	timeStamp := time.Now().Format("Mon Jan 02 15:04:05 MST 2006")
	jarGroup := strings.ReplaceAll(version, ".", "_")

	// 1 - create the world
	createCmd := exec.Command("sudo", "msm", "server", "create", name)
	createErr := createCmd.Run()
	if createErr != nil {
		log.Fatal(createErr)
	}

	jarCmd := exec.Command("sudo", "msm", name, "jar", jarGroup)
	jarErr := jarCmd.Run()
	if jarErr != nil {
		log.Fatal(jarErr)
	}

	// 2 - create the eula.txt file
	eulaFilePath := filepath.Join(msmServerPath, name, "eula.txt")
	eulaFile, err := os.Create(eulaFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer eulaFile.Close()

	_, err = eulaFile.WriteString("#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).\n")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eulaFile.WriteString(fmt.Sprintf("#%s\n", timeStamp))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = eulaFile.WriteString("eula=true")
	if err != nil {
		log.Fatalln(err)
	}

	// # 3 - update server.properties template and copy to the server folder
	// server_properties = MCConfig(self.config['world_config']['server_properties_template'])
	// server_properties.update_config('msm-version', f'minecraft/{version}')
	// server_properties.update_config('motd', f'Autism Up - {name}')
	// server_properties.write_config(f'{msm_server_path}/{name}/server.properties')

	// subprocess.call(['sudo', 'msm', name, 'start'])

	// for operator in self.config['op_usernames']:
	// 	subprocess.call(['sudo', 'msm', name, 'op', 'add', operator])

	// subprocess.call(['sudo', 'msm', name, 'stop', 'now'])
	// subprocess.call(['sudo', 'msm', name, 'worlds', 'ram', 'world'])

	// subprocess.run(['sudo', 'chown', '-R', 'minecraft', f'{msm_server_path}/{name}'])
	// subprocess.run(['sudo', 'chgrp', '-R', 'minecraft', f'{msm_server_path}/{name}'])

	// print(f'World named "{name}" created')
}
