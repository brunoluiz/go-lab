-- name: SaveRadarQuadrant :one
INSERT INTO radar_quadrants (
  uniq_id,
  name
) VALUES ($1, $2)
ON CONFLICT (uniq_id) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING *;

-- name: DeleteRadarByRadarID :exec
DELETE FROM radar_quadrants
WHERE radar_id = $1;
