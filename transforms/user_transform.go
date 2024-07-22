package transforms

import (
	"go-account/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Pincode   string             `json:"pincode"`
	Status    string             `json:"status"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt"`
}

func TransformUser(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Password:  "******",
		Pincode:   "******",
		Status:    user.Status,
		CreatedAt: time.Unix(user.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt: time.Unix(user.UpdatedAt, 0).Format(time.RFC3339),
	}
}

func TransformUsers(users []models.User) []UserResponse {
	userResponses := []UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, TransformUser(user))
	}
	return userResponses
}
