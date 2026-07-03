package db

import (
	"log"
	"os"

	"ticket-system/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "tickets.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("db: failed to open %s: %v", path, err)
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Ticket{}); err != nil {
		log.Fatalf("db: migration failed: %v", err)
	}
}
