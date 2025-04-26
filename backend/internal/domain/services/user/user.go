package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/user"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGenerateID = errors.New("id generation error")
	ErrForbidden  = errors.New("permission denied error")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.UserInfo) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.UserInfo, error)
	UpdateUser(ctx context.Context, user *entity.UserInfo) error
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

func (u *UserService) CreateUser(ctx context.Context, user *entity.UserInfo) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrCreateUser, err)
		}
	}()

	id, err := uuid.NewRandom()
	if err != nil {
		return ErrGenerateID
	}

	user.ID = id
	user.RegistrationDate = time.Now()

	return u.repo.CreateUser(ctx, user)
}

func (u *UserService) GetUser(ctx context.Context, userID uuid.UUID) (_ *entity.UserInfo, err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrGetUser, err)
		}
	}()

	return u.repo.GetUser(ctx, userID)
}

func (u *UserService) UpdateUser(ctx context.Context, claims *entity.Claims, user *entity.UserInfo) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrUpdateUser, err)
		}
	}()

	if claims.UserID != user.ID {
		return ErrForbidden
	}

	return u.repo.UpdateUser(ctx, user)
}

func (u *UserService) DeleteUser(ctx context.Context, claims *entity.Claims, userID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteUser, err)
		}
	}()

	if claims.UserID != userID && claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return u.repo.DeleteUser(ctx, userID)
}
