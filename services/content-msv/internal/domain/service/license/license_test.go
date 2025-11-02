package license

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/pkg/commonerr"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/license"
	"github.com/hahaclassic/orpheon/services/content-msv/mocks"
)

type LicenseObjectMother struct{}

func (LicenseObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.AdminLvl}
}

func (LicenseObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.UserLvl}
}

func (LicenseObjectMother) ValidLicense() *entity.License {
	return &entity.License{ID: uuid.New(), Title: "Pro"}
}

type LicenseServiceSuite struct {
	suite.Suite

	ctx    context.Context
	repo   *mocks.LicenseRepository
	svc    *LicenseService
	mother LicenseObjectMother
}

func TestLicenseServiceSuite(t *testing.T) {
	suite.Run(t, new(LicenseServiceSuite))
}

func (s *LicenseServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewLicenseRepository(s.T())
	s.svc = NewLicenseService(s.repo)
	s.mother = LicenseObjectMother{}
}

func (s *LicenseServiceSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

// --- CreateLicense ---

func (s *LicenseServiceSuite) TestCreateLicense() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	valid := s.mother.ValidLicense()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("Create", s.ctx, valid).Return(nil)
		err := s.svc.CreateLicense(s.ctx, admin, valid)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		err := s.svc.CreateLicense(s.ctx, user, valid)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("Create", s.ctx, valid).Return(errors.New("db error"))
		err := s.svc.CreateLicense(s.ctx, admin, valid)
		s.ErrorIs(err, usecase.ErrCreateLicense)
	})
}

// --- GetLicenseByID ---

func (s *LicenseServiceSuite) TestGetLicenseByID() {
	valid := s.mother.ValidLicense()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("GetByID", s.ctx, valid.ID).Return(valid, nil)
		res, err := s.svc.GetLicenseByID(s.ctx, valid.ID)
		s.NoError(err)
		s.Equal(valid, res)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		res, err := s.svc.GetLicenseByID(s.ctx, uuid.Nil)
		s.ErrorIs(err, ErrInvalidLicenseID)
		s.Nil(res)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("GetByID", s.ctx, valid.ID).Return(nil, errors.New("db error"))
		_, err := s.svc.GetLicenseByID(s.ctx, valid.ID)
		s.ErrorIs(err, usecase.ErrGetLicense)
	})
}

// --- UpdateLicense ---

func (s *LicenseServiceSuite) TestUpdateLicense() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	valid := s.mother.ValidLicense()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("Update", s.ctx, valid).Return(nil)
		err := s.svc.UpdateLicense(s.ctx, admin, valid)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		err := s.svc.UpdateLicense(s.ctx, user, valid)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		err := s.svc.UpdateLicense(s.ctx, admin, &entity.License{ID: uuid.Nil})
		s.ErrorIs(err, ErrInvalidLicenseID)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("Update", s.ctx, valid).Return(errors.New("db error"))
		err := s.svc.UpdateLicense(s.ctx, admin, valid)
		s.ErrorIs(err, usecase.ErrUpdateLicense)
	})
}

// --- DeleteLicense ---

func (s *LicenseServiceSuite) TestDeleteLicense() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	validID := uuid.New()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("Delete", s.ctx, validID).Return(nil)
		err := s.svc.DeleteLicense(s.ctx, admin, validID)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		err := s.svc.DeleteLicense(s.ctx, user, validID)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		err := s.svc.DeleteLicense(s.ctx, admin, uuid.Nil)
		s.ErrorIs(err, ErrInvalidLicenseID)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("Delete", s.ctx, validID).Return(errors.New("db error"))
		err := s.svc.DeleteLicense(s.ctx, admin, validID)
		s.ErrorIs(err, usecase.ErrDeleteLicense)
	})
}

// --- GetAllLicenses ---

func (s *LicenseServiceSuite) TestGetAllLicenses() {
	valid := []*entity.License{s.mother.ValidLicense()}

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("GetAll", s.ctx).Return(valid, nil)
		res, err := s.svc.GetAllLicenses(s.ctx)
		s.NoError(err)
		s.Equal(valid, res)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("GetAll", s.ctx).Return(nil, errors.New("db error"))
		res, err := s.svc.GetAllLicenses(s.ctx)
		s.ErrorIs(err, usecase.ErrGetAllLicenses)
		s.Nil(res)
	})
}
