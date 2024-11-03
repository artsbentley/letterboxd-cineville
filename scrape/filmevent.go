package scrape

import (
	"encoding/json"
	"fmt"
	"letterboxd-cineville/db"
	"letterboxd-cineville/model"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

func (s *FilmEventScraper) Scrape() {
	// _, err := s.cron.AddFunc("0 2 * * 0", func() {
	_, err := s.cron.AddFunc("* * * * *", func() {
		s.logger.Info("FilmEventScraper Scheduled task running...")
		filmEvents, err := CollectFilmEvents("https://www.filmvandaag.nl/filmladder/stad/13-amsterdam")
		if err != nil {
			log.Fatal(err)
		}
		for _, event := range filmEvents {
			err := s.db.InsertFilmEvent(event)
			if err != nil {
				s.logger.Error("Failed to insert FilmEvent: ", "error", err)
			}
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling FilmEventScraper task: %v", err)
	}
	s.cron.Start() // Start the cron scheduler
}

func NewFilmEventScraper(db *db.Sqlite) *FilmEventScraper {
	return &FilmEventScraper{
		logger: slog.Default(),
		db:     db,
		cron:   cron.New(),
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
