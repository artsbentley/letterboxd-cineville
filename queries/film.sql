-- name: CreateFilmEvent :one
INSERT INTO film_event (
    name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, name, url, start_date, end_date, location_name, location_address, city, organizer_name, organizer_url, performer_name;

-- name: MatchFilmEventsWithUser :many
select fe.*
from film_event fe
join users u on u.email = sqlc.arg(email)::string
join user_locations ul on u.id = ul.user_id
join locations l on ul.location_id = l.id
where exists (
    select 1
    from unnest(u.watchlist) as w
    where lower(w) = lower(fe.name)
)
and lower(fe.city) = lower(l.city);

-- name: GetFilmEventsByUserEmail :many
SELECT *
FROM film_event fe
JOIN users u ON fe.name = ANY(u.watchlist)
WHERE u.email = $1;

-- name: GetFilmEventByID :one
SELECT *
FROM film_event
WHERE id = $1;

-- name: ListFilmEvents :many
SELECT *
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



