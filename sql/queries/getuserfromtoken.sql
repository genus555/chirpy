-- name: GetUserFromToken :one
SELECT * FROM users
WHERE token = $1;