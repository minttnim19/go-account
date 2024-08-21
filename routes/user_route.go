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

	users := r.Group("/users")
	{
		users.POST("/", ctrl.CreateUser)
		users.GET("/", ctrl.GetUsers)
		users.GET("/:id", ctrl.GetUserByID)
		users.PATCH("/:id", ctrl.UpdateUser)
		users.DELETE("/:id", ctrl.DeleteUser)
	}
}
