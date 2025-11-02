package main

import (
	"ai-ticketing-backend/services/user"
	"ai-ticketing-backend/services/user/handlers"
	"ai-ticketing-backend/services/user/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting User Service on :8080...")
	svc := user.Setup()
	h := handlers.NewUserHandlers(svc)

	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.Default())
	api := r.Group("/api/v1/users")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}
	protected := r.Group("/api/v1/users")
	protected.Use(middleware.AuthMiddleware(svc))
	{
		protected.GET("/:id", h.GetUser)
		protected.GET("/", h.ListUsers)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
