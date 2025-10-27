package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/pkg/errwrap"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/entity"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/usecase"
)

var (
	ErrGenerateID = errors.New("id generation error")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
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

func (u *UserService) CreateUser(ctx context.Context, user *entity.User) (_ uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateUser, err)
	}()

	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, ErrGenerateID
	}

	user.ID = id
	user.RegistrationDate = time.Now()
	user.AccessLvl = entity.UserLvl

	if err = u.repo.CreateUser(ctx, user); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (u *UserService) GetUser(ctx context.Context, userID uuid.UUID) (_ *entity.User, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUser, err)
	}()

	return u.repo.GetUser(ctx, userID)
}

func (u *UserService) UpdateUser(ctx context.Context, claims *entity.Claims, user *entity.User) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateUser, err)
	}()

	if claims == nil || claims.UserID != user.ID {
		return commonerr.ErrForbidden
	}

	return u.repo.UpdateUser(ctx, user)
}

func (u *UserService) DeleteUser(ctx context.Context, claims *entity.Claims, userID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteUser, err)
	}()

	if claims == nil || claims.UserID != userID && claims.AccessLvl != entity.AdminLvl {
		return commonerr.ErrForbidden
	}

	return u.repo.DeleteUser(ctx, userID)
}
