package main

import (
	database "letterboxd-cineville/db"
	"letterboxd-cineville/handlers"
	"letterboxd-cineville/scrape"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// NOTE: nice to have:
// - scrape and persist movie posters per movie title

// NOTE: considerations:
// - lower case before insert of values or right before matching?
// - image storage
// - smtp server

// TODO:
//   - setup scraping of filmevents into main file, first every time program is
//     run, later as concurrent cron
func main() {
	Sqlite := database.Sql
	FilmEventScraper := scrape.NewFilmEventScraper(Sqlite)
	WatchlistScraper := scrape.NewWatchlistScraper(Sqlite)

	go FilmEventScraper.Scrape()
	go WatchlistScraper.Scrape()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(Sqlite)

	// Routes
	e.GET("/", userHandler.HandleGetUsers)
	e.POST("/users", userHandler.HandleCreateUser)

	// Start server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

// users, err := Sqlite.GetAllUsers()
// if err != nil {
// 	log.Fatal(err)
// }
//
// fmt.Println("Users in the database:")
// for _, user := range users {
// 	fmt.Printf("Email: %s, Username: %s\n", user.Email, user.Username)
// }
//
// filmEvents, err := scrape.CollectFilmEvents("https://www.filmvandaag.nl/filmladder/stad/13-amsterdam")
// if err != nil {
// 	log.Fatal(err)
// }
//
// for _, event := range filmEvents {
// 	err := Sqlite.InsertFilmEvent(event)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// watchlist, err := scrape.ScrapeWatchlist("deltore")
// if err != nil {
// 	log.Fatal(err)
// }
//
// lbox := model.Letterboxd{
// 	Email:     "arnoarts@hotmail.com",
// 	Username:  "deltore",
// 	Watchlist: watchlist,
// }
//
// err = Sqlite.InsertWatchlist(lbox)
// if err != nil {
// 	log.Fatal(err)
// }
//
// matches, err := Sqlite.GetMatchingFilmEventsByEmail("arnoarts@hotmail.com")
// if err != nil {
// 	log.Fatal(err)
// }
//
// for _, match := range matches {
// 	fmt.Println(match.Name, match.LocationName, match.StartDate)
// }
