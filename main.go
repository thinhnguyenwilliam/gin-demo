package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Thinh 🚀 Gin is working!",
		})
	})

	r.Run(":8085") // localhost:8085
}
