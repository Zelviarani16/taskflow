package entity

import (
	"time"

	"gorm.io/gorm"
)

// TimeStamp ditanamkan ke setiap entity agar semua tabel memiliki
// kolom created_at / updated_at / deleted_at (soft delete) yang sama.
type TimeStamp struct {
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// Role adalah tipe string kustom agar role yang tidak valid gagal saat
// kompilasi, bukan diam-diam menerima string yang salah ketik.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func IsValidRole(r Role) bool {
	return r == RoleAdmin || r == RoleMember
}

// TaskStatus melacak posisi task dalam siklus hidupnya.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func IsValidStatus(s TaskStatus) bool {
	return s == StatusTodo || s == StatusInProgress || s == StatusDone
}

// TaskPriority menentukan tingkat urgensi task.
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

func IsValidPriority(p TaskPriority) bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh
}
