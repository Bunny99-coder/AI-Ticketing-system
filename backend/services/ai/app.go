package ai

import (
	"ai-ticketing-backend/internal/pkg/db"
	"ai-ticketing-backend/services/ai/consumer"
	"ai-ticketing-backend/services/ai/repository"
	"ai-ticketing-backend/services/ai/service"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func NewApp() {
	_ = godotenv.Load()

	// DB setup (shared)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "ticket_user"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "ticket123"
	}
	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "Ticket"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, username, pass, dbname, port)

	dbConn, err := db.New(dsn)
	if err != nil {
		panic(err)
	}

	repo := repository.NewTicketRepository(dbConn)
	svc := service.NewAIService(repo)

	// Start consumer
	consumer.StartConsumer(svc)
}
