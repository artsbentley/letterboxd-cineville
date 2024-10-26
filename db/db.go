package db

import (
	"context"
	"database/sql"
	"fmt"
	"letterboxd-cineville/model"
	"log"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Sql Sqlite

type Sqlite struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func (s *Sqlite) deleteExpiredFilmEvents() error {
	// TODO
	return nil
}

// TODO:
// needs to check if the values already exist or not
func (s *Sqlite) InsertFilmEvent(event model.FilmEvent) error {
	_, err := s.DB.NamedExec(`
		INSERT INTO film_event (
			name, 
			url, 
			start_date, 
			end_date, 
			location_name, 
			location_address, 
			organizer_name, 
			organizer_url, 
			performer_name
		) VALUES (
			:name, 
			:url, 
			:start_date, 
			:end_date, 
			:location_name, 
			:location_address, 
			:organizer_name, 
			:organizer_url, 
			:performer_name
		)`, event)
	if err != nil {
		s.Logger.Error("error inserting film event: ", "error: ", err)
		return err
	}
	s.Logger.Info("film event inserted successfully", "name: ", event.Name)
	return nil
}

func (s *Sqlite) InsertWatchlist(letterboxd model.Letterboxd) error {
	ctx := context.Background()

	// Begin a transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Check if the user exists by email
	var userID int
	err = tx.QueryRowContext(ctx, `SELECT id FROM letterboxd WHERE email = ?`, letterboxd.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If user does not exist, insert a new record
			result, err := tx.ExecContext(ctx, `INSERT INTO letterboxd (email, username) VALUES (?, ?)`, letterboxd.Email, letterboxd.Username)
			if err != nil {
				tx.Rollback()
				return err
			}
			lastID, err := result.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
			userID = int(lastID)
		} else {
			tx.Rollback()
			return err
		}
	}

	// Remove existing watchlist items not in the new list
	placeholders := strings.Repeat("?,", len(letterboxd.Watchlist))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	// Delete old watchlist items not in the new list
	deleteQuery := fmt.Sprintf(`
		DELETE FROM watchlist 
		WHERE letterboxd_id = ? AND film_title NOT IN (%s)
	`, placeholders)

	args := make([]interface{}, len(letterboxd.Watchlist)+1)
	args[0] = userID
	for i, film := range letterboxd.Watchlist {
		args[i+1] = film
	}

	_, err = tx.ExecContext(ctx, deleteQuery, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert new watchlist items, ignoring duplicates
	for _, film := range letterboxd.Watchlist {
		_, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO watchlist (letterboxd_id, film_title) VALUES (?, ?)`, userID, film)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	s.Logger.Info("Watchlist updated successfully", "email", letterboxd.Email)
	return nil
}

// init is automatically called when the package is imported
func init() {
	var err error

	// Assign to the global DB variable
	DB, err := sqlx.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Initialize the global Sql instance
	Sql = Sqlite{
		DB,
		slog.Default(),
	}

	fmt.Println("Database connection initialized.")
}
