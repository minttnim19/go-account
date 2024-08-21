package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "Authorization header required",
			})
			return
		}

		encodedCredentials := strings.TrimPrefix(auth, "Basic ")
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
			return
		}

		credentials := strings.SplitN(string(decodedBytes), ":", 2)
		if len(credentials) != 2 || credentials[0] == "" || credentials[1] == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "Invalid credentials format",
			})
			return
		}

		username, password := credentials[0], credentials[1]

		ctx.Set("username", username)
		ctx.Set("password", password)
		ctx.Next()
	}
}
