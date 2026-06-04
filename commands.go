package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HarryCoburn/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	elem, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s is not registered", cmd.name)
	}
	err := elem(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("Login command requires a username.")
	}

	// Check if the user exists

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not login with user %s. Error: %v", cmd.args[0], err)
	}

	s.cfg.SetUser(user.Name)

	fmt.Printf("Username set to %s\n", user.Name)
	return nil

}

func handlerRegister(s *state, cmd command) error {
	if cmd.args == nil {
		return fmt.Errorf("Register command requires a username.")
	}
	fmt.Printf("Attempting to register user: %s\n", cmd.args[0])
	var params database.CreateUserParams
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = params.CreatedAt
	params.Name = cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Could not register user: %v", err)
	}
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("The returned user is: %v\n Setting current user to this.", user)

	return nil
}
