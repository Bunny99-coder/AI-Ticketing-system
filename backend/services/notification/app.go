package notification

import (
	"ai-ticketing-backend/services/notification/consumer"
	"ai-ticketing-backend/services/notification/service"

	"github.com/joho/godotenv"
)

func NewApp() {
	_ = godotenv.Load()

	svc := service.NewNotificationService()

	// Start consumer
	consumer.StartConsumer(svc)
}
