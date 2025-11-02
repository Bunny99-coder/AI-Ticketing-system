package handlers

import (
	"net/http"

	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlers struct {
	svc user.UserService
}

func NewUserHandlers(svc user.UserService) *UserHandlers {
	return &UserHandlers{svc: svc}
}

// RegisterHandler for POST /api/v1/users/register
func (h *UserHandlers) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	user.Password = "" // Hide
	c.JSON(http.StatusCreated, user)
}

// LoginHandler for POST /api/v1/users/login
func (h *UserHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.svc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetUserHandler for GET /api/v1/users/:id (needs JWT middleware later)
func (h *UserHandlers) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	user, err := h.svc.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUsersHandler for GET /api/v1/users
func (h *UserHandlers) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
