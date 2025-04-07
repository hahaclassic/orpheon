package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
