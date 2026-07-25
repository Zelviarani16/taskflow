package cmd

import (
	"log"
	"os"

	"github.com/Zelviarani16/taskflow-api/migrations"
	"gorm.io/gorm"
)

// Command membaca flag CLI (go run . --migrate / --seed / --rollback)
// dan menjalankan aksi yang sesuai alih-alih menjalankan server HTTP.
func Command(db *gorm.DB) {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--migrate":
			if err := migrations.Migrate(db); err != nil {
				log.Fatalf("error migrate: %v", err)
			}
			log.Println("migrate selesai")

		case "--seed":
			if err := migrations.Seed(db); err != nil {
				log.Fatalf("error seed: %v", err)
			}
			log.Println("seed selesai")

		case "--rollback":
			if err := migrations.Rollback(db); err != nil {
				log.Fatalf("error rollback: %v", err)
			}
			log.Println("rollback selesai")

		default:
			log.Printf("argumen tidak dikenali: %s", arg)
		}
	}
}
