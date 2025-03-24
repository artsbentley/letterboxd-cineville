-- name: CreateLocation :one
INSERT INTO locations (city)
VALUES ($1)
ON CONFLICT (city) DO UPDATE SET city = EXCLUDED.city
RETURNING id;

-- name: GetLocationByCity :one
SELECT id FROM locations WHERE city = $1;

-- name: AssignUserLocation :exec
INSERT INTO user_locations (user_id, location_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetUserLocationCities :many
SELECT l.city
FROM locations l
JOIN user_locations ul ON l.id = ul.location_id
WHERE ul.user_id = $1;
