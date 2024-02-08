-- name: SaveRadarItem :one
INSERT INTO radar_items (
  uniq_id,
  radar_id,
  quadrant_id,
  name,
  description
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (uniq_id) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description
RETURNING *;

