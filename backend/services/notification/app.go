package notification

import (
	"github.com/joho/godotenv"
)

func Setup() NotificationService {
	_ = godotenv.Load()

	svc := NewNotificationService()

	return svc
}