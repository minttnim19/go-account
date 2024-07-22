package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required,min=12,max=24"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6,max=16"`
	Pincode   string             `bson:"pincode" json:"pincode" validate:"omitempty,len=6,numeric"`
	Status    string             `bson:"status" json:"status" validate:"required,status"`
	Deleted   bool               `bson:"deleted" json:"deleted"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt int64              `bson:"deletedAt" json:"deletedAt"`
}

type UpdateUser struct {
	Status    string `json:"status" validate:"omitempty,status"`
	UpdatedAt int64  `bson:"updatedAt" json:"updatedAt"`
}
