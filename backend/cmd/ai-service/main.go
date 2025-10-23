package main

import (
	"ai-ticketing-backend/services/ai"
	"ai-ticketing-backend/services/ai/consumer"
	"log"
)

func main() {
	log.Println("Starting AI Service...")
	svc := ai.Setup()
	consumer.StartConsumer(svc)
}
