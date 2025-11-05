package consumer

import (
	"ai-ticketing-backend/internal/models"
	ai "ai-ticketing-backend/services/ai"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func StartConsumer(aiSvc ai.AIService) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(kafkaBroker, ","),
		Topic:    "ticket-events",
		GroupID:  "ai-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	log.Println("AI Consumer listening on ticket-events...")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			log.Println("Caught signal", sig)
			return
		default:
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			var createdEvent models.TicketCreatedEvent
			if err := json.Unmarshal(msg.Value, &createdEvent); err == nil && createdEvent.TicketID != uuid.Nil {
				if err := aiSvc.ProcessTicketEvent(&createdEvent); err != nil {
					log.Printf("Failed to process created event %s: %v", createdEvent.TicketID, err)
				} else {
					log.Printf("AI processed created ticket %s successfully", createdEvent.TicketID)
				}
				continue
			}

			var contentUpdatedEvent models.TicketContentUpdatedEvent
			if err := json.Unmarshal(msg.Value, &contentUpdatedEvent); err == nil && contentUpdatedEvent.TicketID != uuid.Nil {
				if err := aiSvc.ProcessTicketContent(contentUpdatedEvent.TicketID, contentUpdatedEvent.Title, contentUpdatedEvent.Description); err != nil {
					log.Printf("Failed to process content updated event %s: %v", contentUpdatedEvent.TicketID, err)
				} else {
					log.Printf("AI processed content updated ticket %s successfully", contentUpdatedEvent.TicketID)
				}
				continue
			}

			log.Printf("Failed to unmarshal event: %s", string(msg.Value))
		}
	}
}
