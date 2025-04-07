package authservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type refreshTokenRepository interface {
	SetRefreshToken(ctx context.Context, refreshToken string, claims *entities.Claims) error
	GetClaims(ctx context.Context, refreshToken string) (*entities.Claims, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type AuthRepository interface {
	SavePassword(ctx context.Context, userID uuid.UUID, passwords *entities.UserPasswords) error
	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	GetClaims(ctx context.Context, credentials *entities.UserCredentials) (*entities.Claims, error)
}

type UserCreaterService interface {
	CreateUser(ctx context.Context, user *entities.UserInfo) error
}

type AuthService struct {
	user        UserCreaterService
	refreshRepo refreshTokenRepository
	auth        AuthRepository
}

func (a *AuthService) RegisterUser(ctx context.Context, credentials *entities.UserCredentials) (*entities.AuthTokens, error) {

}

func (a *AuthService) Login(ctx context.Context, credentials *entities.UserCredentials) (*entities.AuthTokens, error) {

}

func (a *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {

}

func (a *AuthService) UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entities.UserPasswords) error {

}

func (a *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*entities.AuthTokens, error) {

}

func (a *AuthService) GetClaims(ctx context.Context, accessToken string) (*entities.Claims, error) {

}
