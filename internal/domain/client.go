package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuthClient struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name" binding:"required" validate:"min=1,max=64"`
	Secret     string             `bson:"secret" json:"secret"`
	Redirects  []string           `bson:"redirects" json:"redirects" validate:"omitempty,dive,url"`
	Scopes     []string           `bson:"scopes" json:"scopes" validate:"omitempty,dive,required"`
	GrantTypes []string           `bson:"grantTypes" json:"grantTypes" validate:"omitempty,grant_types"`
	Revoked    int                `bson:"revoked" json:"revoked"`
	Deleted    bool               `bson:"deleted" json:"-"`
	CreatedAt  int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt  int64              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  int64              `bson:"deletedAt" json:"-"`
}

type UpdateOAuthClient struct {
	Name       *string   `bson:"name" json:"name" validate:"omitempty,min=1,max=64"`
	Redirects  *[]string `bson:"redirects" json:"redirects" validate:"omitempty,dive,url"`
	Scopes     *[]string `bson:"scopes" json:"scopes" validate:"omitempty,dive,required"`
	GrantTypes *[]string `bson:"grantTypes" json:"grantTypes" validate:"omitempty,grant_types"`
	Revoked    *int      `bson:"revoked" json:"revoked"`
	UpdatedAt  *int64    `bson:"updatedAt" json:"updatedAt"`
}
