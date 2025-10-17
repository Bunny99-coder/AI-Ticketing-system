package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system (e.g., customer or agent)
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // Auto-gen UUID
	Email     string    `json:"email" gorm:"unique;not null"`                              // Unique email
	Password  string    `json:"-" gorm:"not null"`                                         // Hide from JSON output
	Role      string    `json:"role" gorm:"default:customer"`                              // e.g., "customer" or "agent"
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}

// LoginRequest for incoming login data (no ID/times)
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest for incoming registration (validated)
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=customer agent"` // Optional, defaults to customer
}

// MustParseUUID parses a UUID string, panics on error (for simplicity)
func MustParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		panic("invalid UUID: " + s)
	}
	return u
}
