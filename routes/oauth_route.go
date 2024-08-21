package routes

import (
	"go-account/controllers"
	"go-account/middlewares"
	"go-account/repositories"
	"go-account/services"
	"go-account/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func OAuthRoutes(r *gin.RouterGroup, db *mongo.Database) {
	validate := validator.New()
	validate.RegisterValidation("status", utils.StatusValidator)
	validate.RegisterValidation("grant_types", utils.GrantTypesValidator)

	userRepo := repositories.NewUserRepository(db)
	clientRepo := repositories.NewOAuthClientRepository(db)
	tokenRepo := repositories.NewOAuthAccessTokenRepository(db)
	refreshTokenRepo := repositories.NewOAuthRefreshTokenRepository(db)
	s := services.NewOauthService(validate, userRepo, clientRepo, tokenRepo, refreshTokenRepo)

	ctrl := controllers.NewOauthController(s)

	// oauth := r.Group("/oauth").Use(middlewares.BasicAuth())
	// {
	// 	oauth.POST("/token", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	// 	})
	// 	oauth.POST("/token/refresh", func(ctx *gin.Context) {
	// 		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	// 	})
	// }
	r.POST("/oauth/token", middlewares.Token(), ctrl.Token)
	// r.POST("/oauth/refresh-token", middlewares.RefreshToken(), ctrl.RefreshToken)
	oauth := r.Group("/oauth")
	{
		oauth.POST("/clients", ctrl.CreateOAuthClient)
		oauth.GET("/clients", ctrl.GetOAuthClients)
	}
}
