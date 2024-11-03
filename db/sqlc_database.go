package db

import (
	"context"
	"database/sql"
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

func NewSqlite(db *sql.DB, logger *slog.Logger) *Sqlite {
	return &Sqlite{
		queries: sqlc.New(db),
		DB:      db,
		Logger:  logger,
	}
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
	userID, err := s.queries.GetOrCreateUserID(ctx, email)
	if err == sql.ErrNoRows {
		// User does not exist; insert a new user
		if err := s.queries.InsertUser(ctx, sqlc.InsertUserParams{
			Email:              email,
			LetterboxdUsername: username,
		}); err != nil {
			return 0, err
		}

		// Attempt to retrieve the user ID again after insertion
		userID, err = s.queries.GetOrCreateUserID(ctx, email)
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

	// Begin transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Use the transaction with Queries
	q := s.queries.WithTx(tx)

	// Get or create user ID
	userID, err := s.GetOrCreateUserID(user.Email, user.LetterboxdUsername)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Clear the current watchlist
	if err := q.DeleteUserWatchlist(ctx, userID); err != nil {
		tx.Rollback()
		return err
	}

	// Insert new watchlist items
	for _, film := range user.Watchlist {
		if err := q.InsertWatchlistItem(ctx, sqlc.InsertWatchlistItemParams{
			UserID:    userID,
			FilmTitle: film,
		}); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
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

	Sql = NewSqlite(db, slog.Default())
}
