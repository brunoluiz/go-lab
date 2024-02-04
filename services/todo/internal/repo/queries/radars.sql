-- name: GetRadars :many
SELECT * FROM radars;

-- name: GetRadarByID :one
SELECT * FROM radars
WHERE uniq_id = $1 LIMIT 1;

-- name: SaveRadar :one
INSERT INTO radars (
  uniq_id,
  title
) VALUES ($1, $2)
ON CONFLICT (uniq_id) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING *;

-- name: DeleteRadar :exec
DELETE FROM radars
WHERE uniq_id = $1;
