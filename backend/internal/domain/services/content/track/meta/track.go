package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrGetTrackMeta    = errors.New("get track meta error")
	ErrCreateTrackMeta = errors.New("create track meta error")
	ErrUpdateTrackMeta = errors.New("update track meta error")
	ErrDeleteTrackMeta = errors.New("delete track meta error")

	ErrGenerateTrackID = errors.New("generate track id error")
	ErrForbidden       = errors.New("permission denied")
)

type TrackMetaRepository interface {
	GetByID(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error)
	Create(ctx context.Context, track *entity.TrackMeta) error
	Update(ctx context.Context, track *entity.TrackMeta) error
	Delete(ctx context.Context, trackID uuid.UUID) error
}

type TrackMetaService struct {
	repo TrackMetaRepository
}

func NewTrackMetaService(repo TrackMetaRepository) *TrackMetaService {
	return &TrackMetaService{repo: repo}
}

func (s *TrackMetaService) GetTrackMeta(ctx context.Context, trackID uuid.UUID) (*entity.TrackMeta, error) {
	meta, err := s.repo.GetByID(ctx, trackID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetTrackMeta, err)
	}

	return meta, nil
}

func (s *TrackMetaService) CreateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (uuid.UUID, error) {
	if claims.AccessLvl != entity.Admin {
		return uuid.Nil, ErrForbidden
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, errwrap.Wrap(ErrCreateTrackMeta, ErrGenerateTrackID)
	}

	track.ID = id

	if err := s.repo.Create(ctx, track); err != nil {
		return uuid.Nil, errwrap.Wrap(ErrGetTrackMeta, err)
	}

	return id, nil
}

func (s *TrackMetaService) UpdateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) error {
	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if err := s.repo.Update(ctx, track); err != nil {
		return errwrap.Wrap(ErrUpdateTrackMeta, err)
	}

	return nil
}

func (s *TrackMetaService) DeleteTrackMeta(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) error {
	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, trackID); err != nil {
		return errwrap.Wrap(ErrUpdateTrackMeta, err)
	}

	return nil
}
