package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	authRepo := mocks.NewAuthRepository(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	authService := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)

	creds := &entity.UserCredentials{
		Login:    "testuser",
		Password: "securepass",
	}
	userID := uuid.New()
	hashedPassword := "hashed_securepass"
	claims := &entity.Claims{
		UserID:    userID,
		AccessLvl: entity.User,
	}
	accessToken := "access.token"
	refreshToken := "refresh.token"

	hasher.On("GenerateFromPassword", creds.Password).Return(hashedPassword, nil)
	userCreator.On("CreateUser", ctx, mock.MatchedBy(func(info *entity.UserInfo) bool {
		return info.Name == creds.Login
	})).Return(userID, nil)

	authRepo.On("SaveCredentials", ctx, userID, &entity.UserCredentials{
		Login:    creds.Login,
		Password: hashedPassword,
	}).Return(nil)

	authRepo.On("GetPasswordByLogin", ctx, creds.Login).Return(hashedPassword, nil)
	hasher.On("CompareHashAndPassword", hashedPassword, creds.Password).Return(nil)
	authRepo.On("GetClaimsByLogin", ctx, creds.Login).Return(claims, nil)
	tokenService.On("GenerateAccessToken", claims).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshToken, nil)
	refreshRepo.On("Set", ctx, refreshToken, claims).Return(nil)

	tokens, err := authService.RegisterUser(ctx, creds)

	assert.NoError(t, err)
	assert.Equal(t, accessToken, tokens.Access)
	assert.Equal(t, refreshToken, tokens.Refresh)

	hasher.AssertExpectations(t)
	userCreator.AssertExpectations(t)
	authRepo.AssertExpectations(t)
	refreshRepo.AssertExpectations(t)
	tokenService.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	authRepo := mocks.NewAuthRepository(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)
	ctx := context.Background()

	credentials := &entity.UserCredentials{
		Login:    "user",
		Password: "password",
	}

	hashedPassword := "hashed"
	claims := &entity.Claims{UserID: uuid.New(), AccessLvl: entity.User}
	accessToken := "access"
	refreshToken := "refresh"

	authRepo.On("GetPasswordByLogin", ctx, credentials.Login).Return(hashedPassword, nil)
	hasher.On("CompareHashAndPassword", hashedPassword, credentials.Password).Return(nil)
	authRepo.On("GetClaimsByLogin", ctx, credentials.Login).Return(claims, nil)
	tokenService.On("GenerateAccessToken", claims).Return(accessToken, nil)
	tokenService.On("GenerateRefreshToken").Return(refreshToken, nil)
	refreshRepo.On("Set", ctx, refreshToken, claims).Return(nil)

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.Login(ctx, credentials)

	assert.NoError(t, err)
	assert.Equal(t, &entity.AuthTokens{Access: accessToken, Refresh: refreshToken}, tokens)
}

func TestAuthService_Logout(t *testing.T) {
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)
	ctx := context.Background()

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	userID := uuid.New()
	refreshToken := "refresh"

	refreshRepo.On("Delete", ctx, refreshToken).Return(nil)

	err := s.Logout(ctx, userID, refreshToken)
	assert.NoError(t, err)
}

func TestAuthService_RefreshTokens(t *testing.T) {
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)
	ctx := context.Background()

	claims := &entity.Claims{UserID: uuid.New(), AccessLvl: entity.User}
	oldRefresh := "old_refresh"
	newRefresh := "new_refresh"
	access := "access"

	refreshRepo.On("Get", ctx, oldRefresh).Return(claims, nil)
	tokenService.On("GenerateAccessToken", claims).Return(access, nil)
	tokenService.On("GenerateRefreshToken").Return(newRefresh, nil)
	refreshRepo.On("Delete", ctx, oldRefresh).Return(nil)
	refreshRepo.On("Set", ctx, newRefresh, claims).Return(nil)

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.RefreshTokens(ctx, oldRefresh)

	assert.NoError(t, err)
	assert.Equal(t, &entity.AuthTokens{Access: access, Refresh: newRefresh}, tokens)
}

func TestAuthService_UpdatePassword(t *testing.T) {
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)
	ctx := context.Background()

	userID := uuid.New()
	passwords := &entity.UserPasswords{Old: "old", New: "new"}
	hashed := "hashed"
	newHashed := "new_hashed"

	authRepo.On("GetPasswordByID", ctx, userID).Return(hashed, nil)
	hasher.On("CompareHashAndPassword", hashed, passwords.Old).Return(nil)
	hasher.On("GenerateFromPassword", passwords.New).Return(newHashed, nil)
	authRepo.On("UpdatePassword", ctx, userID, newHashed).Return(nil)

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)

	err := s.UpdatePassword(ctx, userID, passwords)
	assert.NoError(t, err)
}

func TestAuthService_GetClaims(t *testing.T) {
	tokenService := mocks.NewTokenService(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	ctx := context.Background()

	claims := &entity.Claims{UserID: uuid.New(), AccessLvl: entity.User}
	token := "access"

	tokenService.On("ParseAccessToken", token).Return(claims, nil)

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)

	parsed, err := s.GetClaims(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, claims, parsed)
}
