// gin-demo/internal/api/v2/handler/category_handler.go
package handler

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/gin-demo/utils"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (h *CategoryHandler) UploadMultipleCategories(c *gin.Context) {
	var req uploadCategoryRequest

	// Bind form fields
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	// Get multiple files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read multipart form",
		})
		return
	}

	files := form.File["images"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "at least one image is required",
		})
		return
	}

	var savedPaths []string

	for _, file := range files {
		// Validate
		if err := validateImageFile(file); err != nil {
			// 🔥 rollback
			deleteFiles(savedPaths)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid file: " + file.Filename + " - " + err.Error(),
			})
			return
		}
		// Save
		savePath, err := saveImage(c, file)
		if err != nil {
			// 🔥 rollback
			deleteFiles(savedPaths)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save file: " + file.Filename,
			})
			return
		}
		baseURL := "http://192.168.1.8:8085"
		imageURL := baseURL + "/" + savePath
		savedPaths = append(savedPaths, imageURL)

	}

	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Multiple images uploaded successfully",
		"data": gin.H{
			"name":       req.Name,
			"is_active":  isActive,
			"images":     savedPaths,
			"total_file": len(savedPaths),
		},
	})
}

func deleteFiles(paths []string) {
	for _, path := range paths {
		_ = os.Remove(path) // ignore error
	}
}

type uploadCategoryRequest struct {
	Name     string `form:"name" binding:"required,min=3,max=50"`
	IsActive *bool  `form:"is_active" binding:"required"`
}

var allowedImageExt = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)$`)

func validateImageFile(file *multipart.FileHeader) error {

	// 1️⃣ Check size (2MB)
	if file.Size > 2<<20 {
		return errors.New("file size must not exceed 2MB")
	}

	// 2️⃣ Check extension
	if !allowedImageExt.MatchString(file.Filename) {
		return errors.New("only jpg, jpeg, png, webp allowed")
	}

	// 3️⃣ Check MIME type (real security)
	openedFile, err := file.Open()
	if err != nil {
		return errors.New("failed to open file")
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	_, err = openedFile.Read(buffer)
	if err != nil {
		return errors.New("failed to read file")
	}

	fileType := http.DetectContentType(buffer)

	switch fileType {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return errors.New("invalid image type")
	}
}

func saveImage(c *gin.Context, file *multipart.FileHeader) (string, error) {

	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	savePath := filepath.Join(uploadDir, newFileName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", err
	}

	return savePath, nil
}

func (h *CategoryHandler) UploadCategory(c *gin.Context) {
	var req uploadCategoryRequest

	// Bind form fields
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	// Get file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image file is required",
		})
		return
	}

	// Validate file
	if err := validateImageFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save file
	savePath, err := saveImage(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save file",
		})
		return
	}

	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category uploaded successfully",
		"data": gin.H{
			"name":      req.Name,
			"is_active": isActive,
			"image":     savePath,
		},
	})
}

type createCategoryRequest struct {
	Name        string `form:"name" binding:"required,min=3,max=50"`
	Description string `form:"description" binding:"omitempty,max=255"`
	IsActive    *bool  `form:"is_active" binding:"required"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req createCategoryRequest

	// This automatically detects form-urlencoded
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data": gin.H{
			"name":        req.Name,
			"description": req.Description,
			"is_active":   isActive,
		},
	})
}

type getCategoryRequest struct {
	Name string `uri:"name" binding:"required,oneof=php python golang"`
}

// GET /api/v2/categories/:name
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	var req getCategoryRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": req.Name,
		"message":  "Category is valid",
	})
}
