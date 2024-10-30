package usecase

import (
	"go-account/internal/model"
	"go-account/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
	CreateUser(user *model.CreateUser) error
	GetUsers(ctx *gin.Context) ([]model.User, int64, error)
	GetUserByID(id string) (model.User, error)
	UpdateUser(id string, user *model.UpdateUser) error
	DeleteUser(id string) error
}

type userUsecase struct {
	userRepository model.UserRepository
}

func (s *userUsecase) CreateUser(user *model.CreateUser) error {
	hasheds := utils.HashPassword(user.Password)
	user.Password = hasheds[0]
	user.Status = strings.ToLower(user.Status)
	return s.userRepository.Create(user)
}

func (s *userUsecase) GetUsers(ctx *gin.Context) ([]model.User, int64, error) {
	filter := make(map[string]interface{})
	if username := ctx.Query("username"); username != "" {
		filter["username"] = username
	}
	if status := ctx.Query("status"); status != "" {
		filter["status"] = status
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	skip, size := utils.PageAndSize(page, size)
	return s.userRepository.Lists(filter, skip, size)
}

func (s *userUsecase) GetUserByID(id string) (model.User, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.FindByID(objectId)
}

func (s *userUsecase) UpdateUser(id string, user *model.UpdateUser) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.Update(objectId, user)
}

func (s *userUsecase) DeleteUser(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.Delete(objectId)
}

func NewUserUsecase(userRepository model.UserRepository) UserUsecase {
	return &userUsecase{userRepository}
}
