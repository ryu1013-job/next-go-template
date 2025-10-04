package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ryu1013-job/next-go-template/apps/api/gen"
	"github.com/ryu1013-job/next-go-template/apps/api/internal/features/todo/usecase"
)

// TodoController implements the OpenAPI generated ServerInterface
type TodoController struct {
	todoUsecase usecase.TodoUsecase
}

// NewTodoController creates a new TodoController
func NewTodoController(todoUsecase usecase.TodoUsecase) *TodoController {
	return &TodoController{
		todoUsecase: todoUsecase,
	}
}

// ListTodos handles GET /todos
func (c *TodoController) ListTodos(w http.ResponseWriter, r *http.Request, params gen.ListTodosParams) {
	var status *string
	if params.Status != nil {
		s := string(*params.Status)
		status = &s
	}

	limit := int32(20) // default value
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}

	todos, err := c.todoUsecase.ListTodos(r.Context(), status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Items      []gen.Todo `json:"items"`
		NextCursor *string    `json:"nextCursor"`
	}{
		Items:      todos,
		NextCursor: nil, // TODO: Implement pagination if needed
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// CreateTodo handles POST /todos
func (c *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req gen.TodoCreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo, err := c.todoUsecase.CreateTodo(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(todo)
}

// GetTodo handles GET /todos/{id}
func (c *TodoController) GetTodo(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := c.todoUsecase.GetTodo(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

// UpdateTodo handles PATCH /todos/{id}
func (c *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request, id string) {
	var req gen.TodoUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	todo, err := c.todoUsecase.UpdateTodo(r.Context(), id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

// DeleteTodo handles DELETE /todos/{id}
func (c *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	if err := c.todoUsecase.DeleteTodo(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
