package main

import (
	"log"

	"github.com/xyz-corp/xyz-football-api/config"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	seedAdminEmail    = "admin@xyz.co.id"
	seedAdminPassword = "Admin#1234"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := seedAdmin(db); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	log.Println("seeding complete")
}

func seedAdmin(db *gorm.DB) error {
	var existing entity.User
	err := db.Where("LOWER(email) = LOWER(?)", seedAdminEmail).First(&existing).Error
	if err == nil {
		log.Println("admin user already exists, skipping")
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(seedAdminPassword), 12)
	if err != nil {
		return err
	}

	user := &entity.User{Email: seedAdminEmail, Password: string(hash)}
	if err := db.Create(user).Error; err != nil {
		return err
	}

	log.Printf("admin user created: %s / %s\n", seedAdminEmail, seedAdminPassword)
	return nil
}
