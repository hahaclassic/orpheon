package policyservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entities"
)

// TODO: 1. надо учесть, что права могут инвалидироваться (если плейлист стал непубличным, например)
//  	 2. обновить логику с учетом того, что PlaylistViewer тоже может кешироваться

type PlaylistAccessCache interface {
	Set(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID, lvl entities.PlaylistAccessLvl) error
	Get(ctx context.Context, userID uuid.UUID, playlistID uuid.UUID) (lvl entities.PlaylistAccessLvl, err error)
}

type PlaylistRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Playlist, error)
}

type PlaylistPolicyService struct {
	cache PlaylistAccessCache // fast path
	repo  PlaylistRepository  // long path
}

func NewPlaylistPolicyService(localCache PlaylistAccessCache, cache PlaylistAccessCache, repo PlaylistRepository) *PlaylistPolicyService {
	return &PlaylistPolicyService{}
}

// Owner + Admin (if private=false)
func (p *PlaylistPolicyService) CanDelete(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error) {
	var (
		lvl entities.PlaylistAccessLvl
		err error
	)

	lvl, err = p.cache.Get(ctx, claims.UserID, playlistID)
	if err != nil { // CHECK CACHE MISS ERROR !!!
		return false, err
	}

	if lvl == entities.PlaylistOwnerLvl {
		return true, nil
	}

	playlist, err := p.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return false, err
	}

	return playlist.OwnerID == claims.UserID || (claims.AccessLvl == entities.Admin && !playlist.IsPrivate), nil
}

// Owner
func (p *PlaylistPolicyService) CanEdit(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error) {
	var (
		lvl entities.PlaylistAccessLvl
		err error
	)

	lvl, err = p.cache.Get(ctx, claims.UserID, playlistID)
	if err != nil { // CHECK CACHE MISS ERROR !!!
		return false, err
	}

	if lvl == entities.PlaylistOwnerLvl {
		return true, nil
	}

	playlist, err := p.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return false, err
	}

	isOwner := playlist.OwnerID == claims.UserID

	if isOwner {
		err = p.cache.Set(ctx, playlist.OwnerID, playlistID, entities.PlaylistOwnerLvl)
		if err != nil {
			return isOwner, err
		}
	}

	return isOwner, nil
}

// All users (if private=false)
// Если userID != ownerID
func (p *PlaylistPolicyService) CanView(ctx context.Context, claims *entities.Claims, playlistID uuid.UUID) (bool, error) {
	var (
		lvl entities.PlaylistAccessLvl
		err error
	)

	lvl, err = p.cache.Get(ctx, claims.UserID, playlistID)
	if err != nil { // CHECK CACHE MISS ERROR !!!
		return false, err
	}

	if lvl == entities.PlaylistOwnerLvl {
		return true, nil
	}

	playlist, err := p.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return false, err
	}

	isOwner := playlist.OwnerID == claims.UserID

	if isOwner {
		err = p.cache.Set(ctx, playlist.OwnerID, playlistID, entities.PlaylistOwnerLvl)
		if err != nil {
			return isOwner, err
		}
	}

	// TODO: Надо ли кешировать то, что плейлист публичный?

	return isOwner || !playlist.IsPrivate, nil
}
