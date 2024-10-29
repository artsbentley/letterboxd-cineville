package main

import (
	"fmt"
	database "letterboxd-cineville/db"
	"letterboxd-cineville/handle"
	"letterboxd-cineville/model"
	"letterboxd-cineville/scrape"

	_ "github.com/mattn/go-sqlite3"
)

// NOTE: nice to have:
// - scrape and persist movie posters per movie title

// TODO:
//   - setup scraping of filmevents into main file, first every time program is
//     run, later as concurrent cron
func main() {
	Sqlite := database.Sql

	filmEvents, err := scrape.CollectFilmEvents("https://www.filmvandaag.nl/filmladder/stad/13-amsterdam")
	handle.ErrFatal(err)

	for _, event := range filmEvents {
		err := Sqlite.InsertFilmEvent(event)
		handle.ErrFatal(err)
	}

	watchlist, err := scrape.ScrapeWatchlist("deltore")
	handle.ErrFatal(err)

	lbox := model.Letterboxd{
		Email:     "arnoarts@hotmail.com",
		Username:  "deltore",
		Watchlist: watchlist,
	}

	err = Sqlite.InsertWatchlist(lbox)
	handle.ErrFatal(err)

	matches, err := Sqlite.GetMatchingFilmEventsByEmail("arnoarts@hotmail.com")
	handle.ErrFatal(err)

	for _, match := range matches {
		fmt.Println(match.Name)
	}
}
