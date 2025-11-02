package user

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/user/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      repository.UserRepository
	jwtSecret string // Env var in prod
}

func NewUserService(repo repository.UserRepository) UserService {
	secret := os.Getenv("JWT_SECRET") // Load from env (fallback if unset)
	if secret == "" {
		secret = "super-secret-key-change-me" // Fallback for local
	}
	return &userService{repo: repo, jwtSecret: secret}
}

func (s *userService) GetJWTSecret() string {
	return s.jwtSecret
}

func (s *userService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if email exists
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req *models.LoginRequest) (string, *models.User, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24h
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	user.Password = "" // Clear password before returning
	return tokenString, user, nil
}

func (s *userService) GetUser(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// Hide password in response
	user.Password = ""
	return user, nil
}

func (s *userService) ListUsers() ([]models.User, error) {
	users, err := s.repo.List()
	for i := range users {
		users[i].Password = "" // Hide passwords
	}
	return users, err
}
