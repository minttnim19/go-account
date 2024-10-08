package services

import (
	"go-account/models"
	"go-account/repositories"
	"go-account/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(user *models.CreateUser) error
	GetUsers(ctx *gin.Context) ([]models.User, int64, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(id string, user *models.UpdateUser) error
	DeleteUser(id string) error
}

type userService struct {
	validate       *validator.Validate
	userRepository repositories.UserRepository
}

func (s *userService) CreateUser(user *models.CreateUser) error {
	if err := s.validate.Struct(user); err != nil {
		return err
	}
	hasheds := utils.HashPassword(user.Password)
	user.Password = hasheds[0]
	user.Status = strings.ToLower(user.Status)
	return s.userRepository.Create(user)
}

func (s *userService) GetUsers(ctx *gin.Context) ([]models.User, int64, error) {
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

func (s *userService) GetUserByID(id string) (models.User, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.FindByID(objectId)
}

func (s *userService) UpdateUser(id string, user *models.UpdateUser) error {
	if err := s.validate.Struct(user); err != nil {
		return err
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.Update(objectId, user)
}

func (s *userService) DeleteUser(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return s.userRepository.Delete(objectId)
}

func NewUserService(validate *validator.Validate, userRepository repositories.UserRepository) UserService {
	return &userService{validate, userRepository}
}
