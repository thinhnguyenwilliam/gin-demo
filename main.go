// gin-demo/main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v2handler "github.com/thinhnguyenwilliam/gin-demo/internal/api/v2/handler"
	"github.com/thinhnguyenwilliam/gin-demo/utils"
)

func main() {
	r := gin.Default()

	// ✅ Set max upload size (8MB total request size)
	r.MaxMultipartMemory = 8 << 20 // 8MB

	// Serve static uploaded files
	r.Static("/uploads", "./uploads")

	if err := utils.RegisterCustomValidations(); err != nil {
		panic(err)
	}

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
	newsHandler := v2handler.NewNewsHandler()

	news := v2.Group("/news")
	{
		news.GET("", newsHandler.GetNewsList)         // GET /api/v2/news
		news.GET("/:slug", newsHandler.GetNewsBySlug) // GET /api/v2/news/:slug
	}

	categories := v2.Group("/categories")
	{
		categories.POST("/upload-multiple", categoryHandlerV2.UploadMultipleCategories)
		categories.POST("/upload", categoryHandlerV2.UploadCategory)
		categories.POST("", categoryHandlerV2.CreateCategory)
		categories.GET("/:name", categoryHandlerV2.GetCategory)
	}

	products := v2.Group("/products")
	{
		products.POST("", productHandlerV2.CreateProduct)
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
