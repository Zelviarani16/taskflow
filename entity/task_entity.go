package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Task merepresentasikan tabel "tasks". Setiap task milik tepat satu
// project dan bisa secara opsional ditugaskan ke seorang user.
type Task struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID   uuid.UUID    `gorm:"type:uuid;not null;index" json:"project_id"`
	Title       string       `gorm:"not null" json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `gorm:"type:varchar(20);not null;default:'todo'" json:"status"`
	Priority    TaskPriority `gorm:"type:varchar(20);not null;default:'medium'" json:"priority"`
	DueDate     *time.Time   `json:"due_date"`

	// AssignedTo adalah FK yang bisa null - task bisa tidak ditugaskan ke siapapun.
	AssignedTo *uuid.UUID `gorm:"type:uuid;index" json:"assigned_to"`
	Assignee   *User      `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`

	TimeStamp
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if !IsValidStatus(t.Status) {
		t.Status = StatusTodo
	}
	if !IsValidPriority(t.Priority) {
		t.Priority = PriorityMedium
	}
	return nil
}
