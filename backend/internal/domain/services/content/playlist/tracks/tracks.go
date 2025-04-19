package tracks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type PlaylistTracksRepository interface {
	AddTrackToPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteTrackFromPlaylist(ctx context.Context, playlistID uuid.UUID, trackID uuid.UUID) error
	DeleteAllTracksFromPlaylist(ctx context.Context, playlistID uuid.UUID) error
	GetAllPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]*entity.TrackMeta, error)
}

type PlaylistTrackService struct {
	repo   PlaylistTracksRepository
	policy usecase.PlaylistPolicyService
}

func NewPlaylistTrackService(repo PlaylistTracksRepository, policy usecase.PlaylistPolicyService) *PlaylistTrackService {
	return &PlaylistTrackService{
		repo:   repo,
		policy: policy,
	}
}

func (s *PlaylistTrackService) AddTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrAddTrack, err)
	}()

	err = s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return s.repo.AddTrackToPlaylist(ctx, playlistID, trackID)
}

func (s *PlaylistTrackService) GetAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (tracks []*entity.TrackMeta, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllTracks, err)
	}()

	err = s.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAllPlaylistTracks(ctx, playlistID)
}

func (s *PlaylistTrackService) DeleteTrack(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID, trackID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteTrack, err)
	}()

	err = s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTrackFromPlaylist(ctx, playlistID, trackID)
}

func (s *PlaylistTrackService) DeleteAllTracks(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteAllTracks, err)
	}()

	err = s.policy.CanEdit(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return s.repo.DeleteAllTracksFromPlaylist(ctx, playlistID)
}
