package usecase

import (
	"errors"
	"go-account/config"

	"go-account/internal/model"

	"go-account/pkg/middleware"
	"go-account/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OAuthUsecase interface {
	Token(clientID string, clientSecret string, request middleware.OAuthToken) (*TokenDetails, error)
	Revoke(clientID string, clientSecret string, request middleware.OAuthRevoke) error
	CreateOAuthClient(client *model.CreateOAuthClient) (model.OAuthClient, error)
	GetOAuthClients(ctx *gin.Context) ([]model.OAuthClient, int64, error)
}

type TokenDetails struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type oAuthUsecase struct {
	userRepository         model.UserRepository
	clientRepository       model.OAuthClientRepository
	tokenRepository        model.OAuthAccessTokenRepository
	refreshTokenRepository model.OAuthRefreshTokenRepository
}

var (
	tokenExpireTime        string
	tokenRefreshExpireTime string
)

func (s *oAuthUsecase) Token(clientID string, clientSecret string, request middleware.OAuthToken) (*TokenDetails, error) {
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
	case "refresh_token":
		return s.refreshToken(&client, &request)
	default:
		return nil, utils.NewErrorUnauthorized("The authorization grant type is not supported by the authorization server.")
	}
}

func (s *oAuthUsecase) Revoke(clientID string, clientSecret string, request middleware.OAuthRevoke) error {
	ocid, _ := primitive.ObjectIDFromHex(clientID)
	client, err := s.clientRepository.FindByID(ocid)
	if err != nil || clientSecret != client.Secret {
		return utils.NewErrorUnauthorized("Client authentication failed")
	}
	if client.Revoked == 1 {
		return utils.NewErrorUnauthorized("Client authentication has been revoked")
	}
	switch request.TokenTypeHint {
	case "refresh_token":
		return s.revokeRefreshToken(&client, &request)
	// case "access_token":
	default:
		return utils.NewErrorUnauthorized("The token_type_hint is not supported by the authorization server.")
	}
}

func (s *oAuthUsecase) CreateOAuthClient(client *model.CreateOAuthClient) (model.OAuthClient, error) {
	// Set default value
	client.Redirects = utils.DefaultStringSlice(client.Redirects)
	client.Scopes = utils.DefaultStringSlice(client.Scopes)
	client.GrantTypes = utils.DefaultStringSlice(client.GrantTypes)

	// Generate client secret
	secret, err := utils.GenerateSecretBase64(32)
	if err != nil {
		return model.OAuthClient{}, err
	}
	client.Secret = secret

	// Create the client in the repository
	result, err := s.clientRepository.Create(client)
	if err != nil {
		return model.OAuthClient{}, err
	}

	// Extract ObjectID and find client by ID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return model.OAuthClient{}, errors.New("failed to convert inserted ID to ObjectID")
	}

	return s.clientRepository.FindByID(oid)
}

func (s *oAuthUsecase) GetOAuthClients(ctx *gin.Context) ([]model.OAuthClient, int64, error) {
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

func NewOauthUsecase(conf *config.Config,
	userRepository model.UserRepository,
	clientRepository model.OAuthClientRepository,
	tokenRepository model.OAuthAccessTokenRepository,
	refreshTokenRepository model.OAuthRefreshTokenRepository) OAuthUsecase {
	tokenExpireTime = conf.TokenExpireTime
	tokenRefreshExpireTime = conf.TokenRefreshExpireTime
	return &oAuthUsecase{userRepository, clientRepository, tokenRepository, refreshTokenRepository}
}

func (s *oAuthUsecase) clientCredentials(client *model.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes)
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &model.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := s.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateClientToken(scopes, request.GrantType, client.ID.Hex(), tokenID, tokenExpireTime)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		TokenType:   "Bearer",
		ExpiresIn:   tokenExpiresIn,
		AccessToken: accessToken,
	}, nil
}

