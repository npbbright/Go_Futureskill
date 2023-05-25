package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var validAPIKeys = map[string]bool{
	"your-api-key": true,
	// Add more valid API keys if necessary
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("API-Key")
		if _, ok := validAPIKeys[apiKey]; !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
