package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuthClient struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Secret     string             `bson:"secret" json:"secret"`
	Redirects  []string           `bson:"redirects" json:"redirects"`
	Scopes     []string           `bson:"scopes" json:"scopes"`
	GrantTypes []string           `bson:"grantTypes" json:"grantTypes"`
	Revoked    int                `bson:"revoked" json:"revoked"`
	Deleted    bool               `bson:"deleted" json:"-"`
	CreatedAt  int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt  int64              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  int64              `bson:"deletedAt" json:"-"`
}

type CreateOAuthClient struct {
	Name       string   `bson:"name" json:"name" binding:"required,min=1,max=64"`
	Secret     string   `bson:"secret" json:"secret"`
	Redirects  []string `bson:"redirects" json:"redirects" binding:"omitempty,dive,url"`
	Scopes     []string `bson:"scopes" json:"scopes" binding:"omitempty,dive,required"`
	GrantTypes []string `bson:"grantTypes" json:"grantTypes" binding:"omitempty,grant_types"`
	Revoked    int      `bson:"revoked" json:"revoked"`
	Deleted    bool     `bson:"deleted" json:"deleted"`
	CreatedAt  int64    `bson:"createdAt" json:"createdAt"`
	UpdatedAt  int64    `bson:"updatedAt" json:"updatedAt"`
	DeletedAt  int64    `bson:"deletedAt" json:"deletedAt"`
}
