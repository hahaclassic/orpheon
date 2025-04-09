package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

// TODO: обработка ошибок

var (
	ErrForbidden      = errors.New("access denied")
	ErrNotFound       = errors.New("playlist not found")
	ErrTrackExists    = errors.New("track already in playlist")
	ErrTrackNotExists = errors.New("track not in playlist")
)

type PlaylistPolicyService interface {
	CanDelete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanEdit(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
	CanView(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (bool, error)
}

type PlaylistMetaRepository interface {
	Create(ctx context.Context, playlist *entity.Playlist) error
	GetByID(ctx context.Context, playlistID uuid.UUID) (*entity.Playlist, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Playlist, error)
	Update(ctx context.Context, playlist *entity.Playlist) error
	Delete(ctx context.Context, playlistID uuid.UUID) error
}

type PlaylistMetaService struct {
	repo   PlaylistMetaRepository
	policy PlaylistPolicyService
}

func NewPlaylistMetaService(repo PlaylistMetaRepository, policy PlaylistPolicyService) *PlaylistMetaService {
	return &PlaylistMetaService{
		repo:   repo,
		policy: policy,
	}
}

func (p *PlaylistMetaService) CreatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error {
	if playlist.Name == "" {
		return errors.New("playlist name cannot be empty")
	}

	return p.repo.Create(ctx, playlist)
}

func (p *PlaylistMetaService) GetPlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (*entity.Playlist, error) {
	ok, err := p.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}

	return p.repo.GetByID(ctx, playlistID)
}

func (p *PlaylistMetaService) GetUserPlaylists(ctx context.Context, claims *entity.Claims, userID uuid.UUID) ([]*entity.Playlist, error) {
	playlists, err := p.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// if user is owner, show all playlists
	if claims.UserID == userID {
		return playlists, nil
	}

	publicPlaylists := playlists[:]
	currIdx := 0
	for i := range playlists {
		if !playlists[i].IsPrivate {
			publicPlaylists[currIdx] = playlists[i]
			currIdx++
		}
	}

	return publicPlaylists, nil
}

func (p *PlaylistMetaService) UpdatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) error {
	ok, err := p.policy.CanEdit(ctx, claims, playlist.ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}

	return p.repo.Update(ctx, playlist)
}

func (p *PlaylistMetaService) DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	ok, err := p.policy.CanDelete(ctx, claims, playlistID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}

	return p.repo.Delete(ctx, playlistID)
}
