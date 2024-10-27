-- +goose Up
-- name: CreateFilmEventTable
CREATE TABLE film_event (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    location_name TEXT NOT NULL,
    location_address TEXT NOT NULL,
    organizer_name TEXT NOT NULL,
    organizer_url TEXT NOT NULL,
    performer_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE film_event;

