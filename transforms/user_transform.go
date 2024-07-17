package transforms

import (
	"go-account/models"
	"time"
)

func TransformUser(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Password:  "******",
		Pincode:   "******",
		Status:    user.Status,
		CreatedAt: time.Unix(user.CreatedAt, 0).Format(time.RFC3339),
		UpdatedAt: time.Unix(user.UpdatedAt, 0).Format(time.RFC3339),
	}
}

func TransformUsers(users []models.User) []models.UserResponse {
	userResponses := []models.UserResponse{}
	for _, user := range users {
		userResponses = append(userResponses, TransformUser(user))
	}
	return userResponses
}
