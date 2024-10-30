package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func GenerateClientToken(scopes []string, grantType, clientID, tokenID, expTime string) (string, error) {
	expiresAt := expirationTime(expTime)
	claims := &Claims{
		GrantType: grantType,
		ClientID:  clientID,
		Scope:     scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   clientID,
			Id:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func GenerateToken(identity string, scopes []string, grantType, clientID, tokenID, expTime string, originJTI *string) (string, error) {
	expiresAt := expirationTime(expTime)
	claims := &Claims{
		GrantType: grantType,
		ClientID:  clientID,
		Scope:     scopes,
		UserId:    identity,
		OriginJTI: originJTI,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   identity,
			Id:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func GenerateRefreshToken(identity string, clientID, refreshTokenID, expTime string) (string, error) {
	expiresAt := expirationTime(expTime)
	claims := &Claims{
		ClientID: clientID,
		UserId:   identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        refreshTokenID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ValidateJWT(tokenString string) (*Claims, error) {
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
