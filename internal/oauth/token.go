package oauth

import (
	"encoding/base64"
	"errors"
	"go-account/pkg/middlewares"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type OAuthToken struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	RefreshToken string `form:"refresh_token"`
	Username     string `form:"username"`
	Password     string `form:"password"`
}

type OAuthRevoke struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint" binding:"required"`
}

func Revoke() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := OAuthRevoke{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, middlewares.AppError{
				Error:   "Bad Request",
				Message: "the request is missing a required parameter",
			})
			return
		}
		if err := extractClientCredentialsFromHeader(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, middlewares.AppError{
				Error:   "Bad Request",
				Message: err.Error(),
			})
			return
		}
		ctx.Set("request", request)
		ctx.Next()
	}
}

func Token() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := OAuthToken{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, middlewares.AppError{
				Error:   "Bad Request",
				Message: "the request is missing a required parameter",
			})
			return
		}

		if request.GrantType == "client_credentials" {
			if err := extractClientCredentialsFromHeader(ctx); err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, middlewares.AppError{
					Error:   "Bad Request",
					Message: err.Error(),
				})
				return
			}
		} else {
			ctx.Set("clientID", ctx.PostForm("client_id"))
			ctx.Set("clientSecret", ctx.PostForm("client_secret"))
		}

		ctx.Set("request", request)
		ctx.Next()
	}
}

func extractClientCredentialsFromHeader(ctx *gin.Context) error {
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		return errors.New("authorization header required")
	}

	encodedCredentials := strings.TrimPrefix(auth, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return errors.New("invalid authorization header format")
	}

	credentials := strings.SplitN(string(decodedBytes), ":", 2)
	if len(credentials) != 2 || credentials[0] == "" || credentials[1] == "" {
		return errors.New("invalid credentials format")
	}

	ctx.Set("clientID", credentials[0])
	ctx.Set("clientSecret", credentials[1])

	return nil
}
