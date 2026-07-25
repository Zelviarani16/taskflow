package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetUpPostgreSQLConnection membuka koneksi DB satu kali saat startup.
// Secara lokal membaca kredensial dari .env; di production, platform
// (Railway, Render, dll.) diharapkan menyuntikkan env vars yang benar.
func SetUpPostgreSQLConnection() *gorm.DB {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("warning: gagal baca .env, lanjut pakai env yang sudah ada:", err)
		}
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("gagal konek ke postgres: %v", err)
	}

	log.Println("koneksi postgres berhasil")
	return db
}

func ClosePostgreSQLConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("gagal ambil raw db saat close: %v", err)
		return
	}
	sqlDB.Close()
	log.Println("koneksi postgres ditutup")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
