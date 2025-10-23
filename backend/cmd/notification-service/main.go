package main

import (
	"ai-ticketing-backend/services/notification"
	"ai-ticketing-backend/services/notification/consumer"
	"log"
)

func main() {
	log.Println("Starting Notification Service...")
	svc := notification.Setup()
	consumer.StartConsumer(svc)
}
