-- name: GetRadars :many
SELECT sqlc.embed(r), sqlc.embed(ri), sqlc.embed(rq)
FROM radars r
JOIN radar_items ri ON ri.radar_id = r.id
JOIN radar_quadrants rq ON ri.quadrant_id = rq.id;

-- name: GetRadarByID :one
SELECT sqlc.embed(r), sqlc.embed(ri), sqlc.embed(rq)
FROM radars r
JOIN radar_items ri ON ri.radar_id = r.id
JOIN radar_quadrants rq ON ri.quadrant_id = rq.id
WHERE r.uniq_id = $1 LIMIT 1;

-- name: SaveRadar :one
INSERT INTO radars (
  uniq_id,
  title
) VALUES ($1, $2)
ON CONFLICT (uniq_id, quadrant_id) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING *;

-- name: DeleteRadar :exec
DELETE FROM radars
WHERE uniq_id = $1;
