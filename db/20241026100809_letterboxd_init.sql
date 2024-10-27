-- +goose Up
-- +goose StatementBegin
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
DROP TABLE watchlist;
DROP TABLE letterboxd;
-- +goose StatementEnd
