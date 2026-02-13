// gin-demo/internal/api/v2/handler/product_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/gin-demo/utils"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

type ProductImageRequest struct {
	Link string `json:"link" binding:"required,url,image_ext"`
}

type createProductRequest struct {
	Name        string                `json:"name" binding:"required,min=3,max=100"`
	Slug        string                `json:"slug" binding:"required,slug"`
	Price       float64               `json:"price" binding:"required,gt=0"`
	Category    string                `json:"category" binding:"required,oneof=php python golang"`
	Description string                `json:"description" binding:"omitempty,max=500"`
	IsActive    *bool                 `json:"is_active" binding:"required"`
	Images      []ProductImageRequest `json:"images" binding:"required,min=1,dive"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req createProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// simulate saving to DB
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data": gin.H{
			"name":        req.Name,
			"slug":        req.Slug,
			"price":       req.Price,
			"category":    req.Category,
			"description": req.Description,
			"is_active":   isActive,
			"images":      req.Images,
		},
	})

}

type searchProductRequest struct {
	Search string `form:"search" binding:"omitempty,min=2,max=100,search"`
	Limit  *int   `form:"limit" binding:"omitempty,gt=0,lte=100"`
}

// GET /api/v2/products?search=&limit=
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	var req searchProductRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	// default limit
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}

	c.JSON(http.StatusOK, gin.H{
		"search": req.Search,
		"limit":  limit,
	})
}

type getProductRequest struct {
	Slug string `uri:"slug" binding:"required,slug,min=5,max=100"`
}

// GET /api/v2/products/:slug
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	var req getProductRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"slug": req.Slug,
	})
}
