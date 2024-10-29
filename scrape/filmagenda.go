package scrape

import (
	"encoding/json"
	"fmt"
	"letterboxd-cineville/model"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Full struct for the JSON-LD "Event"

// ScrapeFilmEvents extracts FilmEvent details from the JSON-LD script
func ScrapeFilmEvents(e *colly.HTMLElement, filmEvents *[]model.FilmEvent) {
	jsonData := e.Text

	if strings.Contains(jsonData, "Event") {
		var rawEvent map[string]interface{}

		// Unmarshal the JSON into a map
		err := json.Unmarshal([]byte(jsonData), &rawEvent)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		// Parse the date strings to time.Time
		startDate, err := time.Parse(time.RFC3339, rawEvent["startDate"].(string))
		if err != nil {
			fmt.Println("Error parsing startDate:", err)
			return
		}

		endDate, err := time.Parse(time.RFC3339, rawEvent["endDate"].(string))
		if err != nil {
			fmt.Println("Error parsing endDate:", err)
			return
		}

		// Extract values directly
		event := model.FilmEvent{
			Name:            rawEvent["name"].(string),
			URL:             rawEvent["url"].(string),
			StartDate:       startDate,
			EndDate:         endDate,
			LocationName:    rawEvent["location"].(map[string]interface{})["name"].(string),
			LocationAddress: rawEvent["location"].(map[string]interface{})["address"].(string), // Adjust if necessary
			OrganizerName:   rawEvent["organizer"].(map[string]interface{})["name"].(string),
			OrganizerURL:    rawEvent["organizer"].(map[string]interface{})["url"].(string),
			PerformerName:   rawEvent["performer"].(map[string]interface{})["name"].(string),
		}

		*filmEvents = append(*filmEvents, event)
	}
}

func CollectFilmEvents(url string) ([]model.FilmEvent, error) {
	var filmEvents []model.FilmEvent
	c := colly.NewCollector(
		colly.AllowedDomains("www.filmvandaag.nl"),
	)
	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		ScrapeFilmEvents(e, &filmEvents)
	})
	err := c.Visit(url)
	return filmEvents, err
}
