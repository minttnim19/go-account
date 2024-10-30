package handler

import (
	"go-account/internal/model"
	"go-account/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func (ctrl *UserHandler) CreateUser(ctx *gin.Context) {
	user := model.CreateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		return
	}

	if err := ctrl.usecase.CreateUser(&user); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *UserHandler) GetUsers(ctx *gin.Context) {
	users, count, err := ctrl.usecase.GetUsers(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Header("x-total-count", strconv.FormatInt(count, 10))
	ctx.JSON(http.StatusOK, users)
}

func (ctrl *UserHandler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := ctrl.usecase.GetUserByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (ctrl *UserHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user := model.UpdateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.Error(err)
		return
	}

	if err := ctrl.usecase.UpdateUser(id, &user); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ctrl *UserHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.usecase.DeleteUser(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}
