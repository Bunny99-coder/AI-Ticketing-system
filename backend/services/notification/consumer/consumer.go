package consumer

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/notification/service"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(notifSvc service.NotificationService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "ticket-events",
		GroupID:  "notification-consumer-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	log.Println("Notification Consumer listening on ticket-events...")
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
			// Filter for updated events

			// In the loop, after ReadMessage
			var event models.TicketUpdatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				var createdEvent models.TicketCreatedEvent
				if json.Unmarshal(msg.Value, &createdEvent) == nil {
					log.Printf("Skipped created event %s (no notification)", createdEvent.TicketID)
				} else {
					log.Printf("Failed to unmarshal event: %v", err)
				}
				continue
			}
			if err := notifSvc.SendUpdatedNotification(&event); err != nil {
				log.Printf("Failed to send notification for %s: %v", event.TicketID, err)
			} else {
				log.Printf("Notification sent for updated ticket %s", event.TicketID)
			}
		}
	}
}
