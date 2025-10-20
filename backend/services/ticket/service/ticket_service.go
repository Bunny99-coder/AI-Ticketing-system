package service

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/ticket/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type TicketService interface {
	Create(req *models.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
	GetByID(id uuid.UUID, userID uuid.UUID) (*models.Ticket, error)
	ListByUser(userID uuid.UUID) ([]models.Ticket, error)
	Update(id uuid.UUID, req *models.UpdateTicketRequest, userID uuid.UUID) (*models.Ticket, error)
}

type ticketService struct {
	repo   repository.TicketRepository
	writer *kafka.Writer // Kafka writer for events
}

// NewTicketService initializes the service with Kafka writer
func NewTicketService(repo repository.TicketRepository) TicketService {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "ticket-events",
		Balancer: &kafka.LeastBytes{},
	}

	return &ticketService{
		repo:   repo,
		writer: writer,
	}
}

func (s *ticketService) Create(req *models.CreateTicketRequest, userID uuid.UUID) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}
	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}

	// Publish event to Kafka
	event := models.TicketCreatedEvent{
		TicketID:    ticket.ID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return ticket, nil // Don't fail create on event error
	}

	err = s.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(ticket.ID.String()),
			Value: eventBytes,
		},
	)
	if err != nil {
		log.Printf("failed to produce event: %v", err)
	} else {
		log.Println("Published ticket_created event for ID:", ticket.ID)
	}

	return ticket, nil
}

func (s *ticketService) GetByID(id uuid.UUID, userID uuid.UUID) (*models.Ticket, error) {
	ticket, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ticket.UserID != userID {
		return nil, fmt.Errorf("unauthorized: not your ticket")
	}
	return ticket, nil
}

func (s *ticketService) ListByUser(userID uuid.UUID) ([]models.Ticket, error) {
	return s.repo.ListByUser(userID)
}

func (s *ticketService) Update(id uuid.UUID, req *models.UpdateTicketRequest, userID uuid.UUID) (*models.Ticket, error) {
	ticket, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ticket.UserID != userID {
		return nil, fmt.Errorf("unauthorized: not your ticket")
	}
	if req.Title != "" {
		ticket.Title = req.Title
	}
	if req.Description != "" {
		ticket.Description = req.Description
	}
	if req.Status != "" {
		ticket.Status = req.Status
	}
	if err := s.repo.Update(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}
