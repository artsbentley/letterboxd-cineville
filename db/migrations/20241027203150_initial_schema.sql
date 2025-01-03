-- +goose Up
-- +goose StatementBegin
CREATE TABLE film_event (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    location_name TEXT NOT NULL,
    location_address TEXT NOT NULL,
    organizer_name TEXT NOT NULL,
    organizer_url TEXT NOT NULL,
    performer_name TEXT NOT NULL,
    UNIQUE (name, start_date, location_name)
);

CREATE TABLE "user" (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    letterboxd_username TEXT NOT NULL,
    email_confirmation BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE watchlist (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    film_title TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    UNIQUE (user_id, film_title)  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS film_event;
DROP TABLE IF EXISTS watchlist;
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd

