package controllers

import (
	"go-account/models"
	"go-account/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	user := models.CreateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		return
	}

	if err := ctrl.userService.CreateUser(&user); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *UserController) GetUsers(ctx *gin.Context) {
	users, count, err := ctrl.userService.GetUsers(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
	ctx.JSON(http.StatusOK, users)
}

func (ctrl *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := ctrl.userService.GetUserByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user := models.UpdateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		return
	}

	if err := ctrl.userService.UpdateUser(id, &user); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ctrl *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.userService.DeleteUser(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}
