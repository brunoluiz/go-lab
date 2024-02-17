-- name: SaveRadarQuadrant :one
INSERT INTO radar_quadrants (
  uniq_id,
  radar_id,
  name
) VALUES ($1, $2, $3)
ON CONFLICT (uniq_id) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING *;

-- name: DeleteRadarQuadrantByID :exec
DELETE FROM radar_quadrants
WHERE radar_id = $1;
