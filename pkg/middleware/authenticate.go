package middleware

import (
	"fmt"
	"go-account/pkg/oauth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Identities struct {
	UserId string
}

func Authenticate(jwtVerify *oauth.JWTVerify) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "authorization header required",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtVerify.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: err.Error(),
			})
			return
		}
		fmt.Println("Claims", claims)
		fmt.Println("Subject", claims.Subject)
		ctx.Next()
	}
}
