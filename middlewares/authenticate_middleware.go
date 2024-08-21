package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// token := ctx.GetHeader("Authorization")
		// if token == "" {
		// 	ctx.JSON(http.StatusForbidden, AppError{
		// 		Error:   "Forbidden",
		// 		Message: "You do not have permission to access this resource",
		// 	})
		// 	ctx.Abort()
		// 	return
		// }
		// claims, err := utils.ValidateJWT(tokenStr)
		// if err != nil {
		// 	response, _ := utils.JSON(nil, "Unauthorized", false)
		// 	c.Data(http.StatusUnauthorized, "application/json", response)
		// 	c.Abort()
		// 	return
		// }

		// c.Set("username", claims.Username)
		// c.Set("role", claims.Role)
		ctx.Next()
	}
}

// func validateJWT(tokenStr string) (*Claims, error) {
// 	claims := &Claims{}
// 	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return claims, nil
// }
