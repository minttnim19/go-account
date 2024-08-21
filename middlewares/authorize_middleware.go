package middlewares

import (
	"go-account/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorize(conditions ...func(ctx *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, condition := range conditions {
			if condition(ctx) {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, AppError{
			Error:   "Forbidden",
			Message: "You do not have scope to access this resource",
		})
	}
}

func HasRole(roles []string) func(ctx *gin.Context) bool {
	return func(ctx *gin.Context) bool {
		existingRole, exists := ctx.Get("role")
		if !exists || !utils.InSlice(roles, existingRole.(string)) {
			return false
		}
		return true
	}
}

func HasScope(scopes []string) func(ctx *gin.Context) bool {
	return func(ctx *gin.Context) bool {
		existingScopes, exists := ctx.Get("scopes")
		if !exists || !hasScope(scopes, existingScopes.([]string)) {
			return false
		}
		return true
	}
}

func hasScope(scopes []string, items []string) bool {
	for _, scope := range scopes {
		if checked := scopeMatch(items, scope); !checked {
			return checked
		}
	}
	return true
}

func scopeMatch(items []string, searchString string) bool {
	for _, item := range items {
		if item == searchString {
			return true
		} else if strings.HasSuffix(item, ".*") {
			prefix := strings.TrimSuffix(item, ".*")
			if strings.HasPrefix(searchString, prefix) {
				return true
			}
		}
	}
	return false
}
