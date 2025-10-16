-- name: AddUserToken :exec
UPDATE users
SET token = $2
WHERE id = $1;