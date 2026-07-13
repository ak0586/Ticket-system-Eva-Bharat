package models

import "gorm.io/gorm"

// User represents the users table in the database.
// Python equivalent: An SQLAlchemy Base class combined with a Pydantic model.
type User struct {
	// gorm.Model automatically injects ID, CreatedAt, UpdatedAt, and DeletedAt columns.
	gorm.Model
	// The `gorm:"uniqueIndex"` tag enforces uniqueness in the database.
	// The `json:"email"` tag controls how this field is named when returned as JSON.
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	// `json:"-"` guarantees the password hash is NEVER included in a JSON response.
	Password string `json:"-" gorm:"not null"`
}

// Ticket represents the tickets table.
type Ticket struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	// Default status is 'open'.
	Status      string `json:"status" gorm:"default:open"`
	// UserID acts as a foreign key linking back to the User who created it.
	// The `gorm:"index"` tag speeds up queries when we filter by user_id.
	UserID      uint   `json:"user_id" gorm:"index"`
}
