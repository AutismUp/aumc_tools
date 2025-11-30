package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("AUMC - Autism Up Minecraft Tool")
	fmt.Printf("Cobra imported: %T\n", &cobra.Command{})
	fmt.Printf("Viper loaded: %v\n", viper.New() != nil)
	os.Exit(0)
}
