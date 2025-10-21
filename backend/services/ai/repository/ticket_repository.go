package repository

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/internal/pkg/db"

	"github.com/google/uuid"
)

type TicketRepository interface {
	GetByID(id uuid.UUID) (*models.Ticket, error) // Fetch for update
	Update(ticket *models.Ticket) error
}

type ticketRepository struct {
	db *db.DB
}

func NewTicketRepository(db *db.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("User").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}
