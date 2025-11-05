package handlers

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/ticket"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandlers struct {
	svc ticket.TicketService
}

func NewTicketHandlers(svc ticket.TicketService) *TicketHandlers {
	return &TicketHandlers{svc: svc}
}

// CreateHandler for POST /api/v1/tickets
func (h *TicketHandlers) Create(c *gin.Context) {
	userIDStr, exists := c.Get("user_id") // From middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, _ := userIDStr.(uuid.UUID)

	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.svc.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

// GetByIDHandler for GET /api/v1/tickets/:id
func (h *TicketHandlers) GetByID(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := userIDStr.(uuid.UUID)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var ticket *models.Ticket
	// If the user is an agent, they can get any ticket
	if role == "agent" {
		ticket, err = h.svc.GetByID(id, uuid.Nil) // uuid.Nil bypasses ownership check
	} else {
		ticket, err = h.svc.GetByID(id, userID)
	}

	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// ListByUserHandler for GET /api/v1/tickets
func (h *TicketHandlers) ListByUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	userID, _ := userIDStr.(uuid.UUID)

	tickets, err := h.svc.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// New: ListAll handler for agents
func (h *TicketHandlers) ListAll(c *gin.Context) {
	tickets, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// UpdateHandler for PUT /api/v1/tickets/:id
func (h *TicketHandlers) Update(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := userIDStr.(uuid.UUID)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req models.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent customers from updating status
	if role != "agent" && req.Status != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "only agents can update ticket status"})
		return
	}

	ticket, err := h.svc.Update(id, &req, userID, role.(string))
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}
