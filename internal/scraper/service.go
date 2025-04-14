package scraper

import (
	"letterboxd-cineville/internal/filmevent"
	"letterboxd-cineville/internal/user"
	"letterboxd-cineville/internal/watchlist"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type ScraperProvider interface {
	Scrape() error
}

type Service struct {
	scrapers []ScraperProvider
	cron     *cron.Cron
}

func NewService(
	userService user.Provider,
	watchlistService watchlist.Provider,
	filmEventService filmevent.Provider,
) *Service {
	watchlistScraper := NewWatchlistScraper(userService, watchlistService)
	filmEventScraper := NewFilmEventScraper(userService, filmEventService)

	return &Service{
		scrapers: []ScraperProvider{watchlistScraper, filmEventScraper},
		cron:     cron.New(cron.WithLocation(time.FixedZone("CET", 1*60*60))),
	}
}

// TODO:
// does this func need a Go routine?
func (s *Service) Start() {
	cronExpr := "* * * * *"
	for _, scraper := range s.scrapers {
		scraper := scraper // avoid closure capture issues
		_, err := s.cron.AddFunc(cronExpr, func() {
			slog.Debug("Scheduled task running...")
			if err := scraper.Scrape(); err != nil {
				slog.Error("Scraper failed to run", "error", err)
			}
		})
		if err != nil {
			slog.Error("Error adding cron function", "error", err)
		}
	}
	// Start the cron scheduler only once after adding all jobs
	s.cron.Start()
}
