package notification // Same as consumer.go

import (
	"ai-ticketing-backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type NotificationService interface {
	SendUpdatedNotification(event *models.TicketUpdatedEvent) error
}

type notificationService struct {
	slackWebhook string
}

func NewNotificationService() NotificationService { // Exported (capital N)—visible in same package
	webhook := os.Getenv("SLACK_WEBHOOK_URL")
	return &notificationService{slackWebhook: webhook}
}

func (s *notificationService) SendUpdatedNotification(event *models.TicketUpdatedEvent) error {
	message := fmt.Sprintf("Ticket %s updated: Status changed from %s to %s. User: %s", event.TicketID, event.OldStatus, event.NewStatus, event.UserID)

	// Real email via Gmail SMTP
	from := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	to := os.Getenv("EMAIL_RECEIVER")
	if from == "" || password == "" || to == "" {
		log.Printf("Email env vars not set—skipping real email")
	} else {
		auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
		msg := []byte("To: " + to + "\r\n" +
			"Subject: Ticket Update: " + event.TicketID.String() + "\r\n" +
			"\r\n" + message + "\r\n")

		err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
		if err != nil {
			log.Printf("Email send failed: %v", err)
			return err
		}
		log.Printf("Real email sent to %s: %s", to, message)
	}

	// Optional Slack
	if s.slackWebhook != "" {
		payload := map[string]string{"text": message}
		payloadBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", s.slackWebhook, bytes.NewBuffer(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Slack send failed: %v", err)
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("Slack error: %d", resp.StatusCode)
		} else {
			log.Printf("Slack notification sent for ticket %s", event.TicketID)
		}
	}

	log.Printf("Notification processed for ticket %s", event.TicketID)
	return nil
}
