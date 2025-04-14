package main

import (
	"context"
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/filmevent"
	"letterboxd-cineville/internal/scraper"
	"letterboxd-cineville/internal/user"
	"letterboxd-cineville/internal/watchlist"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
)

var (
	dbName = "postgres"
	dbHost = "localhost:5432"
	dbUser = "app"
	url    = fmt.Sprintf("postgresql://%s@%s/%s", dbName, dbHost, dbUser)
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))
	// TODO: use gum?
	// slog.SetDefault(slog.New(log.NewWithOptions(os.Stderr, log.Options{
	// 	ReportTimestamp: true,
	// })))

	pgx, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to create a connection pool: %v\n", err)
	}

	conn := db.New(pgx)

	// userService := service.NewUserService(conn)
	userService := user.NewService(conn)
	// watchlistService := service.NewWatchlistService(conn, userService)
	watchlistService := watchlist.NewService(conn, userService)
	// filmEventService := service.NewFilmEventService(conn)
	filmEventService := filmevent.NewService(conn)

	_ = userService.RegisterUser("arnoarts@hotmail.com", "deltore", []string{"amsterdam", "leiden"})
	_ = userService.RegisterUser("sannilehtonen@gmail.com", "sannisideup", []string{"leiden", "utrecht"})
	_ = userService.RegisterUser("bananaman@gmail.com", "hahahah", []string{"leiden", "utrecht", "amsterdam", "lisse"})

	scrapeService := scraper.NewService(userService, watchlistService, filmEventService)

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
