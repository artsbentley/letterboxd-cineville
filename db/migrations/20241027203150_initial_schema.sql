-- +goose Up
-- +goose StatementBegin
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
    performer_name TEXT NOT NULL,
	UNIQUE (name, start_date, location_name)
);

CREATE TABLE letterboxd (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL
);

CREATE TABLE watchlist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    letterboxd_id INTEGER NOT NULL,
    film_title TEXT NOT NULL,
    FOREIGN KEY (letterboxd_id) REFERENCES letterboxd(id),
    UNIQUE (letterboxd_id, film_title)  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE film_event;
DROP TABLE watchlist;
DROP TABLE letterboxd;
-- +goose StatementEnd
