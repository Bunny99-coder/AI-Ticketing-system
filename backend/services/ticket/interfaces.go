package ticket

import (
	"ai-ticketing-backend/internal/models"
	"github.com/google/uuid"
)

type TicketService interface {
	Create(req *models.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
	GetByID(id uuid.UUID, userID uuid.UUID) (*models.Ticket, error)
	ListByUser(userID uuid.UUID) ([]models.Ticket, error)
	Update(id uuid.UUID, req *models.UpdateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
}
