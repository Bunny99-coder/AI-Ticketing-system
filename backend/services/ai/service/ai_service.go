package service

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
	if key == "" {
		panic("GEMINI_API_KEY not set")
	}
	return &aiService{repo: repo, apiKey: key}
}

func (s *aiService) ProcessTicketEvent(event *models.TicketCreatedEvent) error {
	// Prompt for classification & suggestion
	prompt := fmt.Sprintf(`Classify this support ticket and suggest an auto-reply.

Title: %s
Description: %s

Respond with JSON only (no extra text):
{"category": "Billing|Bug|Feature|Support", "priority": "low|medium|high", "suggestion": "Brief auto-reply for agent (1-2 sentences)"}`, event.Title, event.Description)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.2, // Low for consistent JSON
			"maxOutputTokens": 150,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Fixed: Use "gemini-1.5-flash" model (free, v1beta supported)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + s.apiKey
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Gemini API call failed: %w", err)
	}
	defer resp.Body.Close()

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

	// Fixed: Safe parsing with nil checks
	candidatesIface, ok := response["candidates"].(interface{})
	if !ok || candidatesIface == nil {
		log.Printf("No candidates in Gemini response: %s", string(body))
		return fmt.Errorf("no candidates in Gemini response")
	}
	candidates, ok := candidatesIface.([]interface{})
	if !ok || len(candidates) == 0 {
		log.Printf("Empty candidates in Gemini response: %s", string(body))
		return fmt.Errorf("empty candidates in Gemini response")
	}
	candidate := candidates[0].(map[string]interface{})
	contentIface, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no content in candidate")
	}
	partsIface, ok := contentIface["parts"].(interface{})
	if !ok || partsIface == nil {
		return fmt.Errorf("no parts in content")
	}
	parts, ok := partsIface.([]interface{})
	if !ok || len(parts) == 0 {
		return fmt.Errorf("empty parts in content")
	}
	part := parts[0].(map[string]interface{})
	text, ok := part["text"].(string)
	if !ok {
		return fmt.Errorf("no text in part")
	}

	// Parse JSON from text
	var aiResponse struct {
		Category   string `json:"category"`
		Priority   string `json:"priority"`
		Suggestion string `json:"suggestion"`
	}
	if err := json.Unmarshal([]byte(text), &aiResponse); err != nil {
		log.Printf("Failed to parse AI JSON: %v (raw: %s)", err, text)
		return err
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
