package dto

import (
	"errors"
	"time"
)

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2"`
	Description *string `json:"description"`
}

type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	TaskCount   int64     `json:"task_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrProjectNotFound = errors.New("project tidak ditemukan")
	ErrProjectForbidden = errors.New("kamu bukan pemilik project ini")
)
