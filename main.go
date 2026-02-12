// gin-demo/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v2handler "github.com/thinhnguyenwilliam/gin-demo/internal/api/v2/handler"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Thinh 🚀 Gin is working!",
		})
	})

	// =========================
	// API V2 Group
	// =========================
	v2 := r.Group("/api/v2")

	categoryHandlerV2 := v2handler.NewCategoryHandler()
	productHandlerV2 := v2handler.NewProductHandler()
	userHandlerV2 := v2handler.NewUserHandler()

	categories := v2.Group("/categories")
	{
		categories.GET("/:name", categoryHandlerV2.GetCategory)
	}

	products := v2.Group("/products")
	{
		products.GET("", productHandlerV2.SearchProducts) // query search
		products.GET("/:slug", productHandlerV2.GetProductBySlug)
	}

	user := v2.Group("/users")
	{
		user.GET("/:id", userHandlerV2.GetUserByID)
		user.GET("", userHandlerV2.GetUsers)
	}

	r.Run(":8085")
}
