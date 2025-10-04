package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryu1013-job/next-go-template/apps/api/gen"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/features/todo/repository"
)

// TodoUsecase defines the interface for todo business logic operations
type TodoUsecase interface {
	ListTodos(ctx context.Context, status *string, limit int32) ([]gen.Todo, error)
	GetTodo(ctx context.Context, id string) (gen.Todo, error)
	CreateTodo(ctx context.Context, req gen.TodoCreateInput) (gen.Todo, error)
	UpdateTodo(ctx context.Context, id string, req gen.TodoUpdateInput) (gen.Todo, error)
	DeleteTodo(ctx context.Context, id string) error
}

// todoUsecase implements TodoUsecase
type todoUsecase struct {
	todoRepo repository.TodoRepository
}

// NewTodoUsecase creates a new TodoUsecase
func NewTodoUsecase(todoRepo repository.TodoRepository) TodoUsecase {
	return &todoUsecase{
		todoRepo: todoRepo,
	}
}

// ListTodos retrieves a list of todos with optional filters
func (u *todoUsecase) ListTodos(ctx context.Context, status *string, limit int32) ([]gen.Todo, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	return u.todoRepo.List(ctx, status, limit)
}

// GetTodo retrieves a single todo by ID
func (u *todoUsecase) GetTodo(ctx context.Context, id string) (gen.Todo, error) {
	return u.todoRepo.Get(ctx, id)
}

// CreateTodo creates a new todo
func (u *todoUsecase) CreateTodo(ctx context.Context, req gen.TodoCreateInput) (gen.Todo, error) {
	// Generate unique ID
	id := uuid.New().String()

	// Create todo with default status "open"
	return u.todoRepo.Create(ctx, id, req.Title, req.Description, req.DueDate)
}

// UpdateTodo updates an existing todo
func (u *todoUsecase) UpdateTodo(ctx context.Context, id string, req gen.TodoUpdateInput) (gen.Todo, error) {
	// Get current todo to merge fields
	current, err := u.todoRepo.Get(ctx, id)
	if err != nil {
		return gen.Todo{}, err
	}

	// Merge update fields with current values
	title := current.Title
	if req.Title != nil {
		title = *req.Title
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	} else {
		description = current.Description
	}

	status := string(current.Status)
	if req.Status != nil {
		status = string(*req.Status)
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		dueDate = req.DueDate
	} else {
		dueDate = current.DueDate
	}

	return u.todoRepo.Update(ctx, id, title, description, status, dueDate)
}

// DeleteTodo deletes a todo by ID
func (u *todoUsecase) DeleteTodo(ctx context.Context, id string) error {
	return u.todoRepo.Delete(ctx, id)
}
