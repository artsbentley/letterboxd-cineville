package model

import (
	"time"
)

type FilmEvent struct {
	Name            string    `db:"name"`
	URL             string    `db:"url"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
	LocationName    string    `db:"location_name"`
	LocationAddress string    `db:"location_address"`
	OrganizerName   string    `db:"organizer_name"`
	OrganizerURL    string    `db:"organizer_url"`
	PerformerName   string    `db:"performer_name"`
}
