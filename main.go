package main

import (
	"fmt"
	"os"

	"github.com/HarryCoburn/blog-aggregator/internal/config"
)

type state struct {
	state *config.Config
}

func main() {
	var appState state
	var commands commands

	// Initialize state
	file, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	appState.state = &file

	// Initialize command map
	commands.commands = make(map[string]func(*state, command) error)

	// Register commands
	commands.register("login", handlerLogin)

	// Check arguments
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments provided. Exiting.")
		os.Exit(1)
	}

	// Build command
	var command command
	command.name = args[1]
	if len(args) > 2 {
		command.args = args[2:]
	}

	// Run command
	err = commands.run(&appState, command)
	if err != nil {
		fmt.Println("Error occurred:", err)
		os.Exit(1)
	}

}
