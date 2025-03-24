package main

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/service"
	"letterboxd-cineville/internal/service/scraper"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbName = "postgres"
	dbHost = "localhost:5432"
	dbUser = "app"
	url    = fmt.Sprintf("postgresql://%s@%s/%s", dbName, dbHost, dbUser)
)

func main() {
	pgx, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to create a connection pool: %v\n", err)
	}

	conn := db.New(pgx)

	userService := service.NewUserService(conn)
	watchlistService := service.NewWatchlistService(conn, userService)
	filmEventService := service.NewFilmEventService(conn)

	_ = userService.RegisterUser("arnoarts@hotmail.com", "deltore", []string{"amsterdam", "leiden"})
	_ = userService.RegisterUser("sannilehtonen@gmail.com", "sannisideup", []string{"leiden", "utrecht"})
	_ = userService.RegisterUser("bananaman@gmail.com", "hahahah", []string{"leiden", "utrecht", "amsterdam", "lisse"})

	scrapeService := scraper.NewScraperService(userService, watchlistService, filmEventService)

	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return
	}
	fmt.Fprintf(os.Stderr, "🚀 : main.go:38: users=%+v\n", users)
	scrapeService.Start()

	// Block forever so the cron scheduler can continue running
	select {}
}
