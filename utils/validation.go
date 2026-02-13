// gin-demo/utils/validation.go
package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	slugRegex       = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	searchRegex     = regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`)
	allowedImageExt = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)$`)
)

// RegisterCustomValidations registers custom validation rules into Gin
func RegisterCustomValidations() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	// slug validation
	if err := v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	// search validation
	if err := v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	// https only validation
	if err := v.RegisterValidation("https_url", func(fl validator.FieldLevel) bool {
		return strings.HasPrefix(fl.Field().String(), "https://")
	}); err != nil {
		return err
	}

	// image extension validation
	if err := v.RegisterValidation("image_ext", func(fl validator.FieldLevel) bool {
		return allowedImageExt.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	return nil
}

// HandleValidationErrors converts validator errors into readable JSON
func HandleValidationErrors(err error) gin.H {
	errors := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return gin.H{"error": err.Error()}
	}

	for _, e := range validationErrors {

		// Example:
		// createProductRequest.Images[0].Link
		fieldPath := e.StructNamespace()

		// Remove struct name prefix
		// if idx := strings.Index(fieldPath, "."); idx != -1 {
		// 	fieldPath = fieldPath[idx+1:]
		// }

		// Convert to lowercase JSON-style
		fieldPath = toJSONFieldPath(fieldPath)

		switch e.Tag() {
		case "required":
			errors[fieldPath] = "is required"
		case "uuid":
			errors[fieldPath] = "must be a valid UUID"
		case "slug":
			errors[fieldPath] = "must be a valid slug (lowercase, numbers, dash only)"
		case "min":
			errors[fieldPath] = "must be at least " + e.Param()
		case "max":
			errors[fieldPath] = "must not exceed " + e.Param()
		case "oneof":
			values := strings.Split(e.Param(), " ")
			errors[fieldPath] = "must be one of: " + strings.Join(values, ", ")
		case "search":
			errors[fieldPath] = "must contain only letters, numbers, spaces or dash"
		case "gt":
			errors[fieldPath] = "must be greater than " + e.Param()
		case "gte":
			errors[fieldPath] = "must be greater than or equal to " + e.Param()
		case "lt":
			errors[fieldPath] = "must be less than " + e.Param()
		case "lte":
			errors[fieldPath] = "must be less than or equal to " + e.Param()
		default:
			errors[fieldPath] = "is invalid"
		}
	}

	return gin.H{"errors": errors}
}

func toJSONFieldPath(field string) string {
	// Convert:
	// Images[0].Link → images[0].link

	parts := strings.Split(field, ".")
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}

	return strings.Join(parts, ".")
}
