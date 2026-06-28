package state

import (
	"database/sql"
	"fmt"

	"github.com/shadyendless/gator/internal/config"
	"github.com/shadyendless/gator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

func New() (State, error) {
	conf, err := config.Read()
	if err != nil {
		return State{}, fmt.Errorf("An error occurred: %v\n", err)
	}

	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		return State{}, fmt.Errorf("An error occurred: %v\n", err)
	}

	dbQueries := database.New(db)

	return State{
		Db:     dbQueries,
		Config: &conf,
	}, nil
}
