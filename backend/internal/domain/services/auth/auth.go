package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrHashingPassword = errors.New("hashing password error")
	ErrNewUserCreation = errors.New("new user creation error")
	ErrLogin           = errors.New("login error")
)

type TokenService interface {
	CreateAccessToken(claims *entity.Claims) (string, error)
	ParseAccessToken(tokenStr string) (*entity.Claims, error)
}

type PasswordHasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashed string, password string) error
}

type RefreshTokenRepository interface {
	Set(ctx context.Context, token string, claims *entity.Claims) error
	Get(ctx context.Context, token string) (*entity.Claims, error)
	Delete(ctx context.Context, token string) error
}

type AuthRepository interface {
	SaveCredentials(ctx context.Context, userID uuid.UUID, credentials *entity.UserCredentials) error
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
	GetClaimsByLogin(ctx context.Context, login string) (*entity.Claims, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

type UserCreatorService interface {
	CreateUser(ctx context.Context, info *entity.UserInfo) (uuid.UUID, error)
}

type AuthConfig struct {
	SecretKey  []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	HashCost   int
}

type AuthService struct {
	authRepo     AuthRepository
	refreshRepo  RefreshTokenRepository
	userCreator  UserCreatorService
	hasher       PasswordHasher
	tokenService TokenService
}

func NewAuthService(authRepo AuthRepository, refreshRepo RefreshTokenRepository,
	userCreator UserCreatorService, hasher PasswordHasher, token TokenService) *AuthService {
	return &AuthService{
		authRepo:     authRepo,
		refreshRepo:  refreshRepo,
		userCreator:  userCreator,
		hasher:       hasher,
		tokenService: token,
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error) {
	hashedPassword, err := a.hasher.GenerateFromPassword(credentials.Password)
	if err != nil {
		return nil, err
	}

	userID, err := a.userCreator.CreateUser(ctx, &entity.UserInfo{Name: credentials.Login})
	if err != nil {
		return nil, err
	}

	hashedCreds := &entity.UserCredentials{
		Login:    credentials.Login,
		Password: string(hashedPassword),
	}

	err = a.authRepo.SaveCredentials(ctx, userID, hashedCreds)
	if err != nil {
		return nil, err
	}

	return a.Login(ctx, credentials)
}

func (a *AuthService) Login(ctx context.Context, credentials *entity.UserCredentials) (*entity.AuthTokens, error) {
	hashedPassword, err := a.authRepo.GetPasswordByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if a.hasher.CompareHashAndPassword(hashedPassword, credentials.Password) != nil {
		return nil, errors.New("invalid credentials")
	}

	claims, err := a.authRepo.GetClaimsByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenService.CreateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refresh, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	refreshToken := refresh.String()

	err = a.refreshRepo.Set(ctx, refreshToken, claims)
	if err != nil {
		return nil, err
	}

	return &entity.AuthTokens{
		Access:  accessToken,
		Refresh: refreshToken}, nil
}

func (a *AuthService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return a.refreshRepo.Delete(ctx, refreshToken)
}

func (a *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*entity.AuthTokens, error) {
	claims, err := a.refreshRepo.Get(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenService.CreateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	newRefresh, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	newRefreshToken := newRefresh.String()

	err = a.refreshRepo.Set(ctx, newRefreshToken, claims)
	if err != nil {
		return nil, err
	}

	return &entity.AuthTokens{Access: accessToken, Refresh: newRefreshToken}, nil
}

func (a *AuthService) UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entity.UserPasswords) error {
	hashedPassword, err := a.authRepo.GetPasswordByID(ctx, userID)
	if err != nil {
		return err
	}

	if a.hasher.CompareHashAndPassword(hashedPassword, passwords.Old) != nil {
		return errors.New("invalid old password")
	}

	newHashed, err := a.hasher.GenerateFromPassword(passwords.New)
	if err != nil {
		return err
	}

	return a.authRepo.UpdatePassword(ctx, userID, string(newHashed))
}

func (a *AuthService) GetClaims(ctx context.Context, accessToken string) (*entity.Claims, error) {
	return a.tokenService.ParseAccessToken(accessToken)
}
