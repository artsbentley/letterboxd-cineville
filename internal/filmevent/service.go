package filmevent

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/model"
)

type Provider interface {
	InsertFilmEvent(model.FilmEvent) error
	DeletePastFilmEvents() error
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// TODO: remember to implement this into main logic
func (s *Service) DeletePastFilmEvents() error {
	err := db.Store.DeletePastFilmEvents(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete film event rows that are in the past: ", err)
	}
	return nil
}

func (s *Service) InsertFilmEvent(event model.FilmEvent) error {
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
	_, err := db.Store.CreateFilmEvent(context.Background(), arg)
	return err
}
