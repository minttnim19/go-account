package middlewares

import (
	"go-account/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "Authorization header required",
			})
			return
		}

		tokenString := strings.TrimPrefix(token, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: err.Error(),
			})
			return
		}

		if claims.UserId == "" || claims.Id == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, AppError{
				Error:   "Unauthorized",
				Message: "Invalid refresh token format",
			})
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Set("jti", claims.Id)
		ctx.Next()
	}
}
