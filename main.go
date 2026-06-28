package main

import (
	"fmt"
	"os"

	"github.com/shadyendless/gator/internal/commands"
	"github.com/shadyendless/gator/internal/state"
)

func main() {
	s, err := state.New()
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	coms := commands.New()
	coms.Register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("[ERROR]: Not enough arguments were passed")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err = coms.Run(&s, commands.Command{
		Name: command,
		Args: args,
	}); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func handlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]
	if err := s.Config.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", username)
	return nil
}
