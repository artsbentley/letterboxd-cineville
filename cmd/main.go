package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	var filmNames []string

	c.OnHTML("li.poster-container", func(e *colly.HTMLElement) {
		e.DOM.Find("img").Each(func(i int, s *goquery.Selection) {
			filmName, exists := s.Attr("alt")
			if exists {
				filmNames = append(filmNames, filmName)
			}
		})
	})

	c.OnHTML("li.paginate-page a", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if strings.Contains(nextPage, "/watchlist/page/") {
			e.Request.Visit(nextPage)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("counted %d films\n", len(filmNames))
		fmt.Println("Watchlist Film Names:")
		for _, name := range filmNames {
			fmt.Println(name)
		}
	})

	err := c.Visit("https://letterboxd.com/deltore/watchlist/")
	if err != nil {
		fmt.Println("Error visiting page:", err)
	}
}
