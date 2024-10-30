package oauth

import (
	"crypto/rsa"
	"fmt"
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

type JWTIssue struct {
	privateKey         *rsa.PrivateKey
	expireToken        string
	expireRefreshToken string
}

func NewJWTIssue() *JWTIssue {
	privateKeyData, err := os.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		panic(fmt.Errorf("failed to load private key: %w", err))
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}
	expireToken := os.Getenv("PRIVATE_KEY_PATH")
	expireRefreshToken := os.Getenv("PRIVATE_KEY_PATH")
	return &JWTIssue{privateKey, expireToken, expireRefreshToken}
}

func (jwtIssue *JWTIssue) GenerateClientToken(scopes []string, grantType, clientID, tokenID string) (string, error) {
	expiresAt := expirationTime(jwtIssue.expireToken)
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jwtIssue.privateKey)
}

func (jwtIssue *JWTIssue) GenerateToken(identity string, scopes []string, grantType, clientID, tokenID string, originJTI *string) (string, error) {
	expiresAt := expirationTime(jwtIssue.expireToken)
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jwtIssue.privateKey)
}

func (jwtIssue *JWTIssue) GenerateRefreshToken(identity string, clientID, refreshTokenID string) (string, error) {
	expiresAt := expirationTime(jwtIssue.expireRefreshToken)
	claims := &Claims{
		ClientID: clientID,
		UserId:   identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        refreshTokenID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jwtIssue.privateKey)
}

func expirationTime(second string) time.Duration {
	expire, err := strconv.Atoi(second)
	if err != nil {
		return 24 * time.Hour
	}
	return time.Duration(expire) * time.Second
}
