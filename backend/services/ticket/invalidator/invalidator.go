package invalidator

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/internal/pkg/redis"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func StartInvalidator(cache *redis.Client) {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    "ticket-events",
		GroupID:  "cache-invalidator-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	log.Println("Cache Invalidator listening on ticket-events...")
	ctx := context.Background()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		// Handle any event (created/updated) — invalidate user tickets
		var event models.TicketCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			var updatedEvent models.TicketUpdatedEvent
			if json.Unmarshal(msg.Value, &updatedEvent) == nil {
				log.Printf("Invalidated caches for updated ticket %s (user %s)", updatedEvent.TicketID, updatedEvent.UserID)
				cache.CacheDel(ctx, "ticket:"+updatedEvent.TicketID.String())
				cache.CacheDel(ctx, "user_tickets:"+updatedEvent.UserID.String())
			}
			continue
		}
		log.Printf("Invalidated user cache for created ticket %s (user %s)", event.TicketID, event.UserID)
		cache.CacheDel(ctx, "user_tickets:"+event.UserID.String())
	}
}
