package models

import (
	"time"

	"github.com/google/uuid"
)

// Ticket represents a support ticket
type Ticket struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:open"`                  // e.g., "open", "in_progress", "closed"
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`           // Foreign key to User
	User        User      `json:"user" gorm:"foreignKey:UserID;references:ID"` // Optional: Load user on query
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}

// CreateTicketRequest for incoming data
type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required,min=5"`
	Description string `json:"description" binding:"required,min=10"`
}

// UpdateTicketRequest for updates
type UpdateTicketRequest struct {
	Title       string `json:"title" binding:"omitempty,min=5"`
	Description string `json:"description" binding:"omitempty,min=10"`
	Status      string `json:"status" binding:"omitempty,oneof=open in_progress closed"`
}
