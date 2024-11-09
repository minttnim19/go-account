package server

import (
	"go-account/config"
	"go-account/internal/handler"
	"go-account/internal/repository"
	"go-account/internal/usecase"
	"go-account/pkg/database"
	"go-account/pkg/middleware"
	"go-account/pkg/oauth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GonicServer struct {
	app  *gin.Engine
	db   database.Database
	conf *config.Config
}

func NewGonicServer(conf *config.Config, db database.Database) *GonicServer {
	r := gin.Default()
	r.Use(gin.Recovery()) // Recovery when system die.

	// Set up routes
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Pong!"})
	})

	return &GonicServer{
		app:  r,
		db:   db,
		conf: conf,
	}
}

func (srv *GonicServer) Start() {
	// Apply middleware
	srv.app.Use(middleware.ErrorHandler())

	srv.initializeUserHttpHandler()  // User handler
	srv.initializeOAuthHttpHandler() // OAuth handler

	if err := srv.app.Run(srv.conf.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func (srv *GonicServer) initializeUserHttpHandler() {
	repo := repository.NewUserRepository(srv.db.GetDb())
	u := usecase.NewUserUsecase(repo)
	h := handler.NewUserHandler(u)

	jwtVerify := oauth.NewJWTVerify()
	users := srv.app.Group("/api/v1/users").Use(middleware.Authenticate(jwtVerify))
	{
		users.POST("/", h.CreateUser)
		users.GET("/", h.GetUsers)
		users.GET("/:id", h.GetUserByID)
		users.PATCH("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

func (srv *GonicServer) initializeOAuthHttpHandler() {
	userRepository := repository.NewUserRepository(srv.db.GetDb())
	clientRepository := repository.NewOAuthClientRepository(srv.db.GetDb())
	tokenRepository := repository.NewOAuthAccessTokenRepository(srv.db.GetDb())
	refreshTokenRepository := repository.NewOAuthRefreshTokenRepository(srv.db.GetDb())

	jwtIssue := oauth.NewJWTIssue()
	jwtVerify := oauth.NewJWTVerify()
	u := usecase.NewOauthUsecase(jwtIssue, jwtVerify, userRepository, clientRepository, tokenRepository, refreshTokenRepository)
	h := handler.NewOAuthHandler(u)

	srv.app.POST("/api/v1/oauth/token", middleware.Token(), h.Token)
	srv.app.POST("/api/v1/oauth/revoke", middleware.Revoke(), h.Revoke)
	oauth := srv.app.Group("/api/v1/oauth/clients")
	{
		oauth.POST("/", h.CreateOAuthClient)
		oauth.GET("/", h.GetOAuthClients)
		oauth.PATCH("/:id", h.UpdateOAuthClient)
	}
}
