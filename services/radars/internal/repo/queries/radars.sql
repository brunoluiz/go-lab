-- name: GetRadars :many
SELECT sqlc.embed(r), sqlc.embed(ri), sqlc.embed(rq)
FROM radars r
JOIN radar_items ri ON ri.radar_id = r.id
JOIN radar_quadrants rq ON ri.quadrant_id = rq.id
WHERE r.deleted_at IS NULL;

-- name: GetRadarByID :one
SELECT *
FROM radars r
WHERE
  1 = 1
  AND r.uniq_id = $1
  AND deleted_at IS NULL LIMIT 1;

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
UPDATE radars SET deleted_at = NOW()
WHERE uniq_id = $1;
