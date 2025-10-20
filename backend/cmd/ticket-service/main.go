package main

import (
	"ai-ticketing-backend/services/ticket"
	"ai-ticketing-backend/services/ticket/consumer"
	"log"
)

func main() {
	log.Println("Starting Ticket Service on :8081...")
	app := ticket.NewApp()
	go consumer.StartConsumer()
	if err := app.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
