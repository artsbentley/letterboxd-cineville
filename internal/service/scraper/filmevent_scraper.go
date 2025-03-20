package scraper

import (
	"encoding/json"
	"fmt"
	"letterboxd-cineville/internal/model"
	"letterboxd-cineville/internal/service"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/lmittmann/tint"
)

type FilmEventScraper struct {
	UserService      service.UserProvider
	FilmEventService service.FilmEventProvider
	Logger           *slog.Logger
}

func NewFilmEventScraper(userService service.UserProvider, filmEventService service.FilmEventProvider) *FilmEventScraper {
	return &FilmEventScraper{
		UserService:      userService,
		FilmEventService: filmEventService,
		Logger:           slog.New(tint.NewHandler(os.Stderr, nil)),
	}
}

// TODO: implement every city
func (s *FilmEventScraper) Scrape() error {
	filmEvents, err := CollectFilmEvents("https://www.filmvandaag.nl/filmladder/stad/13-amsterdam")
	// filmEvents, err := CollectFilmEvents("https://www.filmvandaag.nl/filmladder/stad/159-rotterdam")
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range filmEvents {
		err := s.FilmEventService.InsertFilmEvent(event)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			} else {
				s.Logger.Error("Failed to insert FilmEvent", "error", err)
			}
		} else {
			s.Logger.Info("Successfully inserted FilmEvent", "event", event.Name)
		}
	}
	return nil
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
