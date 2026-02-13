// gin-demo/internal/middleware/logger.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// before handler
		log.Println("Incoming request:", c.Request.Method, c.Request.URL.Path)

		// continue to next
		c.Next()

		// after handler
		duration := time.Since(start)
		log.Println("Request completed in:", duration)
	}
}
