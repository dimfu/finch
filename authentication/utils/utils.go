package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GenerateValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "no_whitespace":
		return fmt.Sprintf("%s must not have whitespace", e.Field())
	case "username_chars":
		return fmt.Sprintf("%s must only contains letters, numbers, _ or -", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", e.Field(), e.Param())
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", e.Field(), e.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is not valid", e.Field())
	}
}
