package db

import (
	"context"
	"fmt"
	sqlc "letterboxd-cineville/db/sqlc"
	"letterboxd-cineville/model"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	queries *sqlc.Queries
	Logger  *slog.Logger
	DB      *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) (*Store, error) {
	return &Store{
		queries: sqlc.New(db),
		DB:      db,
		Logger:  slog.Default(),
	}, nil
}

// NOTE: do i really want token in the business logic model?
func (s *Store) CreateNewUser(user model.User) error {
	err := s.queries.InsertUser(context.Background(), sqlc.InsertUserParams{
		Email:              user.Email,
		LetterboxdUsername: user.LetterboxdUsername,
		Token:              user.Token,
	})
	return err
}

func (s *Store) ConfirmUserEmail(id int) error {
	err := s.queries.UpdateUserEmailConfirmation(context.Background(), sqlc.UpdateUserEmailConfirmationParams{
		ID:                int64(id),
		EmailConfirmation: true,
	})
	return err
}

func (s *Store) GetUserIDByToken(token string) (int, error) {
	id, err := s.queries.GetUserIDByToken(context.Background(), token)
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (s *Store) InsertFilmEvent(event model.FilmEvent) error {
	// Convert time.Time to pgtype.Timestamptz
	startDate := pgtype.Timestamptz{
		Time:  event.StartDate,
		Valid: true,
	}
	endDate := pgtype.Timestamptz{
		Time:  event.EndDate,
		Valid: true,
	}

	err := s.queries.InsertFilmEvent(context.Background(), sqlc.InsertFilmEventParams{
		Name:            event.Name,
		Url:             event.URL,
		StartDate:       startDate,
		EndDate:         endDate,
		LocationName:    event.LocationName,
		LocationAddress: event.LocationAddress,
		OrganizerName:   event.OrganizerName,
		OrganizerUrl:    event.OrganizerURL,
		PerformerName:   event.PerformerName,
	})
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok {
			// 23505 is the PostgreSQL error code for unique_violation
			if pgErr.Code == "23505" {
				s.Logger.Warn("film event already exists", "name", event.Name)
				return nil
			}
		}
		return fmt.Errorf("failed to insert film event: %w", err)
	}

	s.Logger.Info("film event inserted successfully", "name", event.Name)
	return nil
}

func (s *Store) GetAllUsers() ([]model.User, error) {
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
			Watchlist:          []string{},
		})
	}

	s.Logger.Info("Retrieved all users successfully")
	return users, nil
}

func (s *Store) GetUserID(email, username string) (int64, error) {
	ctx := context.Background()
	userID, err := s.queries.GetUserIDByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *Store) InsertWatchlist(user model.User) error {
	ctx := context.Background()

	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	q := s.queries.WithTx(tx)

	if err := q.DeleteUserWatchlist(ctx, user.Email); err != nil {
		return fmt.Errorf("failed to delete existing watchlist: %w", err)
	}

	userWatchlistParams := sqlc.UpdateUserWatchlistParams{
		Email:     user.Email,
		Watchlist: user.Watchlist,
	}

	if err := q.UpdateUserWatchlist(ctx, userWatchlistParams); err != nil {
		return fmt.Errorf("failed to insert watchlist for user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *Store) GetMatchingFilmEventsByEmail(email string) ([]model.FilmEvent, error) {
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
			StartDate:       row.StartDate.Time,
			EndDate:         row.EndDate.Time,
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

var Sql *Store

func init() {
	var (
		dbUser = os.Getenv("USER")
		dbHost = "localhost:5432"
		url    = fmt.Sprintf("postgresql://%s@%s/%s", dbUser, dbHost, dbUser)
	)
	fmt.Println(url)

	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create a connection pool: %v\n", err)
		os.Exit(1)
	}

	Sql, err = NewStore(conn)
	if err != nil {
		log.Fatal(err)
	}
}
