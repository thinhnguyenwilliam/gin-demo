// gin-demo/internal/api/v2/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GET /api/v2/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id (must be UUID)",
		})
		return
	}

	// Simulate response
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "User found (v2)",
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users v2",
	})
}
