package routes

import (
	"go-account/controllers"
	"go-account/repositories"
	"go-account/services"
	"go-account/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.RouterGroup, db *mongo.Database) {
	validate := validator.New()
	validate.RegisterValidation("status", utils.StatusValidator)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(validate, userRepository)
	userController := controllers.NewUserController(userService)

	r.POST("/users", userController.CreateUser)
	r.GET("/users", userController.GetUsers)
	// r.GET("/users/:id", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	// })
	r.PATCH("/users/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	r.DELETE("/users/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
}
