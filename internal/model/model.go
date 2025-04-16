package model

import (
	"time"

	"github.com/google/uuid"
)

type FilmEvent struct {
	Name            string
	URL             string
	StartDate       time.Time
	EndDate         time.Time
	LocationName    string
	LocationAddress string
	City            string
	OrganizerName   string
	OrganizerURL    string
	PerformerName   string
}

type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	LetterboxdUsername string     `json:"letterboxd_username"`
	CreatedAt          *time.Time `json:"created_at"`
	Watchlist          []string   `json:"watchlist"`
	Locations          []string   `json:"locations"`
}

type UserFilmMatch struct {
	UserEmail   string
	FilmMatches []FilmEvent
}
