package license

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type LicenseRepository interface {
	Create(ctx context.Context, license *entity.License) error
	Get(ctx context.Context, licenseID uuid.UUID) (*entity.License, error)
	Update(ctx context.Context, license *entity.License) error
	Delete(ctx context.Context, licenseID uuid.UUID) error
}

var (
	ErrCreateLicense = errors.New("failed to create license")
	ErrGetLicense    = errors.New("failed to get license")
	ErrUpdateLicense = errors.New("failed to update license")
	ErrDeleteLicense = errors.New("failed to delete license")

	ErrForbidden        = errors.New("permission denied error")
	ErrInvalidLicenseID = errors.New("invalid license ID")
)

type LicenseService struct {
	repo LicenseRepository
}

func NewLicenseService(repo LicenseRepository) *LicenseService {
	return &LicenseService{repo: repo}
}

func (s *LicenseService) CreateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrCreateLicense, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if license.ID == uuid.Nil {
		return ErrInvalidLicenseID
	}

	return s.repo.Create(ctx, license)
}

func (s *LicenseService) GetLicense(ctx context.Context, licenseID uuid.UUID) (_ *entity.License, err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrGetLicense, err)
	}()

	if licenseID == uuid.Nil {
		return nil, ErrInvalidLicenseID
	}

	return s.repo.Get(ctx, licenseID)
}

func (s *LicenseService) UpdateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrUpdateLicense, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if license.ID == uuid.Nil {
		return ErrInvalidLicenseID
	}

	return s.repo.Update(ctx, license)
}

func (s *LicenseService) DeleteLicense(ctx context.Context, claims *entity.Claims, licenseID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(ErrDeleteLicense, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if licenseID == uuid.Nil {
		return ErrInvalidLicenseID
	}

	return s.repo.Delete(ctx, licenseID)
}
