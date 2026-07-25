package migrations

import (
	"github.com/Zelviarani16/taskflow-api/entity"
	"gorm.io/gorm"
)

// Migrate membuat/memperbarui tabel untuk setiap entity yang terdaftar.
// Urutan penting: tabel induk (tanpa FK) harus ada terlebih dahulu
// sebelum tabel anak yang mereferensikannya.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Project{},
		&entity.Task{},
	)
}
