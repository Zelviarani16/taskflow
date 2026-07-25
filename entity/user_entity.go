package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the "users" table.
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	Role     Role      `gorm:"type:varchar(20);not null;default:'member'" json:"role"`

	// Relasi - seorang user bisa memiliki banyak project dan ditugaskan banyak task.
	Projects []Project `gorm:"foreignKey:OwnerID" json:"-"`
	Tasks    []Task    `gorm:"foreignKey:AssignedTo" json:"-"`

	TimeStamp
}

// BeforeCreate adalah hook GORM yang berjalan sebelum INSERT, apapun
// jalur kode yang membuat user (register, seeder, dll). Menaruh
// pembuatan UUID dan role default di sini berarti setiap jalur pembuatan
// mendapatkannya secara otomatis tanpa perlu diingat oleh pemanggil.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if !IsValidRole(u.Role) {
		u.Role = RoleMember
	}

	return nil
}