func (s *oAuthUsecase) password(client *model.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	user, err := s.userRepository.FindUserByUsername(request.Username)
	if err != nil {
		return nil, utils.NewErrorBadRequest("Your username is not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, utils.NewErrorBadRequest("Your password is incorrect")
	}

	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes, []string{})
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &model.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		UserID:    user.ID.Hex(),
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := s.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateToken(user.ID.Hex(), scopes, request.GrantType, client.ID.Hex(), tokenID, tokenExpireTime, nil)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateAndStoreRefreshToken(tokenID, user.ID.Hex(), client.ID.Hex())
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

func (s *oAuthUsecase) refreshToken(client *model.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	claims, err := utils.ValidateJWT(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	oAuthRefreshToken, err := s.refreshTokenRepository.FindByID(utils.StringToBinaryUUID(claims.Id))
	if err != nil || oAuthRefreshToken.Revoked == 1 {
		return nil, utils.NewErrorUnauthorized("Refresh Token has been revoked or invalid")
	}

	oAuthAccessToken, err := s.tokenRepository.FindByID(oAuthRefreshToken.AccessTokenID)
	if err != nil || oAuthAccessToken.ClientID != client.ID.Hex() {
		return nil, utils.NewErrorUnauthorized("Client mismatch")
	}

	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes, []string{})
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &model.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		UserID:    oAuthAccessToken.UserID,
		GrantType: oAuthAccessToken.GrantType,
		ClientID:  oAuthAccessToken.ClientID,
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := s.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	originTokenId := utils.BinaryUUIDToString(oAuthAccessToken.ID)
	accessToken, err := utils.GenerateToken(oAuthAccessToken.UserID, scopes, oAuthAccessToken.GrantType, oAuthAccessToken.ClientID, tokenID, tokenExpireTime, &originTokenId)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		TokenType:   "Bearer",
		ExpiresIn:   tokenExpiresIn,
		AccessToken: accessToken,
	}, nil
}

func (s *oAuthUsecase) revokeRefreshToken(client *model.OAuthClient, request *middleware.OAuthRevoke) error {
	claims, err := utils.ValidateJWT(request.Token)
	if err != nil {
		return err
	}

	refreshTokenId := utils.StringToBinaryUUID(claims.Id)
	oAuthRefreshToken, err := s.refreshTokenRepository.FindByID(refreshTokenId)
	if err != nil || oAuthRefreshToken.Revoked == 1 {
		return utils.NewErrorUnauthorized("Refresh Token has been revoked or invalid")
	}

	oAuthAccessToken, err := s.tokenRepository.FindByID(oAuthRefreshToken.AccessTokenID)
	if err != nil || oAuthAccessToken.ClientID != client.ID.Hex() {
		return utils.NewErrorUnauthorized("Client mismatch")
	}

	if err := s.refreshTokenRepository.Update(refreshTokenId, &model.UpdateOAuthRefreshToken{Revoked: 1}); err != nil {
		return err
	}

	return nil
}

// Helper function to get token expiry time
func getTokenExpiryTime(envVar string, defaultValue int64) int64 {
	expiry, _ := strconv.ParseInt(config.GetEnv(envVar, strconv.FormatInt(defaultValue, 10)), 10, 64)
	return expiry
}

// Helper function to create and store access token
func (s *oAuthUsecase) createAndStoreAccessToken(tokenRequest *model.OAuthAccessToken) (string, error) {
	tokenCreationResult, err := s.tokenRepository.Create(tokenRequest)
	if err != nil {
		return "", err
	}

	binaryTokenID, ok := tokenCreationResult.InsertedID.(primitive.Binary)
	if !ok {
		return "", utils.NewErrorBadRequest("Failed to parse token ID")
	}

	return utils.BinaryUUIDToString(binaryTokenID), nil
}

// Helper function to generate and store refresh token
func (s *oAuthUsecase) generateAndStoreRefreshToken(accessTokenID string, userID string, clientID string) (string, error) {
	refreshTokenExpiresIn := getTokenExpiryTime("TOKEN_REFRESH_EXPIRE_TIME", 604800)

	refreshTokenRequest := &model.OAuthRefreshToken{
		ID:            utils.BinaryUUID(),
		AccessTokenID: utils.StringToBinaryUUID(accessTokenID),
		ExpiresIn:     refreshTokenExpiresIn,
	}

	refreshTokenCreationResult, err := s.refreshTokenRepository.Create(refreshTokenRequest)
	if err != nil {
		return "", err
	}

	binaryRefreshTokenID, ok := refreshTokenCreationResult.InsertedID.(primitive.Binary)
	if !ok {
		return "", utils.NewErrorBadRequest("Failed to parse refresh token ID")
	}

	refreshTokenID := utils.BinaryUUIDToString(binaryRefreshTokenID)
	return utils.GenerateRefreshToken(userID, clientID, refreshTokenID, tokenRefreshExpireTime)
}
