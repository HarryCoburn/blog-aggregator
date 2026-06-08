package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/HarryCoburn/blog-aggregator/internal/config"
	"github.com/HarryCoburn/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	var appState state
	var commands commands

	// Initialize config
	file, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	appState.cfg = &file

	// Initialize database
	db, err := sql.Open("postgres", appState.cfg.Db_url)
	dbQueries := database.New(db)
	appState.db = dbQueries

	// Initialize command map
	commands.commands = make(map[string]func(*state, command) error)

	// Register commands
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerGetUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerFeeds)

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
