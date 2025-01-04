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
	Email              string
	LetterboxdUsername string
	Watchlist          []string
	Token              string
}

// type User struct {
// 	Username     string
// 	CityInterest []string
// }
