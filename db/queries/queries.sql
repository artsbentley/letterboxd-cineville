-- name: GetOrCreateUserID :one
SELECT id FROM user WHERE email = ?;

-- name: InsertUser :exec
INSERT INTO user (email, letterboxd_username) VALUES (?, ?);

-- name: DeleteUserWatchlist :exec
DELETE FROM watchlist WHERE user_id = ?;

-- name: InsertWatchlistItem :exec
INSERT INTO watchlist (user_id, film_title) VALUES (?, ?);

-- name: GetAllUsers :many
SELECT email, letterboxd_username FROM user;

-- name: InsertFilmEvent :exec
INSERT INTO film_event (
    name, url, start_date, end_date,
    location_name, location_address,
    organizer_name, organizer_url, performer_name
) VALUES (
    ?, ?, ?, ?,
    ?, ?, ?, ?, ?
);

-- name: GetMatchingFilmEventsByEmail :many
SELECT fe.name, fe.url, fe.start_date, fe.end_date,
       fe.location_name, fe.location_address,
       fe.organizer_name, fe.organizer_url, fe.performer_name
FROM film_event AS fe
INNER JOIN watchlist AS wl ON fe.name = wl.film_title
INNER JOIN user ON user.id = wl.user_id
WHERE user.email = ?;


