package migrations

import (
	"errors"
	"log"

	"github.com/Zelviarani16/taskflow-api/entity"
	"github.com/Zelviarani16/taskflow-api/helpers"
	"gorm.io/gorm"
)

// Seed memasukkan data awal untuk development lokal: satu akun admin
// agar kamu punya sesuatu untuk login setelah --migrate.
// Aman dijalankan lebih dari sekali - dia melewati pembuatan jika admin
// sudah ada.
func Seed(db *gorm.DB) error {
	var existing entity.User
	err := db.Where("email = ?", "admin@taskflow.local").Take(&existing).Error
	if err == nil {
		log.Println("admin sudah ada, skip seeding")
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashed, err := helpers.HashPassword("admin12345")
	if err != nil {
		return err
	}

	admin := entity.User{
		Name:     "Admin TaskFlow",
		Email:    "admin@taskflow.local",
		Password: hashed,
		Role:     entity.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("admin default dibuat: admin@taskflow.local / admin12345")
	return nil
}
