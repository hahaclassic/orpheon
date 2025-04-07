package userservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.UserInfo) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entities.UserInfo, error)
	UpdateUser(ctx context.Context, user *entities.UserInfo) error
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

func (u *UserService) CreateUser(ctx context.Context, user *entities.UserInfo) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user.ID = id
	user.RegistrationDate = time.Now()

	return u.repo.CreateUser(ctx, user)
}

func (u *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*entities.UserInfo, error) {
	return u.repo.GetUser(ctx, userID)
}

func (u *UserService) UpdateUser(ctx context.Context, claims *entities.Claims, user *entities.UserInfo) error {
	if claims.UserID != user.ID {
		return errors.New("Forbidden")
	}

	return u.repo.UpdateUser(ctx, user)
}

func (u *UserService) DeleteUser(ctx context.Context, claims *entities.Claims, userID uuid.UUID) error {
	if claims.UserID != userID && claims.AccessLvl != entities.Admin {
		return errors.New("Forbidden")
	}

	return u.repo.DeleteUser(ctx, userID)
}
