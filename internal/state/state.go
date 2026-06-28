package state

import (
	"fmt"

	"github.com/shadyendless/gator/internal/config"
)

type State struct {
	Config *config.Config
}

func New() (State, error) {
	conf, err := config.Read()
	if err != nil {
		return State{}, fmt.Errorf("An error occurred: %v\n", err)
	}

	return State{
		Config: &conf,
	}, nil
}
