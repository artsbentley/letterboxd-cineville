package main

import (
	"context"
	"fmt"
	database "letterboxd-cineville/db"
	"letterboxd-cineville/handlers"
	"letterboxd-cineville/scrape"
	"letterboxd-cineville/service"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

var (
	dbUser = "postgres"
	dbHost = "localhost:5432"
	url    = fmt.Sprintf("postgresql://%s@%s/%s", dbUser, dbHost, dbUser)
)

// TODO:
//   - setup scraping of filmevents into main file, first every time program is
//     run, later as concurrent cron
func main() {
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to create a connection pool: %v\n", err)
	}

	Store := database.Sql
	Service := service.NewService(Store)
	// FilmEventScraper := scrape.NewFilmEventScraper(Store)
	WatchlistScraper := scrape.NewWatchlistScraper(Store)

	// go FilmEventScraper.Scrape()
	go WatchlistScraper.Scrape()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(Service)
	// err := Store Store.DB.

	// Routes
	e.GET("/", userHandler.HandleGetUsers)
	e.POST("/users", userHandler.HandleCreateUser)

	// Start server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
