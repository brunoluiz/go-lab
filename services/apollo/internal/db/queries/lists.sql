-- name: GetListByID :one
SELECT * FROM lists
WHERE uid = $1 LIMIT 1;

-- name: SaveList :one
INSERT INTO lists (
  uid,
  title
) VALUES ($1, $2)
ON CONFLICT (uid) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING *;

-- name: DeleteList :exec
DELETE FROM lists
WHERE uid = $1;
