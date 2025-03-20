-- name: DeleteUserWatchlist :exec
UPDATE users 
SET watchlist = NULL
WHERE email = $1;

-- name: GetUserWatchlist :one
SELECT watchlist 
FROM users 
WHERE email = $1;

-- name: UpdateUserWatchlist :exec
UPDATE users 
SET watchlist = $2 
WHERE email = $1;

