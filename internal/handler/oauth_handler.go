package handler

import (
	"go-account/internal/model"
	"go-account/internal/usecase"
	"go-account/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	usecase usecase.OAuthUsecase
}

func (ctrl *OAuthHandler) Token(ctx *gin.Context) {
	request, _ := ctx.Get("request")
	clientID, _ := ctx.Get("clientID")
	clientSecret, _ := ctx.Get("clientSecret")
	token, err := ctrl.usecase.Token(clientID.(string), clientSecret.(string), request.(middleware.OAuthToken))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func (ctrl *OAuthHandler) Revoke(ctx *gin.Context) {
	request, _ := ctx.Get("request")
	clientID, _ := ctx.Get("clientID")
	clientSecret, _ := ctx.Get("clientSecret")
	if err := ctrl.usecase.Revoke(clientID.(string), clientSecret.(string), request.(middleware.OAuthRevoke)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Token revoked!"})
}

func (ctrl *OAuthHandler) CreateOAuthClient(ctx *gin.Context) {
	client := model.CreateOAuthClient{}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.Error(err)
		return
	}

	result, err := ctrl.usecase.CreateOAuthClient(&client)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (ctrl *OAuthHandler) GetOAuthClients(ctx *gin.Context) {
	clients, count, err := ctrl.usecase.GetOAuthClients(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
	ctx.JSON(http.StatusOK, clients)
}

func NewOAuthHandler(usecase usecase.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{
		usecase: usecase,
	}
}
