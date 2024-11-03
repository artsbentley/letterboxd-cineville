package db

import (
	"context"
	"database/sql"
	"fmt"
	sqlc "letterboxd-cineville/db/sqlc"
	"letterboxd-cineville/model"
	"log"
	"log/slog"

	"github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	queries *sqlc.Queries
	Logger  *slog.Logger
	DB      *sql.DB
}

func NewSqlite(db *sql.DB) (*Sqlite, error) {
	// Configure connection pool
	// db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	// db.SetMaxIdleConns(1)
	// db.SetConnMaxLifetime(time.Hour)

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=EXCLUSIVE"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &Sqlite{
		queries: sqlc.New(db),
		DB:      db,
		Logger:  slog.Default(),
	}, nil
}

func (s *Sqlite) CreateNewUser(user model.User) error {
	err := s.queries.InsertUser(context.Background(), sqlc.InsertUserParams{
		Email:              user.Email,
		LetterboxdUsername: user.LetterboxdUsername,
	})
	return err
}

func (s *Sqlite) ConfirmUserEmail(email string) error {
	err := s.queries.UpdateUserEmailConfirmation(context.Background(), sqlc.UpdateUserEmailConfirmationParams{
		Email:             email,
		EmailConfirmation: 1,
	})
	return err
}

func (s *Sqlite) InsertFilmEvent(event model.FilmEvent) error {
	err := s.queries.InsertFilmEvent(context.Background(), sqlc.InsertFilmEventParams{
		Name:            event.Name,
		Url:             event.URL,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		LocationName:    event.LocationName,
		LocationAddress: event.LocationAddress,
		OrganizerName:   event.OrganizerName,
		OrganizerUrl:    event.OrganizerURL,
		PerformerName:   event.PerformerName,
	})
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			s.Logger.Warn("film event already exists", "name", event.Name)
			return nil
		}
		return err
	}

	s.Logger.Info("film event inserted successfully", "name", event.Name)
	return nil
}

func (s *Sqlite) GetAllUsers() ([]model.User, error) {
	ctx := context.Background()
	rows, err := s.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for _, row := range rows {
		users = append(users, model.User{
			Email:              row.Email,
			LetterboxdUsername: row.LetterboxdUsername,
			Watchlist:          []string{}, // Initialize as empty or populate if needed
		})
	}

	s.Logger.Info("Retrieved all users successfully")
	return users, nil
}

func (s *Sqlite) GetOrCreateUserID(email, username string) (int64, error) {
	ctx := context.Background()

	// Attempt to retrieve the user ID by email
	userID, err := s.queries.GetUserIDByEmail(ctx, email)
	if err == sql.ErrNoRows {
		// User does not exist; insert a new user
		if err := s.queries.InsertUser(ctx, sqlc.InsertUserParams{
			Email:              email,
			LetterboxdUsername: username,
		}); err != nil {
			return 0, err
		}

		// Attempt to retrieve the user ID again after insertion
		userID, err = s.queries.GetUserIDByEmail(ctx, email)
		if err != nil {
			return 0, err
		}
		return userID, nil
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

// InsertWatchlist inserts a watchlist for the user.
func (s *Sqlite) InsertWatchlist(user model.User) error {
	ctx := context.Background()

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back if we return with error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	q := s.queries.WithTx(tx)

	userID, err := s.GetOrCreateUserID(user.Email, user.LetterboxdUsername)
	if err != nil {
		return fmt.Errorf("failed to get/create user ID: %w", err)
	}

	if err := q.DeleteUserWatchlist(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete existing watchlist: %w", err)
	}

	for _, film := range user.Watchlist {
		params := sqlc.InsertWatchlistItemParams{
			UserID:    userID,
			FilmTitle: film,
		}
		if err := q.InsertWatchlistItem(ctx, params); err != nil {
			return fmt.Errorf("failed to insert watchlist item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetMatchingFilmEventsByEmail retrieves film events based on the user's email.
func (s *Sqlite) GetMatchingFilmEventsByEmail(email string) ([]model.FilmEvent, error) {
	rows, err := s.queries.GetMatchingFilmEventsByEmail(context.Background(), email)
	if err != nil {
		s.Logger.Error("error retrieving matching film events", "error", err)
		return nil, err
	}

	var events []model.FilmEvent
	for _, row := range rows {
		events = append(events, model.FilmEvent{
			Name:            row.Name,
			URL:             row.Url,
			StartDate:       row.StartDate,
			EndDate:         row.EndDate,
			LocationName:    row.LocationName,
			LocationAddress: row.LocationAddress,
			OrganizerName:   row.OrganizerName,
			OrganizerURL:    row.OrganizerUrl,
			PerformerName:   row.PerformerName,
		})
	}

	s.Logger.Info("Matching film events retrieved successfully", "email", email)
	return events, nil
}

// Sql is a global variable for the Sqlite instance.
var Sql *Sqlite

// init initializes the database connection and the Sqlite instance.
func init() {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	Sql, err = NewSqlite(db)
	if err != nil {
		log.Fatal(err)
	}
}
