-- name: ByID :one
SELECT * FROM lists
WHERE id = $1 LIMIT 1;
