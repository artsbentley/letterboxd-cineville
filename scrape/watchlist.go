package scrape

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func ScrapeWatchlist(username string) ([]string, error) {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", username)
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
