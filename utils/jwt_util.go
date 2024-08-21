package utils

import (
	"go-account/models"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	GrantType string   `json:"grant_type,omitempty"`
	ClientID  string   `json:"client_id"`
	UserId    string   `json:"user_id,omitempty"`
	Username  string   `json:"username,omitempty"`
	Scope     []string `json:"scope,omitempty"`
	OriginJTI *string  `json:"origin_jti,omitempty"`
	jwt.StandardClaims
}

func GetTokenExpireTime() int64 {
	exp, _ := strconv.ParseInt(GetEnv("TOKEN_EXPIRE_TIME", "86400"), 10, 64)
	return exp
}

func GetRefreshTokenExpireTime() int64 {
	exp, _ := strconv.ParseInt(GetEnv("TOKEN_REFRESH_EXPIRE_TIME", "604800"), 10, 64)
	return exp
}

func GenerateToken(user *models.User, scopes []string, grantType, clientID, tokenID string) (string, error) {
	expiresAt := expirationTime(GetEnv("TOKEN_EXPIRE_TIME", "86400"))
	claims := &Claims{
		GrantType: grantType,
		ClientID:  clientID,
		UserId:    user.ID.Hex(),
		Scope:     scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.ID.Hex(),
			Id:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func GenerateRefreshToken(user *models.User, clientID, refreshTokenID, tokenID string) (string, error) {
	expiresAt := expirationTime(GetEnv("TOKEN_REFRESH_EXPIRE_TIME", "604800"))
	claims := &Claims{
		ClientID: clientID,
		UserId:   user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   tokenID,
			Id:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, NewErrorUnauthorized(err.Error())
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, NewErrorUnauthorized("Invalid token")
}

func expirationTime(second string) time.Duration {
	expire, err := strconv.Atoi(second)
	if err != nil {
		return 24 * time.Hour
	}
	return time.Duration(expire) * time.Second
}
