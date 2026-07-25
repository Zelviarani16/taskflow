package migrations

import (
	"github.com/Zelviarani16/taskflow-api/entity"
	"gorm.io/gorm"
)

// Rollback menghapus semua tabel. Urutannya kebalikan dari Migrate: tabel
// anak (yang punya FK) dihapus duluan agar tidak terjadi pelanggaran constraint.
func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		&entity.Task{},
		&entity.Project{},
		&entity.User{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}
	return nil
}
