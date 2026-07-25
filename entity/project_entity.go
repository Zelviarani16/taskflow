package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project merepresentasikan tabel "projects". Setiap project memiliki
// tepat satu pemilik (user yang membuatnya).
type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`

	Owner *User  `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Tasks []Task `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`

	TimeStamp
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
