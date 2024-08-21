package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func StatusValidator(fl validator.FieldLevel) bool {
	status := strings.ToLower(fl.Field().String())
	return status == "active" || status == "inactive"
}

func GrantTypesValidator(fl validator.FieldLevel) bool {
	grantTypes := fl.Field().Interface().([]string)
	validTypes := map[string]bool{
		"authorization_code": true,
		"password":           true,
		"client_credentials": true,
		"refresh_token":      true,
	}

	for _, grantType := range grantTypes {
		if !validTypes[grantType] {
			return false
		}
	}
	return true
}
