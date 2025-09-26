package auth

import (
	"context"
	"errors"
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

func TestAuthServiceRegisterUser_HasherError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewAuthRepository(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	creds := &entity.UserCredentials{Login: "test", Password: "pass"}
	hasher.On("GenerateFromPassword", creds.Password).Return("", errors.New("hash error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.RegisterUser(ctx, creds)

	assert.Error(t, err)
	assert.Nil(t, tokens)
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

func TestAuthService_Login_GetPasswordError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewAuthRepository(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	creds := &entity.UserCredentials{Login: "user", Password: "pass"}
	authRepo.On("GetPasswordByLogin", ctx, creds.Login).Return("", errors.New("db error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.Login(ctx, creds)

	assert.Error(t, err)
	assert.Nil(t, tokens)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewAuthRepository(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	creds := &entity.UserCredentials{Login: "user", Password: "wrong"}
	hashed := "hashed"
	authRepo.On("GetPasswordByLogin", ctx, creds.Login).Return(hashed, nil)
	hasher.On("CompareHashAndPassword", hashed, creds.Password).Return(errors.New("mismatch"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.Login(ctx, creds)

	assert.Error(t, err)
	assert.Nil(t, tokens)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Logout(t *testing.T) {
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)
	ctx := context.Background()

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	refreshToken := "refresh"

	refreshRepo.On("Delete", ctx, refreshToken).Return(nil)

	err := s.Logout(ctx, refreshToken)
	assert.NoError(t, err)
}

func TestAuthService_Logout_DeleteError(t *testing.T) {
	ctx := context.Background()
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	refreshToken := "refresh"
	refreshRepo.On("Delete", ctx, refreshToken).Return(errors.New("delete error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	err := s.Logout(ctx, refreshToken)

	assert.Error(t, err)
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

func TestAuthService_RefreshTokens_GetError(t *testing.T) {
	ctx := context.Background()
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	oldRefresh := "old"
	refreshRepo.On("Get", ctx, oldRefresh).Return(nil, errors.New("get error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	tokens, err := s.RefreshTokens(ctx, oldRefresh)

	assert.Error(t, err)
	assert.Nil(t, tokens)
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

func TestAuthService_UpdatePassword_GetPasswordError(t *testing.T) {
	ctx := context.Background()
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	userID := uuid.New()
	authRepo.On("GetPasswordByID", ctx, userID).Return("", errors.New("db error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	err := s.UpdatePassword(ctx, userID, &entity.UserPasswords{Old: "old", New: "new"})

	assert.Error(t, err)
}

func TestAuthService_UpdatePassword_InvalidOldPassword(t *testing.T) {
	ctx := context.Background()
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)
	tokenService := mocks.NewTokenService(t)

	userID := uuid.New()
	hashed := "hashed"
	authRepo.On("GetPasswordByID", ctx, userID).Return(hashed, nil)
	hasher.On("CompareHashAndPassword", hashed, "wrong").Return(errors.New("mismatch"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	err := s.UpdatePassword(ctx, userID, &entity.UserPasswords{Old: "wrong", New: "new"})

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
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

func TestAuthService_GetClaims_ParseError(t *testing.T) {
	ctx := context.Background()
	tokenService := mocks.NewTokenService(t)
	refreshRepo := mocks.NewRefreshTokenRepository(t)
	authRepo := mocks.NewAuthRepository(t)
	userCreator := mocks.NewUserCreatorService(t)
	hasher := mocks.NewPasswordHasher(t)

	token := "badtoken"
	tokenService.On("ParseAccessToken", token).Return(nil, errors.New("parse error"))

	s := NewAuthService(authRepo, refreshRepo, userCreator, hasher, tokenService)
	claims, err := s.GetClaims(ctx, token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
