package watchlist

import (
	"context"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
	"letterboxd-cineville/internal/user"
)

type Provider interface {
	InsertWatchlist(model.User) error
}

// TODO:
// should take in a sraper struct/ interface
type Service struct {
	db.Querier
	UserService user.Provider
}

func NewService(conn *db.Queries, userProvider user.Provider) *Service {
	return &Service{
		Querier:     conn,
		UserService: userProvider,
	}
}

func (s *Service) InsertWatchlist(user model.User) error {
	args := db.UpdateUserWatchlistParams{
		Email:     user.Email,
		Watchlist: user.Watchlist,
	}

	err := s.UpdateUserWatchlist(context.Background(), args)
	if err != nil {
		return err
	}
	return nil
}
