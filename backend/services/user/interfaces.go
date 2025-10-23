package user

import (
	"ai-ticketing-backend/internal/models"
	"github.com/google/uuid"
)

type UserService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (string, error)
	GetUser(id uuid.UUID) (*models.User, error)
	ListUsers() ([]models.User, error)
	GetJWTSecret() string
}
