package usecase

import (
	"errors"
	"go-account/config"

	"go-account/internal/domain"
	"go-account/internal/repository"

	"go-account/pkg/middleware"
	"go-account/pkg/oauth"
	"go-account/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type OAuthUsecase interface {
	Token(clientID string, clientSecret string, request middleware.OAuthToken) (*TokenDetails, error)
	Revoke(clientID string, clientSecret string, request middleware.OAuthRevoke) error
	CreateOAuthClient(client *domain.OAuthClient) (*domain.OAuthClient, error)
	GetOAuthClients(ctx *gin.Context) ([]domain.OAuthClient, int64, error)
	UpdateOAuthClient(id string, client *domain.UpdateOAuthClient) error
}

type TokenDetails struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type oAuthUsecase struct {
	jwtIssue               *oauth.JWTIssue
	jwtVerify              *oauth.JWTVerify
	userRepository         repository.UserRepository
	clientRepository       repository.OAuthClientRepository
	tokenRepository        repository.OAuthAccessTokenRepository
	refreshTokenRepository repository.OAuthRefreshTokenRepository
	validator              *validator.Validate
}

func (u *oAuthUsecase) Token(clientID string, clientSecret string, request middleware.OAuthToken) (*TokenDetails, error) {
	ocid, _ := primitive.ObjectIDFromHex(clientID)
	client, err := u.clientRepository.FindByID(ocid)
	if err != nil || clientSecret != client.Secret || !utils.InSlice(client.GrantTypes, request.GrantType) {
		return nil, utils.NewErrorUnauthorized("Client authentication failed")
	}
	if client.Revoked == 1 {
		return nil, utils.NewErrorUnauthorized("Client authentication has been revoked")
	}

	switch request.GrantType {
	case "password":
		return u.password(&client, &request)
	case "client_credentials":
		return u.clientCredentials(&client, &request)
	case "refresh_token":
		return u.refreshToken(&client, &request)
	default:
		return nil, utils.NewErrorUnauthorized("The authorization grant type is not supported by the authorization server.")
	}
}

func (u *oAuthUsecase) Revoke(clientID string, clientSecret string, request middleware.OAuthRevoke) error {
	ocid, _ := primitive.ObjectIDFromHex(clientID)
	client, err := u.clientRepository.FindByID(ocid)
	if err != nil || clientSecret != client.Secret {
		return utils.NewErrorUnauthorized("Client authentication failed")
	}
	if client.Revoked == 1 {
		return utils.NewErrorUnauthorized("Client authentication has been revoked")
	}
	switch request.TokenTypeHint {
	case "refresh_token":
		return u.revokeRefreshToken(&client, &request)
	default:
		return utils.NewErrorUnauthorized("The token_type_hint is not supported by the authorization server.")
	}
}

func (u *oAuthUsecase) CreateOAuthClient(client *domain.OAuthClient) (*domain.OAuthClient, error) {
	if err := u.validator.Struct(client); err != nil {
		return nil, err
	}

	// Generate client secret
	secret, err := utils.GenerateSecretBase64(32)
	if err != nil {
		return nil, err
	}
	client.Secret = secret

	// Create the client in the repository
	result, err := u.clientRepository.Create(client)
	if err != nil {
		return nil, err
	}

	// Extract ObjectID and find client by ID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert inserted ID to ObjectID")
	}

	_client, err := u.clientRepository.FindByID(oid)
	if err != nil {
		return nil, err
	}

	return &_client, err
}

func (u *oAuthUsecase) GetOAuthClients(ctx *gin.Context) ([]domain.OAuthClient, int64, error) {
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
	return u.clientRepository.Lists(filter, skip, size)
}

func (u *oAuthUsecase) UpdateOAuthClient(id string, client *domain.UpdateOAuthClient) error {
	if err := u.validator.Struct(client); err != nil {
		return err
	}
	ocid, _ := primitive.ObjectIDFromHex(id)
	return u.clientRepository.Update(ocid, client)
}

func NewOauthUsecase(
	jwtIssue *oauth.JWTIssue,
	jwtVerify *oauth.JWTVerify,
	userRepository repository.UserRepository,
	clientRepository repository.OAuthClientRepository,
	tokenRepository repository.OAuthAccessTokenRepository,
	refreshTokenRepository repository.OAuthRefreshTokenRepository) OAuthUsecase {

	validate := validator.New()
	validate.RegisterValidation("grant_types", utils.ValidateGrantTypes)

	return &oAuthUsecase{jwtIssue, jwtVerify, userRepository, clientRepository, tokenRepository, refreshTokenRepository, validate}
}

func (u *oAuthUsecase) clientCredentials(client *domain.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes)
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &domain.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := u.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwtIssue.GenerateClientToken(scopes, request.GrantType, client.ID.Hex(), tokenID)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		TokenType:   "Bearer",
		ExpiresIn:   tokenExpiresIn,
		AccessToken: accessToken,
	}, nil
}

