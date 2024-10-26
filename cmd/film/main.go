package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Full struct for the JSON-LD "Event"
type FilmEvent struct {
	Name            string
	URL             string
	StartDate       time.Time
	EndDate         time.Time
	LocationName    string
	LocationAddress string
	OrganizerName   string
	OrganizerURL    string
	PerformerName   string
}

// ScrapeFilmEvents extracts FilmEvent details from the JSON-LD script
func ScrapeFilmEvents(e *colly.HTMLElement, filmEvents *[]FilmEvent) {
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
		event := FilmEvent{
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

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.filmvandaag.nl"),
	)

	var filmEvents []FilmEvent

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		ScrapeFilmEvents(e, &filmEvents)
	})

	err := c.Visit("https://www.filmvandaag.nl/filmladder/stad/13-amsterdam")
	if err != nil {
		fmt.Println("Error visiting page:", err)
	}

	for _, event := range filmEvents {
		fmt.Printf("Event: %+v, Location: %+v, Organizer: %+v, Performer: %+v\n",
			event.Name, event.LocationName, event.OrganizerName, event.PerformerName)
	}
}
