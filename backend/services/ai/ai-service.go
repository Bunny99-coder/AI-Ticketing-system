package ai

import (
	"ai-ticketing-backend/internal/models"
	"ai-ticketing-backend/services/ai/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings" // New: For trim

	"github.com/google/uuid"
)

type AIService interface {
	ProcessTicketEvent(event *models.TicketCreatedEvent) error
}

type aiService struct {
	repo   repository.TicketRepository
	apiKey string
}

func NewAIService(repo repository.TicketRepository) AIService {
	key := os.Getenv("GEMINI_API_KEY")
	log.Printf("GEMINI_API_KEY loaded: %s", key)
	if key == "" {
		panic("GEMINI_API_KEY not set")
	}
	return &aiService{repo: repo, apiKey: key}
}

func (s *aiService) ProcessTicketEvent(event *models.TicketCreatedEvent) error {
	log.Printf("ProcessTicketEvent called for ticket ID: %s", event.TicketID)
	// Shorter prompt
	prompt := fmt.Sprintf(`Classify ticket: %s. Description: %s. JSON only: {"category": "Billing|Bug|Feature|Support", "priority": "low|medium|high", "suggestion": "1-2 sentence reply"}`, event.Title, event.Description)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.1,
			"maxOutputTokens": 1000,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + s.apiKey
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	log.Printf("Making Gemini API call for ticket ID: %s", event.TicketID)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Gemini API call failed: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Received Gemini API response for ticket ID: %s, Status: %d", event.TicketID, resp.StatusCode)
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Gemini error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	// Safe parsing
	candidatesIface, ok := response["candidates"].(interface{})
	if !ok || candidatesIface == nil {
		log.Printf("No candidates in Gemini response: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}
	candidates, ok := candidatesIface.([]interface{})
	if !ok || len(candidates) == 0 {
		log.Printf("Empty candidates in Gemini response: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}
	candidate := candidates[0].(map[string]interface{})
	contentIface, ok := candidate["content"].(map[string]interface{})
	if !ok {
		log.Printf("No content in candidate: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}
	partsIface, ok := contentIface["parts"].(interface{})
	if !ok || partsIface == nil {
		log.Printf("No parts in content: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}
	parts, ok := partsIface.([]interface{})
	if !ok || len(parts) == 0 {
		log.Printf("Empty parts in content: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}
	part := parts[0].(map[string]interface{})
	text, ok := part["text"].(string)
	if !ok {
		log.Printf("No text in part: %s", string(body))
		return s.updateTicketWithDefaults(event.TicketID)
	}

	// Fixed: Strip Markdown backticks and "json" label
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
		text = strings.TrimSuffix(text, "```")
	}
	text = strings.TrimSpace(text)
	log.Printf("Cleaned Gemini text: %s", text)

	// Parse JSON from cleaned text
	var aiResponse struct {
		Category   string `json:"category"`
		Priority   string `json:"priority"`
		Suggestion string `json:"suggestion"`
	}
	if err := json.Unmarshal([]byte(text), &aiResponse); err != nil {
		log.Printf("Failed to parse AI JSON: %v (raw cleaned: %s)", err, text)
		return s.updateTicketWithDefaults(event.TicketID)
	}

	// Fetch and update ticket
	ticket, err := s.repo.GetByID(event.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}
	ticket.Category = aiResponse.Category
	ticket.Priority = aiResponse.Priority
	ticket.Suggestion = aiResponse.Suggestion
	ticket.Status = "classified"
	if err := s.repo.Update(ticket); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	log.Printf("AI processed ticket %s: Category=%s, Priority=%s, Suggestion=%s", event.TicketID, aiResponse.Category, aiResponse.Priority, aiResponse.Suggestion)
	return nil
}

// Fallback update with defaults if Gemini fails
func (s *aiService) updateTicketWithDefaults(ticketID uuid.UUID) error {
	ticket, err := s.repo.GetByID(ticketID)
	if err != nil {
		return err
	}
	ticket.Category = "Unknown"
	ticket.Priority = "low"
	ticket.Suggestion = "Please provide more details for assistance."
	ticket.Status = "classified"
	if err := s.repo.Update(ticket); err != nil {
		return err
	}
	log.Printf("AI fallback updated ticket %s with defaults", ticketID)
	return nil
}
