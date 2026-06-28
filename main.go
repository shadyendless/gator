package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/shadyendless/gator/internal/commands"
	"github.com/shadyendless/gator/internal/database"
	"github.com/shadyendless/gator/internal/state"
	"github.com/shadyendless/gator/internal/xml"
)

func main() {
	s, err := state.New()
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	comms := commands.New()
	comms.Register("login", handlerLogin)
	comms.Register("register", handlerRegister)
	comms.Register("reset", handlerReset)
	comms.Register("users", handlerUsers)
	comms.Register("agg", handlerAgg)

	if len(os.Args) < 2 {
		fmt.Println("[ERROR]: Not enough arguments were passed")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err = comms.Run(&s, commands.Command{
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
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	if err = s.Config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", user.Name)
	return nil
}

func handlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("name is required")
	}

	name := cmd.Args[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Registered the following user: %v\n", user)
	if err = handlerLogin(s, commands.Command{
		Name: "login",
		Args: []string{user.Name},
	}); err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state.State, cmd commands.Command) error {
	if err := s.Db.Reset(context.Background()); err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state.State, cmd commands.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)

		if user.Name == s.Config.CurrentUserName {
			fmt.Print(" (current)")
		}

		fmt.Print("\n")
	}

	return nil
}

func handlerAgg(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("url is required")
	}

	url := cmd.Args[0]
	feed, err := xml.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not fetch from \"%s\"", url)
	}

	fmt.Println(feed)

	return nil
}
