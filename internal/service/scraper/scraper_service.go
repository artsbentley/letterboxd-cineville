package scraper

import (
	"letterboxd-cineville/internal/service"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type Scraper interface {
	Scrape() error
}

type ScraperService struct {
	scrapers []Scraper
	cron     *cron.Cron
}

func NewScraperService(
	userService service.UserProvider,
	watchlistService service.WatchlistProvider,
	filmEventService service.FilmEventProvider,
) *ScraperService {
	watchlistScraper := NewWatchlistScraper(userService, watchlistService)
	filmEventScraper := NewFilmEventScraper(userService, filmEventService)

	return &ScraperService{
		scrapers: []Scraper{watchlistScraper, filmEventScraper},
		cron:     cron.New(cron.WithLocation(time.FixedZone("CET", 1*60*60))),
	}
}

// TODO:
// does this func need a Go routine?
func (s *ScraperService) Start() {
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
