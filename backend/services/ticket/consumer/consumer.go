package consumer

import (
	"ai-ticketing-backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumer() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		// Default to Kafka container hostname when running in Docker
		kafkaBroker = "kafka:9092"
	}
	log.Printf("KAFKA_BROKER: %s", kafkaBroker)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(kafkaBroker, ","),
		GroupID:  "ticket-consumer-group",
		Topic:    "ticket-events",
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
