package consumer

import (
	"ai-ticketing-backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer() {
	// Kafka reader configuration
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "ticket-consumer-group",
		Topic:   "ticket-events",
		// Set a reasonable timeout for reading messages
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer reader.Close()

	fmt.Println("Consuming from topic:", "ticket-events")

	ctx := context.Background()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			time.Sleep(time.Second) // backoff
			continue
		}

		var event models.TicketCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v\n", err)
			continue
		}

		fmt.Printf(
			"Consumed ticket_created: ID=%s, Title=%s, User=%s\n",
			event.TicketID,
			event.Title,
			event.UserID,
		)
	}
}
