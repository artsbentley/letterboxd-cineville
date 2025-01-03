package scrape

import (
	"fmt"
	"letterboxd-cineville/db"
	"log"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

func (s WatchlistScraper) Scrape() {
	// Run every day at 2 AM
	// cronExpr := "0 2 * * *"
	cronExpr := "* * * * *"
	_, err := s.cron.AddFunc(cronExpr, func() {
		s.logger.Info("WatchlistScraper Scheduled task running...")

		// Get all users
		allUsers, err := s.db.GetAllUsers()
		if err != nil {
			s.logger.Error("failed to get users", "error", err)
			return
		}

		for _, user := range allUsers {
			watchlist, err := ScrapeWatchlist(user.LetterboxdUsername)
			if err != nil {
				s.logger.Warn("failed to scrape watchlist",
					"user", user.LetterboxdUsername,
					"email", user.Email,
					"error", err)
				continue
			}

			// Update user with new watchlist
			user.Watchlist = watchlist
			if err = s.db.InsertWatchlist(user); err != nil {
				// Handle unique constraint violation separately
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					s.logger.Info("skipping duplicate watchlist entry",
						"user", user.LetterboxdUsername,
						"email", user.Email)
				} else {
					s.logger.Error("failed to update user watchlist",
						"user", user.LetterboxdUsername,
						"email", user.Email,
						"error", err)
				}
				continue
			}

			s.logger.Info("Successfully updated user watchlist",
				"user", user.LetterboxdUsername,
				"email", user.Email)
		}
	})
	if err != nil {
		log.Fatal("Error scheduling WatchlistScraper task", "error", err)
	}

	s.cron.Start()
}

func NewWatchlistScraper(db *db.Sqlite) *WatchlistScraper {
	return &WatchlistScraper{
		logger: slog.Default(),
		db:     db,
		cron:   cron.New(),
	}
}

func ScrapeWatchlist(letterboxdUsername string) ([]string, error) {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", letterboxdUsername)
	c := colly.NewCollector()
	var filmNames []string

	// Extracts film names from each poster-container
	c.OnHTML("li.poster-container", func(e *colly.HTMLElement) {
		e.DOM.Find("img").Each(func(i int, s *goquery.Selection) {
			filmName, exists := s.Attr("alt")
			if exists {
				filmNames = append(filmNames, filmName)
			}
		})
	})

	// Visits pagination links to collect films on other pages
	c.OnHTML("li.paginate-page a", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if strings.Contains(nextPage, "/watchlist/page/") {
			e.Request.Visit(nextPage)
		}
	})

	// Visit the initial page
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error visiting page: %w", err)
	}

	return filmNames, nil
}
