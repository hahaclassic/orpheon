package meta_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

// --- Object Mother ---

type ArtistMetaObjectMother struct{}

func (ArtistMetaObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (ArtistMetaObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

func (ArtistMetaObjectMother) DefaultArtistMeta() *entity.ArtistMeta {
	return &entity.ArtistMeta{
		ID:          uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		Name:        "Test Artist",
		Description: "Test Bio",
	}
}

// --- Suite ---

type ArtistMetaServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.ArtistMetaRepository
	service   *meta.ArtistMetaService
	objMother *ArtistMetaObjectMother
}

func TestArtistMetaServiceSuite(t *testing.T) {
	suite.Run(t, new(ArtistMetaServiceSuite))
}

func (s *ArtistMetaServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewArtistMetaRepository(s.T())
	s.service = meta.New(s.repo)
	s.objMother = &ArtistMetaObjectMother{}
}

// --- GetArtistMeta ---

func (s *ArtistMetaServiceSuite) TestGetArtistMeta_Success() {
	artist := s.objMother.DefaultArtistMeta()
	s.repo.On("GetByID", s.ctx, artist.ID).Return(artist, nil)

	got, err := s.service.GetArtistMeta(s.ctx, artist.ID)

	s.NoError(err)
	s.Equal(artist, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistMetaServiceSuite) TestGetArtistMeta_RepoError() {
	artist := s.objMother.DefaultArtistMeta()
	s.repo.On("GetByID", s.ctx, artist.ID).Return(nil, errors.New("not found"))

	got, err := s.service.GetArtistMeta(s.ctx, artist.ID)

	s.Nil(got)
	s.ErrorIs(err, usecase.ErrGetArtistMeta)
	s.repo.AssertExpectations(s.T())
}

// --- GetAllArtistMeta ---

func (s *ArtistMetaServiceSuite) TestGetAllArtistMeta_Success() {
	expected := []*entity.ArtistMeta{s.objMother.DefaultArtistMeta()}
	s.repo.On("GetAll", s.ctx).Return(expected, nil)

	got, err := s.service.GetAllArtistMeta(s.ctx)

	s.NoError(err)
	s.Equal(expected, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistMetaServiceSuite) TestGetAllArtistMeta_RepoError() {
	s.repo.On("GetAll", s.ctx).Return(nil, errors.New("db error"))

	got, err := s.service.GetAllArtistMeta(s.ctx)

	s.Nil(got)
	s.ErrorIs(err, usecase.ErrGetAllArtistMeta)
	s.repo.AssertExpectations(s.T())
}

// --- CreateArtistMeta ---

func (s *ArtistMetaServiceSuite) TestCreateArtistMeta_Success() {
	admin := s.objMother.AdminClaims()
	artist := &entity.ArtistMeta{Name: "New Artist"}

	s.repo.On("Create", s.ctx, mock.MatchedBy(func(a *entity.ArtistMeta) bool {
		return a.ID != uuid.Nil && a.Name == "New Artist"
	})).Return(nil)

	err := s.service.CreateArtistMeta(s.ctx, admin, artist)

	s.NoError(err)
	s.NotEqual(uuid.Nil, artist.ID)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistMetaServiceSuite) TestCreateArtistMeta_Forbidden() {
	user := s.objMother.UserClaims()
	artist := s.objMother.DefaultArtistMeta()

	err := s.service.CreateArtistMeta(s.ctx, user, artist)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *ArtistMetaServiceSuite) TestCreateArtistMeta_RepoError() {
	admin := s.objMother.AdminClaims()
	artist := s.objMother.DefaultArtistMeta()

	s.repo.On("Create", s.ctx, mock.Anything).Return(errors.New("db fail"))

	err := s.service.CreateArtistMeta(s.ctx, admin, artist)

	s.ErrorIs(err, usecase.ErrCreateArtistMeta)
	s.repo.AssertExpectations(s.T())
}

// --- UpdateArtistMeta ---

func (s *ArtistMetaServiceSuite) TestUpdateArtistMeta_Success() {
	admin := s.objMother.AdminClaims()
	artist := s.objMother.DefaultArtistMeta()

	s.repo.On("Update", s.ctx, artist).Return(nil)

	err := s.service.UpdateArtistMeta(s.ctx, admin, artist)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistMetaServiceSuite) TestUpdateArtistMeta_Forbidden() {
	user := s.objMother.UserClaims()
	artist := s.objMother.DefaultArtistMeta()

	err := s.service.UpdateArtistMeta(s.ctx, user, artist)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ArtistMetaServiceSuite) TestUpdateArtistMeta_RepoError() {
	admin := s.objMother.AdminClaims()
	artist := s.objMother.DefaultArtistMeta()

	s.repo.On("Update", s.ctx, artist).Return(errors.New("update failed"))

	err := s.service.UpdateArtistMeta(s.ctx, admin, artist)

	s.ErrorIs(err, usecase.ErrUpdateArtistMeta)
	s.repo.AssertExpectations(s.T())
}

// --- DeleteArtistMeta ---

func (s *ArtistMetaServiceSuite) TestDeleteArtistMeta_Success() {
	admin := s.objMother.AdminClaims()
	artist := s.objMother.DefaultArtistMeta()

	s.repo.On("Delete", s.ctx, artist.ID).Return(nil)

	err := s.service.DeleteArtistMeta(s.ctx, admin, artist.ID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ArtistMetaServiceSuite) TestDeleteArtistMeta_Forbidden() {
	user := s.objMother.UserClaims()
	artist := s.objMother.DefaultArtistMeta()

	err := s.service.DeleteArtistMeta(s.ctx, user, artist.ID)

	s.ErrorIs(err, commonerr.ErrForbidden)
	s.repo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *ArtistMetaServiceSuite) TestDeleteArtistMeta_RepoError() {
	admin := s.objMother.AdminClaims()
	artist := s.objMother.DefaultArtistMeta()

	s.repo.On("Delete", s.ctx, artist.ID).Return(errors.New("delete error"))

	err := s.service.DeleteArtistMeta(s.ctx, admin, artist.ID)

	s.ErrorIs(err, usecase.ErrDeleteArtistMeta)
	s.repo.AssertExpectations(s.T())
}
