-- name: ByUID :one
SELECT * FROM lists
WHERE uid = $1 LIMIT 1;

-- name: Create :one
INSERT INTO lists (
  uid,
  title
) VALUES ($1, $2)
RETURNING *;
