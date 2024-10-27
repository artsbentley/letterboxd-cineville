package FilmMatcher

import (
	"letterboxd-cineville/db"
	"log/slog"
)

// NOTE: add username to this struct?
type FilmMatcher struct {
	DB     *db.Sqlite
	Logger *slog.Logger
}

func NewFilmMatcher(DB *db.Sqlite, Logger *slog.Logger) *FilmMatcher {
	return &FilmMatcher{
		DB:     DB,
		Logger: Logger,
	}
}

// func (m *FilmMatcher) GetMatches(email)
