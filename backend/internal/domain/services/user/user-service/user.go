package userservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	repo UserRepository
}

func New(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	return u.repo.CreateUser(ctx, user)
}

func (u *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	return u.repo.GetUser(ctx, userID)
}

func (u *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	return u.repo.UpdateUser(ctx, user)
}

func (u *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return u.repo.DeleteUser(ctx, userID)
}
