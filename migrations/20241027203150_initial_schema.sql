-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE film_event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    letterboxd_username TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
	watchlist TEXT[]
);

CREATE TABLE locations (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT UNIQUE NOT NULL
);


CREATE TABLE user_locations (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	location_id UUID REFERENCES locations(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, location_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_locations;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS film_event;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

