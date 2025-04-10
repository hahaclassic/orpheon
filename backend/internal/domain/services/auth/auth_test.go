package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()
	mockUserCreator := mocks.NewUserCreatorService(t)
	mockTokenService := mocks.NewTokenService(t)
	mockPasswordHasher := mocks.NewPasswordHasher(t)
	mockAuthRepo := mocks.NewAuthRepository(t)
	mockRefreshRepo := mocks.NewRefreshTokenRepository(t)

	service := NewAuthService(mockAuthRepo, mockRefreshRepo, mockUserCreator,
		mockPasswordHasher, mockTokenService)

	credentials := &entity.UserCredentials{
		Login:    "admin",
		Password: "password",
	}

	user := &entity.UserInfo{
		Name: credentials.Login,
	}

	claims := &entity.Claims{
		UserID:    uuid.New(),
		AccessLvl: entity.User,
	}

	// OK
	mockUserCreator.On("Create", ctx, user).Return("6d6fdb15-b669-4edf-bf94-7484af3cfc7b").Once()
	mockTokenService.On("GenerateAccessToken", user.ID).Return("access_token", nil).Once()
	mockTokenService.On("GenerateRefreshToken", user.ID).Return("refresh_token", nil).Once()

	token, err := service.RegisterUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token.AccessToken)
	assert.Equal(t, "refresh_token", token.RefreshToken)

	// ошибка от репозитория
	mockUserCreator.On("Create", ctx, user).Return(assert.AnError).Once()
	token, err = service.Register(ctx, user)
	assert.Error(t, err)
	assert.Nil(t, token)

	// ошибка при генерации токенов
	mockUserCreator.On("Create", ctx, user).Return(nil).Once()
	mockTokenService.On("GenerateAccessToken", user.ID).Return("", assert.AnError).Once()
	token, err = service.Register(ctx, user)
	assert.Error(t, err)
	assert.Nil(t, token)

	mockUserCreator.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	mockUserCreator := mocks.NewUserRepository(t)
	mockTokenService := mocks.NewTokenService(t)

	service := NewAuthService(mockUserCreator, mockTokenService)

	user := &entity.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "password",
	}

	// OK
	mockUserCreator.On("FindByEmail", ctx, user.Email).Return(user, nil).Once()
	mockTokenService.On("GenerateAccessToken", user.ID).Return("access_token", nil).Once()
	mockTokenService.On("GenerateRefreshToken", user.ID).Return("refresh_token", nil).Once()

	token, err := service.Login(ctx, user.Email, user.Password)
	assert.NoError(t, err)
	assert.Equal(t, "access_token", token.AccessToken)
	assert.Equal(t, "refresh_token", token.RefreshToken)

	// пользователь не найден
	mockUserCreator.On("FindByEmail", ctx, user.Email).Return(nil, assert.AnError).Once()
	token, err = service.Login(ctx, user.Email, user.Password)
	assert.Error(t, err)
	assert.Nil(t, token)

	// ошибка при генерации токенов
	mockUserCreator.On("FindByEmail", ctx, user.Email).Return(user, nil).Once()
	mockTokenService.On("GenerateAccessToken", user.ID).Return("", assert.AnError).Once()
	token, err = service.Login(ctx, user.Email, user.Password)
	assert.Error(t, err)
	assert.Nil(t, token)

	mockUserCreator.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()
	mockTokenService := mocks.NewTokenService(t)
	mockRefreshRepo := mocks.NewRefreshTokenRepository(t)

	service := NewAuthService(mockRefreshRepo, mockTokenService)

	tokenStr := "some_refresh_token"
	userID := uuid.New()

	// Мок для парсинга refresh токена
	claims := &entity.Claims{UserID: userID}
	mockTokenService.On("ParseAccessToken", tokenStr).Return(claims, nil).Once()

	// Мок для удаления refresh токена
	mockRefreshRepo.On("Delete", ctx, tokenStr).Return(nil).Once()

	// OK
	err := service.Logout(ctx, tokenStr)
	assert.NoError(t, err)

	// ошибка при парсинге токена
	mockTokenService.On("ParseAccessToken", tokenStr).Return(nil, assert.AnError).Once()
	err = service.Logout(ctx, tokenStr)
	assert.Error(t, err)

	// ошибка при удалении refresh токена
	mockTokenService.On("ParseAccessToken", tokenStr).Return(claims, nil).Once()
	mockRefreshRepo.On("Delete", ctx, tokenStr).Return(assert.AnError).Once()
	err = service.Logout(ctx, tokenStr)
	assert.Error(t, err)

	mockRefreshRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}
