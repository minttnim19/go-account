package api

import (
	"go-account/internal/api/controllers"
	"go-account/internal/api/repositories"
	"go-account/internal/api/services"
	"go-account/internal/oauth"
	"go-account/internal/validators"
	"go-account/pkg/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(db *mongo.Database) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery()) // Recovery when system die.

	// Apply middlewares
	r.Use(middlewares.ErrorHandler())

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("status", validators.ValidateStatus)
		v.RegisterValidation("grant_types", validators.ValidateGrantTypes)
	}

	// Set up routes
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Pong!"})
	})

	version := r.Group("/api/v1")
	{
		userRoutes(version, db)
		oAuthRoutes(version, db)
	}

	return r
}

func userRoutes(r *gin.RouterGroup, db *mongo.Database) {
	repo := repositories.NewUserRepository(db)
	s := services.NewUserService(repo)
	ctrl := controllers.NewUserController(s)

	users := r.Group("/users").Use(middlewares.Authenticate())
	{
		users.POST("/", ctrl.CreateUser)
		users.GET("/", ctrl.GetUsers)
		users.GET("/:id", ctrl.GetUserByID)
		users.PATCH("/:id", ctrl.UpdateUser)
		users.DELETE("/:id", ctrl.DeleteUser)
	}
}

func oAuthRoutes(r *gin.RouterGroup, db *mongo.Database) {
	userRepository := repositories.NewUserRepository(db)
	clientRepository := repositories.NewOAuthClientRepository(db)
	tokenRepository := repositories.NewOAuthAccessTokenRepository(db)
	refreshTokenRepository := repositories.NewOAuthRefreshTokenRepository(db)
	s := services.NewOauthService(userRepository, clientRepository, tokenRepository, refreshTokenRepository)

	ctrl := controllers.NewOAuthController(s)

	r.POST("/oauth/token", oauth.Token(), ctrl.Token)
	r.POST("/oauth/revoke", oauth.Revoke(), ctrl.Revoke)
	oauth := r.Group("/oauth")
	{
		oauth.POST("/clients", ctrl.CreateOAuthClient)
		oauth.GET("/clients", ctrl.GetOAuthClients)
	}
}
