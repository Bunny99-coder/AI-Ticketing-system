package main

import (
	"ai-ticketing-backend/services/user"
	"log"
)

func main() {
	log.Println("Starting User Service on :8080...")
	app := user.NewApp()
	if err := app.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
