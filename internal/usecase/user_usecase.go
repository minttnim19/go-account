package usecase

import (
	"go-account/internal/domain"
	"go-account/internal/repository"
	"go-account/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	GetUsers(ctx *gin.Context) ([]domain.User, int64, error)
	GetUserByID(id string) (domain.User, error)
	UpdateUser(id string, user *domain.UpdateUser) error
	DeleteUser(id string) error
}

type userUsecase struct {
	userRepository repository.UserRepository
	validator      *validator.Validate
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	if err := u.validator.Struct(user); err != nil {
		return err
	}
	hasheds := utils.HashPassword(user.Password)
	user.Password = hasheds[0]
	user.Status = strings.ToLower(user.Status)
	return u.userRepository.Create(user)
}

func (u *userUsecase) GetUsers(ctx *gin.Context) ([]domain.User, int64, error) {
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
	return u.userRepository.Lists(filter, skip, size)
}

func (u *userUsecase) GetUserByID(id string) (domain.User, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return u.userRepository.FindByID(objectId)
}

func (u *userUsecase) UpdateUser(id string, user *domain.UpdateUser) error {
	if err := u.validator.Struct(user); err != nil {
		return err
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	return u.userRepository.Update(objectId, user)
}

func (u *userUsecase) DeleteUser(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return u.userRepository.Delete(objectId)
}

func NewUserUsecase(userRepository repository.UserRepository) UserUsecase {
	validate := validator.New()
	validate.RegisterValidation("status", utils.ValidateStatus)
	return &userUsecase{userRepository, validate}
}
