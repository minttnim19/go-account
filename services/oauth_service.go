package services

import (
	"fmt"
	middy "go-account/middlewares"
	"go-account/models"
	"go-account/repositories"
	"go-account/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OAuthService interface {
	Token(clientID string, clientSecret string, request middy.OAuthToken) (*TokenDetails, error)
	// RefreshToken(userId string, jti string) (*TokenDetails, error)

	CreateOAuthClient(client *models.CreateOAuthClient) (models.OAuthClient, error)
	GetOAuthClients(ctx *gin.Context) ([]models.OAuthClient, int64, error)
}

type TokenDetails struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type oAuthService struct {
	validate               *validator.Validate
	userRepository         repositories.UserRepository
	clientRepository       repositories.OAuthClientRepository
	tokenRepository        repositories.OAuthAccessTokenRepository
	refreshTokenRepository repositories.OAuthRefreshTokenRepository
}

func (s *oAuthService) Token(clientID string, clientSecret string, request middy.OAuthToken) (*TokenDetails, error) {
	ocid, _ := primitive.ObjectIDFromHex(clientID)
	client, err := s.clientRepository.FindByID(ocid)
	if err != nil || clientSecret != client.Secret || !utils.InSlice(client.GrantTypes, request.GrantType) {
		return nil, utils.NewErrorUnauthorized("Client authentication failed")
	}
	if client.Revoked == 1 {
		return nil, utils.NewErrorUnauthorized("Client authentication has been revoked")
	}

	switch request.GrantType {
	case "password":
		return s.password(&client, &request)
	case "client_credentials":
		return s.clientCredentials(&client, &request)
	default:
		return nil, utils.NewErrorUnauthorized("The authorization grant type is not supported by the authorization server.")
	}
}

func (s *oAuthService) CreateOAuthClient(client *models.CreateOAuthClient) (models.OAuthClient, error) {
	// Validate the client
	if err := s.validate.Struct(client); err != nil {
		return models.OAuthClient{}, err
	}

	// Set default value
	client.Redirects = utils.DefaultStringSlice(client.Redirects)
	client.Scopes = utils.DefaultStringSlice(client.Scopes)
	client.GrantTypes = utils.DefaultStringSlice(client.GrantTypes)

	// Generate client secret
	secret, err := utils.GenerateSecretBase64(32)
	if err != nil {
		return models.OAuthClient{}, err
	}
	client.Secret = secret

	// Create the client in the repository
	result, err := s.clientRepository.Create(client)
	if err != nil {
		return models.OAuthClient{}, err
	}

	// Extract ObjectID and find client by ID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.OAuthClient{}, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return s.clientRepository.FindByID(oid)
}

func (s *oAuthService) GetOAuthClients(ctx *gin.Context) ([]models.OAuthClient, int64, error) {
	filter := make(map[string]interface{})
	if username := ctx.Query("username"); username != "" {
		filter["username"] = username
	}
	if status := ctx.Query("status"); status != "" {
		filter["status"] = status
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	skip, size := utils.PageAndSize(page, size)
	return s.clientRepository.Lists(filter, skip, size)
}

func NewOauthService(validate *validator.Validate, userRepository repositories.UserRepository, clientRepository repositories.OAuthClientRepository, tokenRepository repositories.OAuthAccessTokenRepository, refreshTokenRepository repositories.OAuthRefreshTokenRepository) OAuthService {
	return &oAuthService{validate, userRepository, clientRepository, tokenRepository, refreshTokenRepository}
}

func (s *oAuthService) clientCredentials(client *models.OAuthClient, request *middy.OAuthToken) (*TokenDetails, error) {
	// merge client scopes and user scopes
	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes)

	// generate access token for the authenticated user
	tokenExpiresIn := utils.GetTokenExpireTime()
	tokenCreationResult, err := s.tokenRepository.Create(&models.OAuthAccessToken{
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	})
	if err != nil {
		return nil, err
	}
	// generate access token and refresh token for the authenticated user
	tokenID, ok := tokenCreationResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, utils.NewErrorBadRequest("Failed to parse token ID")
	}
	accessToken, err := utils.GenerateClientToken(scopes, request.GrantType, client.ID.Hex(), tokenID.Hex())
	if err != nil {
		return nil, err
	}
	return &TokenDetails{
		TokenType:   "Bearer",
		ExpiresIn:   tokenExpiresIn,
		AccessToken: accessToken,
	}, nil
}

func (s *oAuthService) password(client *models.OAuthClient, request *middy.OAuthToken) (*TokenDetails, error) {
	user, err := s.userRepository.FindUserByUsername(request.Username)
	if err != nil {
		return nil, utils.NewErrorBadRequest("Your username is not found")
	}

	// check if the user password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, utils.NewErrorBadRequest("Your password is incorrect")
	}

	// merge client scopes and user scopes
	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes, []string{"*"})

	// store data access token in oauth_tokens
	tokenExpiresIn := utils.GetTokenExpireTime()
	tokenCreationResult, err := s.tokenRepository.Create(&models.OAuthAccessToken{
		UserID:    user.ID.Hex(),
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	})
	if err != nil {
		return nil, err
	}

	// generate access token and refresh token for the authenticated user
	tokenID, ok := tokenCreationResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, utils.NewErrorBadRequest("Failed to parse token ID")
	}
	accessToken, err := utils.GenerateToken(&user, scopes, request.GrantType, client.ID.Hex(), tokenID.Hex())
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresIn := utils.GetRefreshTokenExpireTime()
	refreshTokenCreationResult, err := s.refreshTokenRepository.Create(&models.OAuthRefreshToken{
		AccessTokenID: tokenID,
		ExpiresIn:     refreshTokenExpiresIn,
	})
	if err != nil {
		return nil, err
	}

	refreshTokenID, ok := refreshTokenCreationResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, utils.NewErrorBadRequest("Failed to parse refresh token ID")
	}
	refreshToken, err := utils.GenerateRefreshToken(&user, client.ID.Hex(), refreshTokenID.Hex(), tokenID.Hex())
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		TokenType:    "Bearer",
		ExpiresIn:    tokenExpiresIn,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
