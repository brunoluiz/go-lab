// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: radar_quadrants.sql

package repo

import (
	"context"
)

const deleteRadarQuadrantByID = `-- name: DeleteRadarQuadrantByID :exec
DELETE FROM radar_quadrants
WHERE radar_id = $1
`

func (q *Queries) DeleteRadarQuadrantByID(ctx context.Context, radarID int32) error {
	_, err := q.db.ExecContext(ctx, deleteRadarQuadrantByID, radarID)
	return err
}

const getRadarQuadrantByUniqID = `-- name: GetRadarQuadrantByUniqID :one
SELECT id, uniq_id, radar_id, name
FROM radar_quadrants
WHERE uniq_id = $1
`

func (q *Queries) GetRadarQuadrantByUniqID(ctx context.Context, uniqID string) (RadarQuadrant, error) {
	row := q.db.QueryRowContext(ctx, getRadarQuadrantByUniqID, uniqID)
	var i RadarQuadrant
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.RadarID,
		&i.Name,
	)
	return i, err
}

const getRadarQuadrantsByRadarID = `-- name: GetRadarQuadrantsByRadarID :many
SELECT id, uniq_id, radar_id, name
FROM radar_quadrants
WHERE radar_id = $1
`

func (q *Queries) GetRadarQuadrantsByRadarID(ctx context.Context, radarID int32) ([]RadarQuadrant, error) {
	rows, err := q.db.QueryContext(ctx, getRadarQuadrantsByRadarID, radarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RadarQuadrant
	for rows.Next() {
		var i RadarQuadrant
		if err := rows.Scan(
			&i.ID,
			&i.UniqID,
			&i.RadarID,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveRadarQuadrant = `-- name: SaveRadarQuadrant :one
INSERT INTO radar_quadrants (
  uniq_id,
  radar_id,
  name
) VALUES ($1, $2, $3)
ON CONFLICT (uniq_id) DO UPDATE
SET
  name = EXCLUDED.name
RETURNING id, uniq_id, radar_id, name
`

type SaveRadarQuadrantParams struct {
	UniqID  string `json:"uniq_id"`
	RadarID int32  `json:"radar_id"`
	Name    string `json:"name"`
}

func (q *Queries) SaveRadarQuadrant(ctx context.Context, arg SaveRadarQuadrantParams) (RadarQuadrant, error) {
	row := q.db.QueryRowContext(ctx, saveRadarQuadrant, arg.UniqID, arg.RadarID, arg.Name)
	var i RadarQuadrant
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.RadarID,
		&i.Name,
	)
	return i, err
}
