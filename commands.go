package main

import "fmt"

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
	s.cfg.SetUser(cmd.args[0])

	fmt.Printf("Username set to %s\n", cmd.args[0])
	return nil

}
