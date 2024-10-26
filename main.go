package main

import (
	"fmt"
	database "letterboxd-cineville/db"
	"letterboxd-cineville/model"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// db := database.DB
	Sqlite := database.Sql
	var tables []string

	event := model.FilmEvent{
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(time.Hour * 2),
		Name:            "My Event",
		URL:             "https://example.com",
		LocationName:    "test",
		LocationAddress: "test",
		OrganizerName:   "hi there",
		OrganizerURL:    "fhuehf",
		PerformerName:   "hi",
	}

	err := Sqlite.InsertFilmEvent(event)
	if err != nil {
		log.Fatal(err)
	}

	lbox := model.Letterboxd{
		Email:     "arnoarts@hotmail.com",
		Username:  "Deltore",
		Watchlist: []string{"Banana", "phone"},
	}

	err = Sqlite.InsertWatchlist(lbox)
	if err != nil {
		log.Fatal(err)
	}

	// Query to get all tables from the sqlite_master system table
	query := `SELECT name 
				FROM sqlite_master 
				WHERE type='table' 
				AND name NOT LIKE 'sqlite_%';`
	err = Sqlite.DB.Select(&tables, query)
	if err != nil {
		log.Fatalf("Error querying tables: %v", err)
	}

	// Print all table names
	fmt.Println("Tables in the database:")
	for _, table := range tables {
		fmt.Println(table)
	}
}
