package scraper

import (
	"letterboxd-cineville/internal/service"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/robfig/cron/v3"
)

type Scraper interface {
	Scrape() error
}

type ScraperService struct {
	scrapers []Scraper
	logger   *slog.Logger
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
		logger:   slog.New(tint.NewHandler(os.Stderr, nil)),
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
			s.logger.Debug("Scheduled task running...")
			if err := scraper.Scrape(); err != nil {
				s.logger.Error("Scraper failed to run", "error", err)
			}
		})
		if err != nil {
			s.logger.Error("Error adding cron function", "error", err)
		}
	}
	// Start the cron scheduler only once after adding all jobs
	s.cron.Start()
}
