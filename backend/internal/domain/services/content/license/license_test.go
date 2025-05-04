package license_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/license"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

func TestCreateLicense(t *testing.T) {
	ctx := context.Background()
	lic := &entity.License{ID: uuid.New()}
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Create", ctx, lic).Return(nil)
		svc := license.NewLicenseService(repo)
		err := svc.CreateLicense(ctx, adminClaims, lic)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.CreateLicense(ctx, userClaims, lic)
		assert.ErrorIs(t, err, license.ErrForbidden)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.CreateLicense(ctx, adminClaims, &entity.License{ID: uuid.Nil})
		assert.ErrorIs(t, err, license.ErrInvalidLicenseID)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Create", ctx, lic).Return(errors.New("db error"))
		svc := license.NewLicenseService(repo)
		err := svc.CreateLicense(ctx, adminClaims, lic)
		assert.ErrorIs(t, err, license.ErrCreateLicense)
	})
}

func TestGetLicense(t *testing.T) {
	ctx := context.Background()
	licenseID := uuid.New()
	lic := &entity.License{ID: licenseID}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Get", ctx, licenseID).Return(lic, nil)
		svc := license.NewLicenseService(repo)
		result, err := svc.GetLicense(ctx, licenseID)
		assert.NoError(t, err)
		assert.Equal(t, lic, result)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		result, err := svc.GetLicense(ctx, uuid.Nil)
		assert.ErrorIs(t, err, license.ErrInvalidLicenseID)
		assert.Nil(t, result)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Get", ctx, licenseID).Return(nil, errors.New("db error"))
		svc := license.NewLicenseService(repo)
		_, err := svc.GetLicense(ctx, licenseID)
		assert.ErrorIs(t, err, license.ErrGetLicense)
	})
}

func TestUpdateLicense(t *testing.T) {
	ctx := context.Background()
	lic := &entity.License{ID: uuid.New()}
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Update", ctx, lic).Return(nil)
		svc := license.NewLicenseService(repo)
		err := svc.UpdateLicense(ctx, adminClaims, lic)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.UpdateLicense(ctx, userClaims, lic)
		assert.ErrorIs(t, err, license.ErrForbidden)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.UpdateLicense(ctx, adminClaims, &entity.License{ID: uuid.Nil})
		assert.ErrorIs(t, err, license.ErrInvalidLicenseID)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Update", ctx, lic).Return(errors.New("db error"))
		svc := license.NewLicenseService(repo)
		err := svc.UpdateLicense(ctx, adminClaims, lic)
		assert.ErrorIs(t, err, license.ErrUpdateLicense)
	})
}

func TestDeleteLicense(t *testing.T) {
	ctx := context.Background()
	licenseID := uuid.New()
	adminClaims := &entity.Claims{AccessLvl: entity.Admin}
	userClaims := &entity.Claims{AccessLvl: entity.User}

	t.Run("success", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Delete", ctx, licenseID).Return(nil)
		svc := license.NewLicenseService(repo)
		err := svc.DeleteLicense(ctx, adminClaims, licenseID)
		assert.NoError(t, err)
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.DeleteLicense(ctx, userClaims, licenseID)
		assert.ErrorIs(t, err, license.ErrForbidden)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		svc := license.NewLicenseService(repo)
		err := svc.DeleteLicense(ctx, adminClaims, uuid.Nil)
		assert.ErrorIs(t, err, license.ErrInvalidLicenseID)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := mocks.NewLicenseRepository(t)
		repo.On("Delete", ctx, licenseID).Return(errors.New("db error"))
		svc := license.NewLicenseService(repo)
		err := svc.DeleteLicense(ctx, adminClaims, licenseID)
		assert.ErrorIs(t, err, license.ErrDeleteLicense)
	})
}
