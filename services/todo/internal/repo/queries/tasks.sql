-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE uniq_id = $1 LIMIT 1;

-- name: SaveTask :one
INSERT INTO tasks (
  uniq_id,
  task_uniq_id,
  title
) VALUES ($1, $2, $3)
ON CONFLICT (uniq_id) DO UPDATE
SET
  title = EXCLUDED.title
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE uniq_id = $1;

