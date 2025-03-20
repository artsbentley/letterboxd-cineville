-- name: CreateUser :one
INSERT INTO users (email, letterboxd_username)
VALUES ($1, $2)
RETURNING id, email, letterboxd_username, created_at, updated_at, watchlist;

-- name: GetUserByID :one
SELECT id, email, letterboxd_username, created_at, updated_at, watchlist
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, letterboxd_username, created_at, updated_at, watchlist
FROM users
WHERE email = $1;

-- name: GetUsers :many
SELECT id, email, letterboxd_username, created_at, updated_at, watchlist
FROM users;

-- name: UpdateUser :one
UPDATE users
SET
  email = $2,
  letterboxd_username = $3,
  watchlist = $4,
  updated_at = NOW()
WHERE id = $1
RETURNING id, email, letterboxd_username, created_at, updated_at, watchlist;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

