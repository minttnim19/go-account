package controllers

import (
	"go-account/internal/api/models"
	"go-account/internal/api/services"
	"go-account/internal/oauth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	service services.OAuthService
}

func (ctrl *OAuthController) Token(ctx *gin.Context) {
	request, _ := ctx.Get("request")
	clientID, _ := ctx.Get("clientID")
	clientSecret, _ := ctx.Get("clientSecret")
	token, err := ctrl.service.Token(clientID.(string), clientSecret.(string), request.(oauth.OAuthToken))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (ctrl *OAuthController) Revoke(ctx *gin.Context) {
	request, _ := ctx.Get("request")
	clientID, _ := ctx.Get("clientID")
	clientSecret, _ := ctx.Get("clientSecret")
	if err := ctrl.service.Revoke(clientID.(string), clientSecret.(string), request.(oauth.OAuthRevoke)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Token revoked!"})
}

func (ctrl *OAuthController) CreateOAuthClient(ctx *gin.Context) {
	client := models.CreateOAuthClient{}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.Error(err)
		return
	}

	result, err := ctrl.service.CreateOAuthClient(&client)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (ctrl *OAuthController) GetOAuthClients(ctx *gin.Context) {
	clients, count, err := ctrl.service.GetOAuthClients(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
	ctx.JSON(http.StatusOK, clients)
}

func NewOAuthController(service services.OAuthService) *OAuthController {
	return &OAuthController{
		service: service,
	}
}
