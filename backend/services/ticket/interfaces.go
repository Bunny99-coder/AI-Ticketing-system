package ticket

import (
	"ai-ticketing-backend/internal/models"
	"github.com/google/uuid"
)

type TicketService interface {
	Create(req *models.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
	GetByID(id uuid.UUID, userID uuid.UUID) (*models.Ticket, error)
	ListByUser(userID uuid.UUID) ([]models.Ticket, error)
	ListAll() ([]models.Ticket, error) // New: For agents
	Update(id uuid.UUID, req *models.UpdateTicketRequest, userID uuid.UUID, role string) (*models.Ticket, error)
	CustomerUpdate(id uuid.UUID, req *models.CustomerUpdateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
}
