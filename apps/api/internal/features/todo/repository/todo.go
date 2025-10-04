package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryu1013-job/next-go-template/apps/api/gen"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/infra/db"
)

// TodoRepository defines the interface for todo repository operations
type TodoRepository interface {
	List(ctx context.Context, status *string, limit int32) ([]gen.Todo, error)
	Get(ctx context.Context, id string) (gen.Todo, error)
	Create(ctx context.Context, id, title string, description *string, dueDate *time.Time) (gen.Todo, error)
	Update(ctx context.Context, id, title string, description *string, status string, dueDate *time.Time) (gen.Todo, error)
	Delete(ctx context.Context, id string) error
}

// todoRepository implements TodoRepository
type todoRepository struct {
	queries *db.Queries
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(queries *db.Queries) TodoRepository {
	return &todoRepository{
		queries: queries,
	}
}

// toGen converts db.Todo to gen.Todo
func (r *todoRepository) toGen(dbTodo db.Todo) gen.Todo {
	var description *string
	if dbTodo.Description.Valid {
		description = &dbTodo.Description.String
	}

	return gen.Todo{
		Id:          dbTodo.ID,
		Title:       dbTodo.Title,
		Description: description,
		Status:      gen.TodoStatus(dbTodo.Status),
		DueDate:     r.toTimePtr(dbTodo.DueDate),
		CreatedAt:   r.mustParse(dbTodo.CreatedAt),
		UpdatedAt:   r.mustParse(dbTodo.UpdatedAt),
	}
}

// toTimePtr converts sql.NullString to *time.Time
func (r *todoRepository) toTimePtr(ns sql.NullString) *time.Time {
	if !ns.Valid || ns.String == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, ns.String)
	if err != nil {
		return nil
	}
	return &t
}

// mustParse parses time string
func (r *todoRepository) mustParse(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// List retrieves todos with optional status filter and limit
func (r *todoRepository) List(ctx context.Context, status *string, limit int32) ([]gen.Todo, error) {
	var st sql.NullString
	if status != nil {
		st = sql.NullString{String: *status, Valid: true}
	}

	rows, err := r.queries.ListTodos(ctx, db.ListTodosParams{
		Status: st,
		Limit:  sql.NullInt64{Int64: int64(limit), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	todos := make([]gen.Todo, 0, len(rows))
	for _, row := range rows {
		todos = append(todos, r.toGen(row))
	}

	return todos, nil
}

// Get retrieves a todo by ID
func (r *todoRepository) Get(ctx context.Context, id string) (gen.Todo, error) {
	row, err := r.queries.GetTodo(ctx, id)
	if err != nil {
		return gen.Todo{}, err
	}
	return r.toGen(row), nil
}

// Create creates a new todo
func (r *todoRepository) Create(ctx context.Context, id, title string, description *string, dueDate *time.Time) (gen.Todo, error) {
	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}

	var due sql.NullString
	if dueDate != nil {
		due = sql.NullString{String: dueDate.Format(time.RFC3339), Valid: true}
	}

	row, err := r.queries.CreateTodo(ctx, db.CreateTodoParams{
		ID:          id,
		Title:       title,
		Description: desc,
		DueDate:     due,
	})
	if err != nil {
		return gen.Todo{}, err
	}

	return r.toGen(row), nil
}

// Update updates an existing todo
func (r *todoRepository) Update(ctx context.Context, id, title string, description *string, status string, dueDate *time.Time) (gen.Todo, error) {
	var desc sql.NullString
	if description != nil {
		desc = sql.NullString{String: *description, Valid: true}
	}

	var due sql.NullString
	if dueDate != nil {
		due = sql.NullString{String: dueDate.Format(time.RFC3339), Valid: true}
	}

	row, err := r.queries.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:          id,
		Title:       title,
		Description: desc,
		Status:      status,
		DueDate:     due,
	})
	if err != nil {
		return gen.Todo{}, err
	}

	return r.toGen(row), nil
}

// Delete deletes a todo by ID
func (r *todoRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteTodo(ctx, id)
}
