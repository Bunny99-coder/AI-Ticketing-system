package models

import "github.com/google/uuid"

// TicketCreatedEvent for Kafka
type TicketCreatedEvent struct {
	TicketID    uuid.UUID `json:"ticket_id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"` // ISO string
}

type TicketUpdatedEvent struct {
	TicketID  uuid.UUID `json:"ticket_id"`
	UserID    uuid.UUID `json:"user_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	UpdatedAt string    `json:"updated_at"`
}
