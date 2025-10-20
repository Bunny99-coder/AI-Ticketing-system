package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Create a new Kafka reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "test-consumer",
		Topic:   "ticket-events",
		// Set some limits for reading
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer reader.Close()

	fmt.Println("✅ Kafka consumer initialized")

	ctx := context.Background()

	// Read a single message to test
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Fatalf("failed to read message: %v", err)
	}

	fmt.Printf("Received message at offset %d: %s\n", msg.Offset, string(msg.Value))

	// Optional: continuously read messages
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %v", err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("Received message at offset %d: %s\n", msg.Offset, string(msg.Value))
	}
}
