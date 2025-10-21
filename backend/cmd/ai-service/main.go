package main

import (
	"ai-ticketing-backend/services/ai"
	"log"
)

func main() {
	log.Println("Starting AI Service...")
	ai.NewApp() // Starts the Kafka consumer loop (blocks)
}
