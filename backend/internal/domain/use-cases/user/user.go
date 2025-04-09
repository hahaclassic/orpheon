package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.UserInfo) (uuid.UUID, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error)
	UpdateUser(ctx context.Context, claims *entity.Claims, user *entity.UserInfo) error
	DeleteUser(ctx context.Context, claims *entity.Claims, userID uuid.UUID) error
}
