package main

import (
	"fmt"
	"letterboxd-cineville/internal/db"
	"letterboxd-cineville/internal/filmevent"
	"letterboxd-cineville/internal/scraper"
	"letterboxd-cineville/internal/user"
	"letterboxd-cineville/internal/watchlist"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/lmittmann/tint"
)

func InitEverything() error {
	if err := db.Init(); err != nil {
		return err
	}
	return nil
}

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))

	if err := InitEverything(); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// TODO: use gum?
	// slog.SetDefault(slog.New(log.NewWithOptions(os.Stderr, log.Options{
	// 	ReportTimestamp: true,
	// })))

	userService := user.NewService()
	filmEventService := filmevent.NewService()
	watchlistService := watchlist.NewService(userService)
	scrapeService := scraper.NewService(userService, watchlistService, filmEventService)

	_ = userService.RegisterUser("arnoarts@hotmail.com", "deltore", []string{"amsterdam", "leiden"})
	_ = userService.RegisterUser("sannilehtonen@gmail.com", "sannisideup", []string{"leiden", "utrecht"})
	_ = userService.RegisterUser("bananaman@gmail.com", "hahahah", []string{"leiden", "utrecht", "amsterdam", "lisse"})

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
