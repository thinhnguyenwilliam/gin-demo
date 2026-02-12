// gin-demo/internal/api/v2/handler/product_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/gin-demo/internal/pkg/validator"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// GET /api/v2/products?search=&limit=
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	search := c.Query("search")

	limit, err := validator.ParseLimit(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"search": search,
		"limit":  limit,
	})
}

// GET /api/v2/products/:slug
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if !validator.IsValidSlug(slug) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid slug",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"slug": slug})
}
