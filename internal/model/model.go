package model

import (
	"time"
)

type FilmEvent struct {
	Name            string
	URL             string
	StartDate       time.Time
	EndDate         time.Time
	LocationName    string
	LocationAddress string
	OrganizerName   string
	OrganizerURL    string
	PerformerName   string
}

type User struct {
	ID                 int64      `json:"id"`
	Email              string     `json:"email"`
	LetterboxdUsername string     `json:"letterboxd_username"`
	CreatedAt          *time.Time `json:"created_at"`
	Watchlist          []string   `json:"watchlist"`
}
