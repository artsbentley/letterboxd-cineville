package user

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
	"log/slog"
)

type Provider interface {
	GetAllUsers() ([]model.User, error)
	RegisterUser(string, string, []string) error
}

type Service struct {
	db.Querier
}

func NewService(conn *db.Queries) *Service {
	return &Service{
		Querier: conn,
	}
}

func (s *Service) GetAllUsers() ([]model.User, error) {
	ctx := context.Background()
	rows, err := s.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	users := make([]model.User, len(rows))
	for _, user := range rows {
		locations, err := s.Querier.GetUserLocationCities(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get the locations belonging to the user:", err)
		}
		users = append(users, model.User{
			ID:                 user.ID,
			Email:              user.Email,
			LetterboxdUsername: user.LetterboxdUsername,
			CreatedAt:          user.CreatedAt,
			Watchlist:          user.Watchlist,
			Locations:          locations,
		})
	}
	slog.Info("Retrieved all users successfully")
	return users, nil
}

func (s *Service) RegisterUser(email string, username string, locations []string) error {
	ctx := context.Background()
	user := db.CreateUserParams{
		Email:              email,
		LetterboxdUsername: username,
	}
	userEntry, err := s.Querier.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user database entry: %v", err)
	}
	// TODO: notify the user that their account has been created
	for _, location := range locations {
		locationID, err := s.Querier.CreateLocation(ctx, location)
		if err != nil {
			return fmt.Errorf("failed to create location entry in database: %v", err)
		}
		err = s.Querier.AssignUserLocation(ctx, db.AssignUserLocationParams{
			UserID:     userEntry.ID,
			LocationID: locationID,
		})
		// FIX: things might go wrong here, i dont fully trust the querying logic
		// yet, validate this later
		if err != nil {
			return fmt.Errorf("failed to assign loction to user: %v", err)
		}
	}
	return nil
}
