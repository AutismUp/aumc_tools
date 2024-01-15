package main

import (
	"bufio"
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
	createWorldOps := createWorldCmd.String("ops", "", "List of operators, separted by commas (ie. op1,op2,op3)")

	// delete-world flags
	deleteWorldCmd := flag.NewFlagSet("delete-world", flag.ExitOnError)
	deleteWorldName := deleteWorldCmd.String("name", "", "Name of the world")

	if len(os.Args) < 2 || os.Args[1] == "-h" {
		helpCmd()
		os.Exit(0)
	}

	switch os.Args[1] {

	case "create-new-world":
		createWorldCmd.Parse(os.Args[2:])
		// fmt.Println(*createWorldName, *createWorldVersion, *createWorldOps)
		createNewWorld(*createWorldName, *createWorldVersion, "/opt/msm/servers", *createWorldOps)

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
	fmt.Print("  server-property   - Manage the Minecraft world properties\n\n")
	fmt.Print("Type '-h' after each subcommand for more help\n\n")

}

// createNewWorld creates a new Minecraft world with the provided name, version, and operators
func createNewWorld(name string, version string, msmServerPath string, ops string) {
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
	fmt.Println("Agreeing the Minecraft EULA")
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
	fmt.Println("Updating the server.properties file")
	sp := CreateServerProperties()
	sp.Set("msm-version", fmt.Sprintf("minecraft/%s", version))
	sp.Set("motd", fmt.Sprintf("Autism Up - %v", name))
	sp.Save(fmt.Sprintf("%s/%s/server.properties", msmServerPath, name))

	startCmd := exec.Command("sudo", "msm", name, "start")
	startErr := startCmd.Run()
	if createErr != nil {
		log.Fatal(startErr)
	}

	fmt.Println("Adding the operators to the world")
	opList := strings.Split(ops, ",")
	for _, op := range opList {
		opAddCmd := exec.Command("sudo", "msm", name, "op", "add", op)
		opAddErr := opAddCmd.Run()
		if opAddErr != nil {
			log.Fatal(opAddErr)
		}
	}

	stopCmd := exec.Command("sudo", "msm", name, "stop")
	stopErr := stopCmd.Run()
	if stopErr != nil {
		log.Fatal(stopErr)
	}

	fmt.Println("Setting world to load into RAM")
	ramCmd := exec.Command("sudo", "msm", name, "worlds", "ram", "world")
	ramErr := ramCmd.Run()
	if ramErr != nil {
		log.Fatal(ramErr)
	}

	fmt.Println("Updating folder owner")
	chownCmd := exec.Command("sudo", "chown", "-R", "minecraft", fmt.Sprintf("%s/%s", msmServerPath, name))
	chownErr := chownCmd.Run()
	if chownErr != nil {
		log.Fatal(chownErr)
	}

	fmt.Println("Updating folder group")
	chgrpCmd := exec.Command("sudo", "chgrp", "-R", "minecraft", fmt.Sprintf("%s/%s", msmServerPath, name))
	chgrpErr := chgrpCmd.Run()
	if chgrpErr != nil {
		log.Fatal(chgrpErr)
	}

	fmt.Printf("World named '%s' created on Minecraft version %s\n", name, version)
}

// ServerProperties represents the properties of the Minecraft server.
type ServerProperties map[string]string

func CreateServerProperties() ServerProperties {
	defaultSp := ServerProperties{
		"msm-version":                       "minecraft/1.18.2",
		"spawn-protection":                  "16",
		"max-tick-time":                     "60000",
		"query.port":                        "25565",
		"generator-settings":                "",
		"sync-chunk-writes":                 "true",
		"force-gamemode":                    "false",
		"allow-nether":                      "true",
		"enforce-whitelist":                 "true",
		"gamemode":                          "creative",
		"broadcast-console-to-ops":          "true",
		"enable-query":                      "false",
		"player-idle-timeout":               "0",
		"text-filtering-config":             "",
		"difficulty":                        "easy",
		"spawn-monsters":                    "true",
		"broadcast-rcon-to-ops":             "true",
		"op-permission-level":               "4",
		"pvp":                               "true",
		"entity-broadcast-range-percentage": "100",
		"snooper-enabled":                   "true",
		"level-type":                        "default",
		"hardcore":                          "false",
		"enable-status":                     "true",
		"enable-command-block":              "false",
		"max-players":                       "20",
		"network-compression-threshold":     "256",
		"resource-pack-sha1":                "",
		"max-world-size":                    "29999984",
		"function-permission-level":         "2",
		"rcon.port":                         "25575",
		"server-port":                       "25565",
		"debug":                             "false",
		"server-ip":                         "",
		"spawn-npcs":                        "true",
		"allow-flight":                      "false",
		"level-name":                        "world",
		"view-distance":                     "10",
		"resource-pack":                     "",
		"spawn-animals":                     "true",
		"white-list":                        "true",
		"rcon.password":                     "",
		"generate-structures":               "true",
		"max-build-height":                  "256",
		"online-mode":                       "true",
		"level-seed":                        "",
		"use-native-transport":              "true",
		"prevent-proxy-connections":         "false",
		"enable-jmx-monitoring":             "false",
		"enable-rcon":                       "false",
		"rate-limit":                        "0",
		"motd":                              "AU",
	}
	return defaultSp
}

// LoadServerProperties loads the server properties from a file.
func LoadServerProperties(filePath string) (ServerProperties, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	properties := ServerProperties{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			properties[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return properties, scanner.Err()
}

// Save writes the server properties back to the file.
func (sp ServerProperties) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range sp {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return err
		}
	}

	return nil
}

// Get a property value by key.
func (sp ServerProperties) Get(key string) (string, bool) {
	value, exists := sp[key]
	return value, exists
}

// Set a property value by key.
func (sp ServerProperties) Set(key, value string) {
	sp[key] = value
}
