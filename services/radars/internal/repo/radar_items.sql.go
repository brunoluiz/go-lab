// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: radar_items.sql

package repo

import (
	"context"
)

const saveRadarItem = `-- name: SaveRadarItem :one
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
RETURNING id, uniq_id, radar_id, quadrant_id, name, description, updated_at, created_at
`

type SaveRadarItemParams struct {
	UniqID      string `json:"uniq_id"`
	RadarID     int32  `json:"radar_id"`
	QuadrantID  int32  `json:"quadrant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) SaveRadarItem(ctx context.Context, arg SaveRadarItemParams) (RadarItem, error) {
	row := q.db.QueryRowContext(ctx, saveRadarItem,
		arg.UniqID,
		arg.RadarID,
		arg.QuadrantID,
		arg.Name,
		arg.Description,
	)
	var i RadarItem
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.RadarID,
		&i.QuadrantID,
		&i.Name,
		&i.Description,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
