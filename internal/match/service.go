package service

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
	"letterboxd-cineville/internal/user"
	"log/slog"
)

// 1. get list of all users
// 1. get film events that match watchlist
// 2. everything lower case
// 3. use regex?

type Service struct {
	userService user.Provider
}

func NewService(userProvider user.Provider) *MatchService {
	return &MatchService{
		userService: userProvider,
	}
}

func (s *Service) match() ([]model.UserFilmMatch, error) {
	users, err := s.userService.GetAllUsers()
	if err != nil {
		slog.Error("failed to retrieve users from the database", "error", err)
	}
	matches, err := s.GetAllUsersFilmMatches(users)
	return matches, nil
}

func (s *Service) GetAllUsersFilmMatches(users []model.User) ([]model.UserFilmMatch, error) {
	var matches []model.UserFilmMatch
	for _, user := range users {
		match, err := s.GetUserFilmMatches(user)
		if err != nil {
			slog.Error("failed to retrieve film event matches for user", "user", user.Email, "error", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}

func (s *Service) GetUserFilmMatches(user model.User) (model.UserFilmMatch, error) {
	dbEvents, err := db.Store.MatchFilmEventsWithUser(context.Background(), user.Email)
	if err != nil {
		return model.UserFilmMatch{}, fmt.Errorf("failed to retrieve film event matches for a user: %v", err)
	}

	var filmEvents []model.FilmEvent
	for _, event := range dbEvents {
		filmEvent := model.FilmEvent{
			Name:            event.Name,
			URL:             event.Url,
			StartDate:       event.StartDate,
			EndDate:         event.EndDate,
			LocationName:    event.LocationName,
			LocationAddress: event.LocationAddress,
			City:            event.City,
			OrganizerName:   event.OrganizerName,
			OrganizerURL:    event.OrganizerUrl,
			PerformerName:   event.PerformerName,
		}
		filmEvents = append(filmEvents, filmEvent)
	}
	match := model.UserFilmMatch{
		UserEmail:   user.Email,
		FilmMatches: filmEvents,
	}
	return match, nil
}

// func (s *MatchService) Start()
