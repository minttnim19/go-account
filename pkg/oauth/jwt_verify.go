package oauth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type JWTVerify struct {
	publicKey *rsa.PublicKey
}

func NewJWTVerify() *JWTVerify {
	publicKeyData, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		panic(fmt.Errorf("failed to load public key: %w", err))
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		panic(fmt.Errorf("failed to parse public key: %w", err))
	}
	return &JWTVerify{publicKey}
}

func (jwtVerify *JWTVerify) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure token uses RSA signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtVerify.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
