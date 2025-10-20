package repository

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/internal/pkg/db"

	"github.com/google/uuid"
)

type TicketRepository interface {
	Create(ticket *models.Ticket) error
	FindByID(id uuid.UUID) (*models.Ticket, error)
	ListByUser(userID uuid.UUID) ([]models.Ticket, error) // User's tickets only
	Update(ticket *models.Ticket) error
}

type ticketRepository struct {
	db *db.DB
}

func NewTicketRepository(db *db.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("User").Where("id = ?", id).First(&ticket).Error // Preload user
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) ListByUser(userID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("User").Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error // Updates timestamps auto
}
