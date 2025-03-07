package main

import (
	"fmt"
	"log"
)

type command struct {
	Name string
	Args []string
}

func handlerLogin(state *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("command %s expects [name] argument", cmd.Name)
	}
	username := cmd.Args[0]
	err := state.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	log.Printf("%s set as current user\n", username)
	return nil
}

type commandRegistry struct {
	reg map[string]func(*state, command) error
}

func (cmdReg *commandRegistry) register(name string, function func(*state, command) error) {
	cmdReg.reg[name] = function
}

func (cmdReg *commandRegistry) run(s *state, cmd command) error {
	command, ok := cmdReg.reg[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s not found\n", cmd.Name)
	}
	return command(s, cmd)
}
