// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: lists.sql

package repo

import (
	"context"
)

const deleteList = `-- name: DeleteList :exec
DELETE FROM lists
WHERE uniq_id = $1
`

func (q *Queries) DeleteList(ctx context.Context, uniqID string) error {
	_, err := q.db.ExecContext(ctx, deleteList, uniqID)
	return err
}

const getListByID = `-- name: GetListByID :one
SELECT id, uniq_id, title, position, updated_at, created_at FROM lists
WHERE uniq_id = $1 LIMIT 1
`

func (q *Queries) GetListByID(ctx context.Context, uniqID string) (List, error) {
	row := q.db.QueryRowContext(ctx, getListByID, uniqID)
	var i List
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.Title,
		&i.Position,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getLists = `-- name: GetLists :one
SELECT id, uniq_id, title, position, updated_at, created_at FROM lists
`

func (q *Queries) GetLists(ctx context.Context) (List, error) {
	row := q.db.QueryRowContext(ctx, getLists)
	var i List
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.Title,
		&i.Position,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const saveList = `-- name: SaveList :one
INSERT INTO lists (
  uniq_id,
  title
) VALUES ($1, $2)
ON CONFLICT (uniq_id) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING id, uniq_id, title, position, updated_at, created_at
`

type SaveListParams struct {
	UniqID string `json:"uniq_id"`
	Title  string `json:"title"`
}

func (q *Queries) SaveList(ctx context.Context, arg SaveListParams) (List, error) {
	row := q.db.QueryRowContext(ctx, saveList, arg.UniqID, arg.Title)
	var i List
	err := row.Scan(
		&i.ID,
		&i.UniqID,
		&i.Title,
		&i.Position,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
