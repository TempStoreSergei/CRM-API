package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Get allowed origins from environment or use defaults
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			// In development, allow localhost origins
			env := os.Getenv("ENVIRONMENT")
			if env == "development" || env == "" {
				allowedOrigins = "http://localhost:3000,http://localhost:8080"
			}
		}

		// Check if origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin, allowedOrigins string) bool {
	if allowedOrigins == "*" {
		return true
	}
	if origin == "" {
		return false
	}
	
	// Parse comma-separated origins
	for _, allowed := range splitOrigins(allowedOrigins) {
		if allowed == origin {
			return true
		}
	}
	return false
}

func splitOrigins(origins string) []string {
	result := []string{}
	current := ""
	for _, c := range origins {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
