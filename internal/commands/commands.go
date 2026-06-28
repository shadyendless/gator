package commands

import (
	"fmt"

	"github.com/shadyendless/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commands map[string]func(*state.State, Command) error
}

func New() Commands {
	return Commands{
		commands: make(map[string]func(*state.State, Command) error),
	}
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	command, exists := c.commands[cmd.Name]
	if !exists {
		return fmt.Errorf("The command `%s` is not registered", cmd.Name)
	}

	return command(s, cmd)
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.commands[name] = f
}
