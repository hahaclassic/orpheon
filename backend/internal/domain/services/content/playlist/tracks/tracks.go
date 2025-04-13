package tracks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrAddTrack        = errors.New("track addition error")
	ErrGetAllTracks    = errors.New("get all tracks error")
	ErrDeleteTrack     = errors.New("delete track error")
	ErrDeleteAllTracks = errors.New("delete all tracks error")

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

func (s *PlaylistTrackService) AddTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return errwrap.Wrap(ErrAddTrack, err)
	}
	if !canEdit {
		return fmt.Errorf("%w: %w: user cannot edit playlist", ErrAddTrack, ErrForbidden)
	}

	err = s.repo.AddTrackToPlaylist(ctx, playlistID, trackID)

	return errwrap.WrapIfErr(ErrAddTrack, err)
}

func (s *PlaylistTrackService) GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) ([]*entity.TrackMeta, error) {
	canView, err := s.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllTracks, err)
	}
	if !canView {
		return nil, fmt.Errorf("%w: %w: user cannot view playlist", ErrGetAllTracks, ErrForbidden)
	}

	tracks, err := s.repo.GetAllPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllTracks, err)
	}

	return tracks, nil
}

func (s *PlaylistTrackService) DeleteTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return errwrap.Wrap(ErrDeleteTrack, err)
	}
	if !canEdit {
		return fmt.Errorf("%w: %w: user cannot edit playlist", ErrDeleteTrack, ErrForbidden)
	}

	err = s.repo.DeleteTrackFromPlaylist(ctx, playlistID, trackID)

	return errwrap.WrapIfErr(ErrDeleteTrack, err)
}

func (s *PlaylistTrackService) DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return errwrap.Wrap(ErrDeleteAllTracks, err)
	}
	if !canEdit {
		return fmt.Errorf("%w: %w: user cannot edit playlist", ErrDeleteAllTracks, ErrForbidden)
	}

	err = s.repo.DeleteAllTracksFromPlaylist(ctx, playlistID)

	return errwrap.WrapIfErr(ErrDeleteAllTracks, err)
}
