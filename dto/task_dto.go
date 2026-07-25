package dto

import (
	"errors"
	"time"
)

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=2"`
	Description string     `json:"description"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *string    `json:"assigned_to"` // string UUID user, opsional
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=2"`
	Description *string    `json:"description"`
	Status      *string    `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *string    `json:"assigned_to"`
}

// TaskFilter di-bind dari query params untuk GET /projects/:id/tasks
type TaskFilter struct {
	Status   string `form:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority string `form:"priority" binding:"omitempty,oneof=low medium high"`
	PaginationQuery
}

type TaskResponse struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *string    `json:"assigned_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

var (
	ErrTaskNotFound      = errors.New("task tidak ditemukan")
	ErrAssigneeNotFound  = errors.New("user yang di-assign tidak ditemukan")
	ErrInvalidTaskStatus = errors.New("status task tidak valid")
)
