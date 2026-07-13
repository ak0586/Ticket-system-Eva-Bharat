package db

import (
	"log"
	"os"

	"ticket-system/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a global variable holding the database connection pool.
// Python equivalent: A global sessionmaker or engine instance in SQLAlchemy.
var DB *gorm.DB

func Init() {
	// Read the database path from the environment.
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "tickets.db" // Default fallback
	}

	var err error
	// Open the SQLite database connection using GORM.
	// Python equivalent: engine = create_engine("sqlite:///tickets.db")
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Keep the console clean
	})
	if err != nil {
		// log.Fatalf prints the error and crashes the program (exit code 1).
		log.Fatalf("db: failed to open %s: %v", path, err)
	}

	// AutoMigrate examines the User and Ticket structs and automatically creates
	// or updates the database tables to match the struct definitions.
	// Python equivalent: Base.metadata.create_all(engine)
	if err := DB.AutoMigrate(&models.User{}, &models.Ticket{}); err != nil {
		log.Fatalf("db: migration failed: %v", err)
	}
}
