package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
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

func (s *TrackMetaService) GetTrackMeta(ctx context.Context, trackID uuid.UUID) (_ *entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetTrackMeta, err)
	}()

	return s.repo.GetByID(ctx, trackID)
}

func (s *TrackMetaService) CreateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (id uuid.UUID, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateTrackMeta, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return uuid.Nil, ErrForbidden
	}

	track.ID, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, ErrGenerateTrackID
	}

	if err = s.repo.Create(ctx, track); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *TrackMetaService) UpdateTrackMeta(ctx context.Context, claims *entity.Claims, track *entity.TrackMeta) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateTrackMeta, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return s.repo.Update(ctx, track)
}

func (s *TrackMetaService) DeleteTrackMeta(ctx context.Context, claims *entity.Claims, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteTrackMeta, err)
	}()

	if claims.AccessLvl != entity.Admin {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, trackID)
}
