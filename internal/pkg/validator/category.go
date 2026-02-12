package validator

var allowedCategories = map[string]bool{
	"php":    true,
	"python": true,
	"golang": true,
}

func IsValidCategory(name string) bool {
	return allowedCategories[name]
}
