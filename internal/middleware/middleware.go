package middleware

import (
	"context"

	"github.com/shadyendless/gator/internal/commands"
	"github.com/shadyendless/gator/internal/database"
	"github.com/shadyendless/gator/internal/state"
)

func LoggedIn(handler func(s *state.State, cmd commands.Command, user database.User) error, s *state.State) func(*state.State, commands.Command) error {
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)

	return func(s *state.State, c commands.Command) error {
		if err != nil {
			return err
		}

		return handler(s, c, user)
	}
}