func (u *oAuthUsecase) password(client *domain.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	user, err := u.userRepository.FindUserByUsername(request.Username)
	if err != nil {
		return nil, utils.NewErrorBadRequest("Your username is not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, utils.NewErrorBadRequest("Your password is incorrect")
	}

	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes, []string{})
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &domain.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		UserID:    user.ID.Hex(),
		GrantType: request.GrantType,
		ClientID:  client.ID.Hex(),
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := u.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwtIssue.GenerateToken(user.ID.Hex(), scopes, request.GrantType, client.ID.Hex(), tokenID, nil)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.generateAndStoreRefreshToken(tokenID, user.ID.Hex(), client.ID.Hex())
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

func (u *oAuthUsecase) refreshToken(client *domain.OAuthClient, request *middleware.OAuthToken) (*TokenDetails, error) {
	claims, err := u.jwtVerify.ValidateToken(request.RefreshToken)
	if err != nil {
		return nil, err
	}

	oAuthRefreshToken, err := u.refreshTokenRepository.FindByID(utils.StringToBinaryUUID(claims.Id))
	if err != nil || oAuthRefreshToken.Revoked == 1 {
		return nil, utils.NewErrorUnauthorized("Refresh Token has been revoked or invalid")
	}

	oAuthAccessToken, err := u.tokenRepository.FindByID(oAuthRefreshToken.AccessTokenID)
	if err != nil || oAuthAccessToken.ClientID != client.ID.Hex() {
		return nil, utils.NewErrorUnauthorized("Client mismatch")
	}

	scopes := utils.MergeSliceAndRemoveDuplicates(client.Scopes, []string{})
	tokenExpiresIn := getTokenExpiryTime("TOKEN_EXPIRE_TIME", 86400)

	tokenRequest := &domain.OAuthAccessToken{
		ID:        utils.BinaryUUID(),
		UserID:    oAuthAccessToken.UserID,
		GrantType: oAuthAccessToken.GrantType,
		ClientID:  oAuthAccessToken.ClientID,
		Scopes:    scopes,
		ExpiresIn: tokenExpiresIn,
	}

	tokenID, err := u.createAndStoreAccessToken(tokenRequest)
	if err != nil {
		return nil, err
	}

	originTokenId := utils.BinaryUUIDToString(oAuthAccessToken.ID)
	accessToken, err := u.jwtIssue.GenerateToken(oAuthAccessToken.UserID, scopes, oAuthAccessToken.GrantType, oAuthAccessToken.ClientID, tokenID, &originTokenId)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		TokenType:   "Bearer",
		ExpiresIn:   tokenExpiresIn,
		AccessToken: accessToken,
	}, nil
}

func (u *oAuthUsecase) revokeRefreshToken(client *domain.OAuthClient, request *middleware.OAuthRevoke) error {
	claims, err := u.jwtVerify.ValidateToken(request.Token)
	if err != nil {
		return err
	}

	refreshTokenId := utils.StringToBinaryUUID(claims.Id)
	oAuthRefreshToken, err := u.refreshTokenRepository.FindByID(refreshTokenId)
	if err != nil || oAuthRefreshToken.Revoked == 1 {
		return utils.NewErrorUnauthorized("Refresh Token has been revoked or invalid")
	}

	oAuthAccessToken, err := u.tokenRepository.FindByID(oAuthRefreshToken.AccessTokenID)
	if err != nil || oAuthAccessToken.ClientID != client.ID.Hex() {
		return utils.NewErrorUnauthorized("Client mismatch")
	}
	revoked := &domain.UpdateOAuthRefreshToken{Revoked: 1}
	if err := u.refreshTokenRepository.Update(refreshTokenId, revoked); err != nil {
		return err
	}

	return nil
}

// Helper function to create and store access token
func (u *oAuthUsecase) createAndStoreAccessToken(tokenRequest *domain.OAuthAccessToken) (string, error) {
	tokenCreationResult, err := u.tokenRepository.Create(tokenRequest)
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
func (u *oAuthUsecase) generateAndStoreRefreshToken(accessTokenID string, userID string, clientID string) (string, error) {
	refreshTokenExpiresIn := getTokenExpiryTime("TOKEN_REFRESH_EXPIRE_TIME", 604800)

	refreshTokenRequest := &domain.OAuthRefreshToken{
		ID:            utils.BinaryUUID(),
		AccessTokenID: utils.StringToBinaryUUID(accessTokenID),
		ExpiresIn:     refreshTokenExpiresIn,
	}

	refreshTokenCreationResult, err := u.refreshTokenRepository.Create(refreshTokenRequest)
	if err != nil {
		return "", err
	}

	binaryRefreshTokenID, ok := refreshTokenCreationResult.InsertedID.(primitive.Binary)
	if !ok {
		return "", utils.NewErrorBadRequest("Failed to parse refresh token ID")
	}

	refreshTokenID := utils.BinaryUUIDToString(binaryRefreshTokenID)
	return u.jwtIssue.GenerateRefreshToken(userID, clientID, refreshTokenID)
}

// Helper function to get token expiry time
func getTokenExpiryTime(envVar string, defaultValue int64) int64 {
	expiry, _ := strconv.ParseInt(config.GetEnv(envVar, strconv.FormatInt(defaultValue, 10)), 10, 64)
	return expiry
}
