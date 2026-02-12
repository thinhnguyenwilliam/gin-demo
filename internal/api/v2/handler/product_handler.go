// gin-demo/internal/api/v2/handler/product_handler.go
package handler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// ✅ package-level variable
var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// GET /api/v2/products/:slug
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if !slugRegex.MatchString(slug) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product slug format",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"slug":    slug,
		"message": "Product found (v2)",
	})
}
