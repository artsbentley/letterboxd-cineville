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
	"github.com/mattn/go-sqlite3"
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
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				s.Logger.Warn("film event already exists", "name: ", event.Name)
				// NOTE: better to return error?
				return nil
			}
		}
		return err // Return other errors
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

func (s *Sqlite) GetMatchingFilmEventsByEmail(email string) ([]model.FilmEvent, error) {
	var filmEvents []model.FilmEvent

	query := `
		SELECT fe.name, fe.url, fe.start_date, fe.end_date, 
		       fe.location_name, fe.location_address, 
		       fe.organizer_name, fe.organizer_url, 
		       fe.performer_name
		FROM film_event AS fe
		INNER JOIN watchlist AS wl ON fe.name = wl.film_title
		INNER JOIN letterboxd AS lb ON lb.id = wl.letterboxd_id
		WHERE lb.email = ?
	`

	err := s.DB.Select(&filmEvents, query, email)
	if err != nil {
		s.Logger.Error("error retrieving matching film events", "error", err)
		return nil, err
	}

	s.Logger.Info("Matching film events retrieved successfully", "email", email)
	return filmEvents, nil
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
}
