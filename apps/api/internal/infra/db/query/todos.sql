-- name: ListTodos :many
SELECT id, title, description, status, due_date, created_at, updated_at
FROM todos
WHERE (sqlc.narg('status') IS NULL OR status = sqlc.narg('status'))
ORDER BY updated_at DESC
LIMIT sqlc.narg('limit');

-- name: CreateTodo :one
INSERT INTO todos (id, title, description, status, due_date, created_at, updated_at)
VALUES (?, ?, ?, 'open', ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, title, description, status, due_date, created_at, updated_at;

-- name: GetTodo :one
SELECT id, title, description, status, due_date, created_at, updated_at
FROM todos WHERE id = ?;

-- name: UpdateTodo :one
UPDATE todos
SET title = ?, description = ?, status = ?, due_date = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, title, description, status, due_date, created_at, updated_at;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;
