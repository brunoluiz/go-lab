-- name: GetLists :one
SELECT * FROM lists;

-- name: GetListByID :one
SELECT * FROM lists
WHERE uniq_id = $1 LIMIT 1;

-- name: SaveList :one
INSERT INTO lists (
  uniq_id,
  title
) VALUES ($1, $2)
ON CONFLICT (uniq_id) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING *;

-- name: DeleteList :exec
DELETE FROM lists
WHERE uniq_id = $1;
