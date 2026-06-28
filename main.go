package main

import (
	"fmt"
	"os"

	"github.com/shadyendless/gator/internal/config"
)

func main() {
	fmt.Println("Loading config file...")
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("An error occurred: %w\n", err)
		os.Exit(1)
	}

	fmt.Println("Config file loaded.")
	fmt.Println(conf)

	fmt.Println("Updating config file...")
	err = conf.SetUser("jacob")
	if err != nil {
		fmt.Printf("An error occurred: %w\n", err)
		os.Exit(1)
	}
	fmt.Println("Config file updated.")

	fmt.Println("Loading config file...")
	conf, err = config.Read()
	if err != nil {
		fmt.Printf("An error occurred: %w\n", err)
		os.Exit(1)
	}

	fmt.Println("Config file loaded.")
	fmt.Println(conf)

	os.Exit(0)
}
