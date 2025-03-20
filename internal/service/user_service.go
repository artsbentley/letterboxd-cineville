package service

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type UserProvider interface {
	GetAllUsers() ([]model.User, error)
	RegisterUser(string, string) error
}

type UserService struct {
	db.Querier
	Logger *slog.Logger
}

func NewUserService(conn *db.Queries) *UserService {
	return &UserService{
		Querier: conn,
		Logger:  slog.New(tint.NewHandler(os.Stderr, nil)),
	}
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	rows, err := s.GetUsers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	users := make([]model.User, len(rows))
	for _, user := range rows {
		users = append(users, model.User{
			ID:                 user.ID,
			Email:              user.Email,
			LetterboxdUsername: user.LetterboxdUsername,
			CreatedAt:          user.CreatedAt,
			Watchlist:          user.Watchlist,
		})
	}
	s.Logger.Info("Retrieved all users successfully")
	return users, nil
}

func (s *UserService) RegisterUser(email string, username string) error {
	user := db.CreateUserParams{
		Email:              email,
		LetterboxdUsername: username,
	}
	_, err := s.Querier.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to create user database entry: %v", err)
	}
	// TODO: notify the user that their account has been created
	return nil
}
