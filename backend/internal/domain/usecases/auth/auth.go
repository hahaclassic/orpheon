package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type AuthService interface {
	RegisterUser(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error)
	Login(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*entity.AuthTokens, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entity.UserPasswords) error
	GetClaims(ctx context.Context, accessToken string) (*entity.Claims, error)
}
