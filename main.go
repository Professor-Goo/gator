package main

import (
	"fmt"
	"os"

	"github.com/Professor-Goo/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.SetUser("jmcgl"); err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		os.Exit(1)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config after update: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config contents:\n")
	fmt.Printf("  DB URL: %s\n", cfg.DbURL)
	fmt.Printf("  Current User: %s\n", cfg.CurrentUserName)
}
