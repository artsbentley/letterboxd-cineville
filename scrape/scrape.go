package scrape

import (
	"letterboxd-cineville/db"
	"log/slog"

	"github.com/robfig/cron/v3"
)

type Scraper interface {
	Scrape()
}

type WatchlistScraper struct {
	logger *slog.Logger
	db     *db.Sqlite
	cron   *cron.Cron
}

type FilmEventScraper struct {
	logger *slog.Logger
	db     *db.Sqlite
	cron   *cron.Cron
}
