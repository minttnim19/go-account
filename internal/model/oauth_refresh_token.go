package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthRefreshToken struct {
	ID            primitive.Binary `bson:"_id,omitempty" json:"id"`
	AccessTokenID primitive.Binary `bson:"accessTokenID" json:"accessTokenID"`
	Revoked       int              `bson:"revoked" json:"revoked"`
	ExpiresIn     int64            `bson:"expiresIn" json:"expiresIn"`
	Deleted       bool             `bson:"deleted" json:"-"`
	CreatedAt     int64            `bson:"createdAt" json:"createdAt"`
	UpdatedAt     int64            `bson:"updatedAt" json:"updatedAt"`
	DeletedAt     int64            `bson:"deletedAt" json:"-"`
}

type UpdateOAuthRefreshToken struct {
	Revoked   int   `bson:"revoked" json:"revoked"`
	UpdatedAt int64 `bson:"updatedAt" json:"updatedAt"`
}

type OAuthRefreshTokenRepository interface {
	Create(token *OAuthRefreshToken) (*mongo.InsertOneResult, error)
	FindByID(id primitive.Binary) (OAuthRefreshToken, error)
	Update(id primitive.Binary, token *UpdateOAuthRefreshToken) error
}
