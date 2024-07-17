package services

import (
	"fmt"
	"go-account/models"
	"go-account/repositories"
	"go-account/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUsers(ctx *gin.Context) ([]models.User, int64, error)
}

type userService struct {
	validate       *validator.Validate
	userRepository repositories.UserRepository
}

func (s *userService) CreateUser(user *models.User) error {
	if err := s.validate.Struct(user); err != nil {
		return err
	}
	hasheds := utils.HashPassword(user.Password, user.Pincode)
	user.Password = hasheds[0]
	user.Pincode = hasheds[1]
	user.Status = strings.ToLower(user.Status)
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUsers(ctx *gin.Context) ([]models.User, int64, error) {
	filter := make(map[string]interface{})
	if username := ctx.Query("username"); username != "" {
		filter["username"] = username
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	skip, size := utils.PageAndSize(page, size)
	fmt.Println(skip, size)
	return s.userRepository.GetUsers(filter, skip, size)
}

func NewUserService(validate *validator.Validate, userRepository repositories.UserRepository) UserService {
	return &userService{validate, userRepository}
}
