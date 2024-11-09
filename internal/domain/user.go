package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" binding:"required" validate:"min=12,max=24"`
	Password  string             `bson:"password" json:"password" binding:"required" validate:"min=6,max=16"`
	Status    string             `bson:"status" json:"status" binding:"required" validate:"status"`
	Deleted   bool               `bson:"deleted" json:"-"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt int64              `bson:"deletedAt" json:"-"`
}

type UpdateUser struct {
	Status    *string `bson:"status" json:"status" validate:"omitempty,status"`
	UpdatedAt *int64  `bson:"updatedAt" json:"updatedAt"`
}
