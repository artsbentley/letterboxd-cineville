package scraper

import (
	"fmt"
	"letterboxd-cineville/internal/user"
	"letterboxd-cineville/internal/watchlist"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type WatchlistScraper struct {
	UserService      user.Provider
	WatchlistService watchlist.Provider
}

func NewWatchlistScraper(userService user.Provider, watchlistService watchlist.Provider) *WatchlistScraper {
	return &WatchlistScraper{
		UserService:      userService,
		WatchlistService: watchlistService,
	}
}

func (s *WatchlistScraper) Scrape() error {
	users, err := s.UserService.GetAllUsers()
	if err != nil {
		// slog.Error("failed to get users", "error", err)
		return err
	}

	for _, user := range users {
		watchlist, err := ScrapeUserWatchlist(user.LetterboxdUsername)
		if err != nil {
			slog.Warn("failed to scrape watchlist",
				"user", user.LetterboxdUsername,
				"email", user.Email,
				"error", err)
			continue
		}

		// Update user with new watchlist
		user.Watchlist = watchlist
		if err = s.WatchlistService.InsertWatchlist(user); err != nil {
			// Handle unique constraint violation separately
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				slog.Info("skipping duplicate watchlist entry",
					"user", user.LetterboxdUsername,
					"email", user.Email)
			} else {
				slog.Error("failed to update user watchlist",
					"user", user.LetterboxdUsername,
					"email", user.Email,
					"error", err)
			}
			continue
		}
		slog.Info("Successfully updated user watchlist",
			"user", user.LetterboxdUsername,
			"email", user.Email)
	}
	return nil
}

func ScrapeUserWatchlist(letterboxdUsername string) ([]string, error) {
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
