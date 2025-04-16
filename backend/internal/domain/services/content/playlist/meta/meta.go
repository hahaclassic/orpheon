package meta

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

// TODO: обработка ошибок

var (
	ErrEmptyPlaylistName = errors.New("playlist name cannot be empty")
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
	policy usecase.PlaylistPolicyService
}

func NewPlaylistMetaService(repo PlaylistMetaRepository, policy usecase.PlaylistPolicyService) *PlaylistMetaService {
	return &PlaylistMetaService{
		repo:   repo,
		policy: policy,
	}
}

func (p *PlaylistMetaService) CreatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreatePlaylist, err)
	}()

	if playlist.Name == "" {
		return ErrEmptyPlaylistName
	}

	return p.repo.Create(ctx, playlist)
}

func (p *PlaylistMetaService) GetPlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (_ *entity.Playlist, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetPlaylist, err)
	}()

	err = p.policy.CanView(ctx, claims, playlistID)
	if err != nil {
		return nil, err
	}

	return p.repo.GetByID(ctx, playlistID)
}

func (p *PlaylistMetaService) GetUserPlaylists(ctx context.Context, claims *entity.Claims, userID uuid.UUID) (_ []*entity.Playlist, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetUserPlaylists, err)
	}()

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

func (p *PlaylistMetaService) UpdatePlaylist(ctx context.Context, claims *entity.Claims, playlist *entity.Playlist) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdatePlaylist, err)
	}()

	err = p.policy.CanEdit(ctx, claims, playlist.ID)
	if err != nil {
		return err
	}

	return p.repo.Update(ctx, playlist)
}

func (p *PlaylistMetaService) DeletePlaylist(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeletePlaylist, err)
	}()

	err = p.policy.CanDelete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return p.repo.Delete(ctx, playlistID)
}
