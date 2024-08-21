package utils

import "net/http"

type CustomError struct {
	Code    int
	Err     string
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewErrorUnauthorized(message string) error {
	return &CustomError{http.StatusUnauthorized, "Unauthorized", message}
}

func NewErrorBadRequest(message string) error {
	return &CustomError{http.StatusBadRequest, "Bad Request", message}
}
