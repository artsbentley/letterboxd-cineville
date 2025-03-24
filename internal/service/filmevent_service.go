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

type FilmEventProvider interface {
	InsertFilmEvent(model.FilmEvent) error
	DeletePastFilmEvents() error
}

type FilmEventService struct {
	db.Querier
	Logger *slog.Logger
}

func NewFilmEventService(conn *db.Queries) *FilmEventService {
	return &FilmEventService{
		Querier: conn,
		Logger:  slog.New(tint.NewHandler(os.Stderr, nil)),
	}
}

// TODO: remember to implement this into main logic
func (s *FilmEventService) DeletePastFilmEvents() error {
	err := s.Querier.DeletePastFilmEvents(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete film event rows that are in the past: ", err)
	}
	return nil
}

func (s *FilmEventService) InsertFilmEvent(event model.FilmEvent) error {
	arg := db.CreateFilmEventParams{
		Name:            event.Name,
		Url:             event.URL,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		LocationName:    event.LocationName,
		LocationAddress: event.LocationAddress,
		City:            event.City,
		OrganizerName:   event.OrganizerName,
		OrganizerUrl:    event.OrganizerURL,
		PerformerName:   event.PerformerName,
	}
	_, err := s.Querier.CreateFilmEvent(context.Background(), arg)
	return err
}
