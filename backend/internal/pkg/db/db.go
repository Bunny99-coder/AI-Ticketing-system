package db

import (
	"ai-ticketing-backend/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Create UUID extension (fixed: proper Exec handling)
	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		return nil, fmt.Errorf("failed to create UUID extension: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	// In New(), after AutoMigrate(&models.User{})
	if err := db.AutoMigrate(&models.Ticket{}); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return &DB{db}, nil
}
