// gin-demo/internal/api/v2/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/gin-demo/utils"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// Struct for validation
type getUserRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GET /api/v2/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	var req getUserRequest

	// Bind URI params + validate automatically
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      req.ID,
		"message": "User found (v2)",
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get users v2",
	})
}
