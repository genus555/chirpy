-- name: GetChirpByChirpId :one
SELECT * FROM chirps
WHERE id = $1;