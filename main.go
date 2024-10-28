package main

import (
	"fmt"
	database "letterboxd-cineville/db"
	"letterboxd-cineville/model"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TODO:
// - scrape and persist movie posters per movie title

func main() {
	Sqlite := database.Sql
	dateStr := "2024-10-27 15:30:00"
	layout := "2006-01-02 15:04:05"

	// Parse the string to time.Time
	parsedTime, _ := time.Parse(layout, dateStr)
	event := model.FilmEvent{
		StartDate:       parsedTime,
		EndDate:         parsedTime,
		Name:            "The Substance",
		URL:             "https://themovies.nl",
		LocationName:    "The Movies",
		LocationAddress: "Haarlemmerstraat",
		OrganizerName:   "The Movies",
		OrganizerURL:    "www.themovies.nl",
		PerformerName:   "Performer",
	}

	err := Sqlite.InsertFilmEvent(event)
	if err != nil {
		log.Fatal(err)
	}

	lbox := model.Letterboxd{
		Email:     "arnoarts@hotmail.com",
		Username:  "Deltore",
		Watchlist: []string{"The Substance", "Persona"},
	}

	err = Sqlite.InsertWatchlist(lbox)
	if err != nil {
		log.Fatal(err)
	}

	match, err := Sqlite.GetMatchingFilmEventsByEmail("arnoarts@hotmail.com")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("We've found matches: ", match)
}
