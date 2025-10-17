package user

import (
	"ai-ticketing-backend/internal/pkg/db"
	"ai-ticketing-backend/services/user/handlers"
	"ai-ticketing-backend/services/user/middleware"
	"ai-ticketing-backend/services/user/repository"
	"ai-ticketing-backend/services/user/service"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	router *gin.Engine
}

func NewApp() *App {
	// Load .env file (if exists)
	_ = godotenv.Load() // Ignores errors if no .env

	// Load from env with fallbacks (using YOUR var names)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1" // Matches your .env
	}
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "ticket_user"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "ticket123" // Matches your .env!
	}
	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "Ticket"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, username, pass, dbname, port)

	dbConn, err := db.New(dsn)
	if err != nil {
		panic(err)
	}

	repo := repository.NewUserRepository(dbConn)
	svc := service.NewUserService(repo)
	h := handlers.NewUserHandlers(svc)

	r := gin.Default()
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

	return &App{router: r}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
