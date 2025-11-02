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
			var event models.TicketCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			if err := aiSvc.ProcessTicketEvent(&event); err != nil {
				log.Printf("Failed to process event %s: %v", event.TicketID, err)
			} else {
				log.Printf("AI processed ticket %s successfully", event.TicketID)
			}
		}
	}
}
