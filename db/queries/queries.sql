-- name: GetUserIDByEmail :one
SELECT id FROM "user" WHERE email = $1;

-- name: InsertUser :exec
INSERT INTO "user" (email, letterboxd_username, token) VALUES ($1, $2, $3);

-- name: DeleteUserWatchlist :exec
UPDATE "user" 
SET watchlist = NULL
WHERE email = $1;

-- name: GetUserWatchlist :one
SELECT watchlist FROM "user" WHERE email = $1;

-- name: UpdateUserWatchlist :exec
UPDATE "user" 
SET watchlist = $2 
WHERE email = $1;

-- name: UpdateUserEmailConfirmation :exec
UPDATE "user"
SET email_confirmation = $1
WHERE id = $2;

-- name: GetAllUsers :many
SELECT email, letterboxd_username FROM "user";

-- name: GetUserIDByToken :one
SELECT id FROM "user" WHERE token = $1;

-- name: InsertFilmEvent :exec
INSERT INTO "film_event" (
    name, url, start_date, end_date,
    location_name, location_address,
    organizer_name, organizer_url, performer_name
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8, $9
);

-- name: GetMatchingFilmEventsByEmail :many
SELECT fe.name, fe.url, fe.start_date, fe.end_date,
       fe.location_name, fe.location_address,
       fe.organizer_name, fe.organizer_url, fe.performer_name
FROM film_event AS fe
INNER JOIN "user" AS u ON fe.name = u.film_title
WHERE u.email = $1;
