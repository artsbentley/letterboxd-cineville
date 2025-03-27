package scraper

import (
	"encoding/json"
	"fmt"
	"letterboxd-cineville/internal/model"
	"letterboxd-cineville/internal/service"
	"letterboxd-cineville/internal/types"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type FilmEventScraper struct {
	UserService      service.UserProvider
	FilmEventService service.FilmEventProvider
}

func NewFilmEventScraper(userService service.UserProvider, filmEventService service.FilmEventProvider) *FilmEventScraper {
	return &FilmEventScraper{
		UserService:      userService,
		FilmEventService: filmEventService,
	}
}

func (s *FilmEventScraper) Scrape() error {
	for _, city := range types.Cities {
		url := s.constructURL(city.URLPartial)
		filmEvents, err := CollectFilmEvents(url)
		if err != nil {
			log.Fatal(err)
		}
		// NOTE: this will be unnecessary if we end up using locationaddress
		// value in the database to retrieve the city instead
		// set the city location for the film events
		for i := range filmEvents {
			filmEvents[i].City = city.Name
		}
		s.processFilmEvents(filmEvents)
	}
	return nil
}

func (s *FilmEventScraper) constructURL(location string) string {
	baseURL := "https://www.filmvandaag.nl/filmladder/stad/"
	return fmt.Sprintf("%s%s", baseURL, location)
}

func (s *FilmEventScraper) processFilmEvents(filmEvents []model.FilmEvent) {
	for _, event := range filmEvents {
		err := s.FilmEventService.InsertFilmEvent(event)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			} else {
				slog.Error("Failed to insert FilmEvent", "error", err)
			}
		} else {
			slog.Info("Successfully inserted FilmEvent", "event", event.Name)
		}
	}
}

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
