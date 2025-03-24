-- name: CreateFilmEvent :one
INSERT INTO film_event (
    name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name;

-- name: GetFilmEventsByUserEmail :many
SELECT fe.id, fe.name, fe.url, fe.start_date, fe.end_date, fe.location_name, fe.location_address, fe.city, fe.organizer_name, fe.organizer_url, fe.performer_name
FROM film_event fe
JOIN users u ON fe.name = ANY(u.watchlist)
WHERE u.email = $1;

-- name: GetFilmEventByID :one
SELECT id, name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name
FROM film_event
WHERE id = $1;

-- name: ListFilmEvents :many
SELECT id, name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name
FROM film_event;

-- name: DeleteFilmEvent :exec
DELETE FROM film_event
WHERE id = $1;

-- name: DeletePastFilmEvents :exec
DELETE FROM film_event
WHERE start_date < NOW();

-- name: GetFilmEventByCity :many
SELECT *
FROM film_event
WHERE TRIM(LOWER(SUBSTRING_INDEX(address, ',', -1))) = LOWER($1);



