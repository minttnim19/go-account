package middlewares

import (
	"encoding/json"
	"go-account/pkg/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Code    int      `json:"-"`
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Details []Detail `json:"details,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
			appError := handleAppError(err)
			ctx.JSON(appError.Code, appError)
		}
	}
}

func handleAppError(err error) *AppError {
	appError := &AppError{}

	switch err := err.(type) {
	case validator.ValidationErrors:
		handleValidationError(appError, err)
	case mongo.CommandError:
		handleCommandError(appError, err)
	case mongo.WriteException:
		handleWriteException(appError, err)
	case *utils.CustomError:
		handleCustomError(appError, *err)
	case *json.UnmarshalTypeError:
		handleUnmarshalTypeError(appError, err)
	default:
		handleDefaultError(appError, err)
	}

	return appError
}

func handleUnmarshalTypeError(appError *AppError, err *json.UnmarshalTypeError) {
	fieldName := strings.ToUpper(string(err.Field[0])) + err.Field[1:]
	errorMessages := []Detail{{
		Field:   err.Field,
		Message: fieldName + " must be type " + err.Type.String(),
	}}

	*appError = AppError{
		Code:    http.StatusUnprocessableEntity,
		Error:   "Unprocessable Entity",
		Message: "validation failed",
		Details: errorMessages,
	}
}

func handleValidationError(appError *AppError, err validator.ValidationErrors) {
	errorMessages := []Detail{}
	for _, fieldError := range err {
		errorMessage := formatErrorDetail(fieldError)
		errorMessages = append(errorMessages, errorMessage)
	}
	*appError = AppError{
		Code:    http.StatusUnprocessableEntity,
		Error:   "Unprocessable Entity",
		Message: "validation failed",
		Details: errorMessages,
	}
}

func handleCustomError(appError *AppError, err utils.CustomError) {
	*appError = AppError{
		Code:    err.Code,
		Error:   err.Err,
		Message: err.Error(),
	}
}

func handleCommandError(appError *AppError, err mongo.CommandError) {
	*appError = AppError{
		Code:    http.StatusBadRequest,
		Error:   "Bad Request",
		Message: err.Error(),
	}
}

func handleWriteException(appError *AppError, err mongo.WriteException) {
	if err.WriteErrors[0].Code == 11000 {
		*appError = AppError{
			Code:    http.StatusConflict,
			Error:   "Conflict",
			Message: parseDuplicateKeyError(err.Error()),
		}
	} else {
		*appError = AppError{
			Code:    http.StatusBadRequest,
			Error:   "Bad Request",
			Message: err.Error(),
		}
	}
}

func handleDefaultError(appError *AppError, err error) {
	if err == mongo.ErrNoDocuments {
		*appError = AppError{
			Code:    http.StatusNotFound,
			Error:   "Not Found",
			Message: "the requested resource was not found",
		}
	} else {
		*appError = AppError{
			Code:    http.StatusInternalServerError,
			Error:   "Internal Server Error",
			Message: err.Error(),
		}
	}
}

func parseDuplicateKeyError(errorMessage string) string {
	re := regexp.MustCompile(`dup key: { (\w+): "(.+)" }`)
	matches := re.FindStringSubmatch(errorMessage)
	if len(matches) > 2 {
		return matches[2] + " already exists"
	}
	return "the request could not be completed due to a conflict with the current state of the resource"
}

func formatErrorDetail(fe validator.FieldError) Detail {
	fieldName := strings.ToLower(string(fe.Field()[0])) + fe.Field()[1:]
	message := ""

	switch fe.Tag() {
	case "required":
		message = "is required"
	case "min":
		message = "must be at least " + fe.Param() + " characters"
	case "max":
		message = "must be at most " + fe.Param() + " characters"
	case "len":
		message = "must be " + fe.Param() + " characters long"
	case "numeric":
		message = "must be a number"
	case "oneof":
		message = "must be one of " + fe.Param()
	case "status":
		message = "must be one of active or inactive"
	case "url":
		message = "is invalid url format"
	case "grant_types":
		message = "must be one or more in authorization_code ,password ,client_credentials ,refresh_token"
	default:
		message = "is invalid"
	}

	return Detail{
		Field:   fieldName,
		Message: fe.Field() + " " + message,
	}
}
