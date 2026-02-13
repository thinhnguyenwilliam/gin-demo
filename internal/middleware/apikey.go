package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/gin-demo/internal/config"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientKey := c.GetHeader("x-api-key")

		if clientKey == "" || clientKey != config.AppConfig.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
