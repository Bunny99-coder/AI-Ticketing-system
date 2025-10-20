package ticket

import (
	"ai-ticketing-backend/internal/pkg/db"
	"ai-ticketing-backend/services/ticket/handlers"
	"ai-ticketing-backend/services/ticket/middleware"
	"ai-ticketing-backend/services/ticket/repository"
	"ai-ticketing-backend/services/ticket/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	router *gin.Engine
}

func NewApp() *App {
	_ = godotenv.Load()

	// Load from env with fallbacks (using YOUR var names)
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

	// Explicit DSN (use "dbname=", not "database=")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, username, pass, dbname, port)

	dbConn, err := db.New(dsn)
	if err != nil {
		panic(err)
	}

	repo := repository.NewTicketRepository(dbConn)
	svc := service.NewTicketService(repo)
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

	return &App{router: r}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
