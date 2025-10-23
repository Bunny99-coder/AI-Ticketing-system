package main

import (
	"ai-ticketing-backend/internal/pkg/redis"
	"ai-ticketing-backend/services/ticket"
	"ai-ticketing-backend/services/ticket/consumer"
	"ai-ticketing-backend/services/ticket/handlers"
	"ai-ticketing-backend/services/ticket/invalidator"
	"ai-ticketing-backend/services/ticket/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Ticket Service on :8081...")
	svc := ticket.Setup()
	h := handlers.NewTicketHandlers(svc)

	r := gin.Default()
	api := r.Group("/api/v1/tickets")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/", h.Create)
		api.GET("/:id", h.GetByID)
		api.GET("/", h.ListByUser)
		api.PUT("/:id", h.Update)
	}

	cache := redis.New()
	go invalidator.StartInvalidator(cache)

	go consumer.StartConsumer()
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
