package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type AuthService interface {
	RegisterUser(ctx context.Context, credentials *entities.UserCredentials) (*entities.AuthTokens, error)
	Login(ctx context.Context, credentials *entities.UserCredentials) (*entities.AuthTokens, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entities.UserPasswords) error
	RefreshTokens(ctx context.Context, refreshToken string) (*entities.AuthTokens, error)
	GetClaims(ctx context.Context, accessToken string) (*entities.Claims, error)
}
