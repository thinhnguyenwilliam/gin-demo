package handler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct{}

func NewNewsHandler() *NewsHandler {
	return &NewsHandler{}
}

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ===============================
// GET /api/v2/news
// ===============================
func (h *NewsHandler) GetNewsList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []string{
			"golang-1-22-released",
			"gin-framework-update",
		},
	})
}

// ===============================
// GET /api/v2/news/:slug
// ===============================
func (h *NewsHandler) GetNewsBySlug(c *gin.Context) {
	slug := c.Param("slug")

	if !slugRegex.MatchString(slug) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid news slug",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"slug": slug,
		"data": "News detail content",
	})
}
