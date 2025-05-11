package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usersvc "github.com/hahaclassic/orpheon/backend/internal/domain/services/user"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func(repo *mocks.UserRepository)
		prepare    func() *entity.UserInfo
		expectsErr bool
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
			},
			prepare: func() *entity.UserInfo {
				return &entity.UserInfo{}
			},
			expectsErr: false,
		},
		{
			name: "repo error",
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			prepare: func() *entity.UserInfo {
				return &entity.UserInfo{}
			},
			expectsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			service := usersvc.New(repo)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			_, err := service.CreateUser(context.Background(), tt.prepare())
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "CreateUser", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	service := usersvc.New(repo)

	userID := uuid.New()
	expected := &entity.UserInfo{ID: userID, Name: "testuser"}
	repo.On("GetUser", mock.Anything, userID).Return(expected, nil)

	result, err := service.GetUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertCalled(t, "GetUser", mock.Anything, userID)
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		claims     *entity.Claims
		user       *entity.UserInfo
		setupMock  func(repo *mocks.UserRepository)
		expectsErr bool
	}{
		{
			name:   "success",
			claims: &entity.Claims{UserID: uuid.New()},
			user:   &entity.UserInfo{},
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
			},
			expectsErr: false,
		},
		{
			name:       "forbidden",
			claims:     &entity.Claims{UserID: uuid.New()},
			user:       &entity.UserInfo{ID: uuid.New()},
			setupMock:  nil,
			expectsErr: true,
		},
		{
			name:   "repo error",
			claims: &entity.Claims{UserID: uuid.New()},
			user:   &entity.UserInfo{},
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("UpdateUser", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			service := usersvc.New(repo)
			if tt.setupMock != nil {
				tt.user.ID = tt.claims.UserID
				tt.setupMock(repo)
			}
			err := service.UpdateUser(context.Background(), tt.claims, tt.user)
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "UpdateUser", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	adminID := uuid.New()
	userID := uuid.New()
	tests := []struct {
		name       string
		claims     *entity.Claims
		userID     uuid.UUID
		setupMock  func(repo *mocks.UserRepository)
		expectsErr bool
	}{
		{
			name:       "user can delete self",
			claims:     &entity.Claims{UserID: userID},
			userID:     userID,
			setupMock:  func(repo *mocks.UserRepository) { repo.On("DeleteUser", mock.Anything, userID).Return(nil) },
			expectsErr: false,
		},
		{
			name:       "admin can delete any",
			claims:     &entity.Claims{UserID: adminID, AccessLvl: entity.Admin},
			userID:     userID,
			setupMock:  func(repo *mocks.UserRepository) { repo.On("DeleteUser", mock.Anything, userID).Return(nil) },
			expectsErr: false,
		},
		{
			name:       "forbidden",
			claims:     &entity.Claims{UserID: uuid.New()},
			userID:     userID,
			setupMock:  nil,
			expectsErr: true,
		},
		{
			name:   "repo error",
			claims: &entity.Claims{UserID: userID},
			userID: userID,
			setupMock: func(repo *mocks.UserRepository) {
				repo.On("DeleteUser", mock.Anything, userID).Return(errors.New("db error"))
			},
			expectsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			service := usersvc.New(repo)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			err := service.DeleteUser(context.Background(), tt.claims, tt.userID)
			if tt.expectsErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "DeleteUser", mock.Anything, tt.userID)
			}
		})
	}
}
