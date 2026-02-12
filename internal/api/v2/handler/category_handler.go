package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// Allowed categories
var allowedCategories = map[string]bool{
	"php":    true,
	"python": true,
	"golang": true,
}

// GET /api/v2/categories/:name
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	name := c.Param("name")

	if !allowedCategories[name] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category (allowed: php, python, golang)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": name,
		"message":  "Category is valid",
	})
}
