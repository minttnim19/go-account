package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func StatusValidator(fl validator.FieldLevel) bool {
	status := strings.ToLower(fl.Field().String())
	return status == "active" || status == "inactive"
}
