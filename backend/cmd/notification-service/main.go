package main

import (
	"ai-ticketing-backend/services/notification"
	"log"
)

func main() {
	log.Println("Starting Notification Service...")
	notification.NewApp()
}
