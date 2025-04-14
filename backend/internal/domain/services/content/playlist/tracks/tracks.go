package tracks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrForbidden = errors.New("permission denied")
)

type playlistTracksRepository interface {
	AddTrackToPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteTrackFromPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteAllTracksFromPlaylist(ctx context.Context, playlistID uuid.UUID) error
	GetAllPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]*entity.TrackMeta, error)
}

type playlistPolicyService interface {
	CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
}

type PlaylistTrackService struct {
	repo   playlistTracksRepository
	policy playlistPolicyService
}

func NewPlaylistTrackService(repo playlistTracksRepository, policy playlistPolicyService) *PlaylistTrackService {
	return &PlaylistTrackService{
		repo:   repo,
		policy: policy,
	}
}

func (s *PlaylistTrackService) AddTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrAddTrack, err)
		}
	}()

	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("%w: user cannot edit playlist", ErrForbidden)
	}

	err = s.repo.AddTrackToPlaylist(ctx, playlistID, trackID)
	return err
}

func (s *PlaylistTrackService) GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrGetAllTracks, err)
		}
	}()

	canView, err := s.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, fmt.Errorf("%w: user cannot view playlist", ErrForbidden)
	}

	tracks, err = s.repo.GetAllPlaylistTracks(ctx, playlistID)
	return tracks, err
}

func (s *PlaylistTrackService) DeleteTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteTrack, err)
		}
	}()

	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("%w: user cannot edit playlist", ErrForbidden)
	}

	err = s.repo.DeleteTrackFromPlaylist(ctx, playlistID, trackID)
	return err
}

func (s *PlaylistTrackService) DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			err = errwrap.Wrap(usecase.ErrDeleteAllTracks, err)
		}
	}()

	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("%w: user cannot edit playlist", ErrForbidden)
	}

	err = s.repo.DeleteAllTracksFromPlaylist(ctx, playlistID)
	return err
}
