package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OAuthAccessToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userID" json:"userID"`
	ClientID  string             `bson:"clientID" json:"clientID"`
	GrantType string             `bson:"grantType" json:"grantType"`
	Scopes    []string           `bson:"scopes" json:"scopes"`
	Revoked   int                `bson:"revoked" json:"revoked"`
	ExpiresIn int64              `bson:"expiresIn" json:"expiresIn"`
	Deleted   bool               `bson:"deleted" json:"-"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt int64              `bson:"updatedAt" json:"updatedAt"`
	DeletedAt int64              `bson:"deletedAt" json:"-"`
}

type UpdateOAuthAccessToken struct {
	Revoked   int   `bson:"revoked" json:"revoked"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
}
