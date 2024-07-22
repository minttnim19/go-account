package routes

import (
	"go-account/controllers"
	"go-account/repositories"
	"go-account/services"
	"go-account/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(r *gin.RouterGroup, db *mongo.Database) {
	validate := validator.New()
	validate.RegisterValidation("status", utils.StatusValidator)

	repo := repositories.NewUserRepository(db)
	s := services.NewUserService(validate, repo)
	ctrl := controllers.NewUserController(s)

	r.POST("/users", ctrl.CreateUser)
	r.GET("/users", ctrl.GetUsers)
	r.GET("/users/:id", ctrl.GetUserByID)
	r.PATCH("/users/:id", ctrl.UpdateUser)
	r.DELETE("/users/:id", ctrl.DeleteUser)
}
