package main

import (
	"ai-ticketing-backend/internal/pkg/metrics"
	"ai-ticketing-backend/internal/pkg/redis"
	"ai-ticketing-backend/services/ticket"
	"ai-ticketing-backend/services/ticket/consumer"
	"ai-ticketing-backend/services/ticket/handlers"
	"ai-ticketing-backend/services/ticket/invalidator"
	"ai-ticketing-backend/services/ticket/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Ticket Service on :8081...")
	svc := ticket.Setup()
	h := handlers.NewTicketHandlers(svc)

	r := gin.New()
	r.Use(cors.Default()) // Add CORS middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)
	r.Use(metricsMiddleware) // Metrics for requests

	// Customer routes
	customerApi := r.Group("/api/v1/tickets")
	customerApi.Use(middleware.AuthMiddleware())
	{
		customerApi.POST("/", h.Create)
		customerApi.GET("/:id", h.GetByID)
		customerApi.GET("/", h.ListByUser)
		customerApi.PUT("/:id", h.Update)
	}

	// Agent routes
	agentApi := r.Group("/api/v1/agent/tickets")
	agentApi.Use(middleware.AuthMiddleware(), middleware.AgentAuthMiddleware())
	{
		agentApi.GET("/", h.ListAll)
	}

	metrics.RegisterMetrics() // /metrics endpoint

	cache := redis.New()
	go invalidator.StartInvalidator(cache)

	go consumer.StartConsumer()
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Metrics middleware
func metricsMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start).Seconds()
	metrics.RecordRequest(c.Request.Method, c.FullPath(), c.Writer.Status(), duration)
}
