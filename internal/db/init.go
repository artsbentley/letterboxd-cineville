package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Store *Queries

func Init() error {
	dbName := "postgres"
	dbHost := "localhost:5432"
	dbUser := "app"
	url := fmt.Sprintf("postgresql://%s@%s/%s", dbName, dbHost, dbUser)

	pgx, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to create DB connection: %w", err)
	}

	Store = New(pgx)
	slog.Info("database initialized succesfully")
	return nil
}
