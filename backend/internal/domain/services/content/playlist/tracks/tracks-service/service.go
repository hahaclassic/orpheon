package tracksservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

type playlistTracksRepository interface {
	AddTrackToPlaylist(playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteTrackFromPlaylist(playlistID uuid.UUID, trackID uuid.UUID) error
	GetTracksInPlaylist(playlistID uuid.UUID) ([]entities.TrackMeta, error)
}

type playlistPolicyService interface {
	CanDelete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error)
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

func (s *PlaylistTrackService) AddTrack(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID, trackID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return errors.New("permission denied: user cannot edit playlist")
	}

	return s.repo.AddTrackToPlaylist(playlistID, trackID)
}

func (s *PlaylistTrackService) GetAllTracks(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) ([]entities.TrackMeta, error) {
	canView, err := s.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, errors.New("permission denied: user cannot view playlist")
	}

	return s.repo.GetTracksInPlaylist(playlistID)
}

func (s *PlaylistTrackService) DeleteTrack(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID, trackID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return errors.New("permission denied: user cannot edit playlist")
	}

	return s.repo.DeleteTrackFromPlaylist(playlistID, trackID)
}

func (s *PlaylistTrackService) DeleteAllTracks(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) error {
	canEdit, err := s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !canEdit {
		return errors.New("permission denied: user cannot edit playlist")
	}

	tracks, err := s.repo.GetTracksInPlaylist(playlistID)
	if err != nil {
		return err
	}

	for _, track := range tracks {
		err := s.repo.DeleteTrackFromPlaylist(playlistID, track.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
