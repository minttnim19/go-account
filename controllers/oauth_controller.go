package controllers

import (
	middy "go-account/middlewares"
	"go-account/models"
	"go-account/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	RefreshToken string `form:"refresh_token"`
	Username     string `form:"username"`
	Password     string `form:"password"`
}

type OAuthController struct {
	oAuthService services.OAuthService
}

func (ctrl *OAuthController) Token(ctx *gin.Context) {
	request, _ := ctx.Get("request")
	clientID, _ := ctx.Get("clientID")
	clientSecret, _ := ctx.Get("clientSecret")
	token, err := ctrl.oAuthService.Token(clientID.(string), clientSecret.(string), request.(middy.OAuthToken))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

// func (ctrl *OAuthController) RefreshToken(ctx *gin.Context) {
// 	userId, _ := ctx.Get("userId")
// 	jti, _ := ctx.Get("jti")
// 	token, err := ctrl.oAuthService.RefreshToken(userId.(string), jti.(string))
// 	if err != nil {
// 		ctx.Error(err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, token)
// }

func (ctrl *OAuthController) CreateOAuthClient(ctx *gin.Context) {
	client := models.CreateOAuthClient{}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.Error(err)
		return
	}

	result, err := ctrl.oAuthService.CreateOAuthClient(&client)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (ctrl *OAuthController) GetOAuthClients(ctx *gin.Context) {
	clients, count, err := ctrl.oAuthService.GetOAuthClients(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
	ctx.JSON(http.StatusOK, clients)
}

func NewOauthController(service services.OAuthService) *OAuthController {
	return &OAuthController{
		oAuthService: service,
	}
}
