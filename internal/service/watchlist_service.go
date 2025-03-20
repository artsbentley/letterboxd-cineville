package service

import (
	"context"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type WatchlistProvider interface {
	InsertWatchlist(model.User) error
}

// TODO:
// should take in a sraper struct/ interface
type WatchlistService struct {
	db.Querier
	UserService UserProvider
	Logger      *slog.Logger
}

func NewWatchlistService(conn *db.Queries, userProvider UserProvider) *WatchlistService {
	return &WatchlistService{
		Querier:     conn,
		UserService: userProvider,
		Logger:      slog.New(tint.NewHandler(os.Stderr, nil)),
	}
}

func (s *WatchlistService) InsertWatchlist(user model.User) error {
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
