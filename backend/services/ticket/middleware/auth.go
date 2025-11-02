package middleware

import (
	"ai-ticketing-backend/internal/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT and sets user_id/role in context (no DB lookup)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT secret not configured"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id in token claims"})
				c.Abort()
				return
			}
			userID := models.MustParseUUID(userIDStr)
			role, _ := claims["role"].(string) // Optional fallback to ""
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Next()
		}
	}
}

// New: AgentAuthMiddleware checks for 'agent' role
func AgentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "agent" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Agent role required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
