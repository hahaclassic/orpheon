package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entities.UserInfo) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.UserInfo, error)
	UpdateUser(ctx context.Context, claims *entities.Claims, user *entities.UserInfo) error
	DeleteUser(ctx context.Context, claims *entities.Claims, userID uuid.UUID) error
}
