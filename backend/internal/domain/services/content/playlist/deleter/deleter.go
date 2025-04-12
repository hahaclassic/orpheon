package deleter

import (
	"context"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

// TODO: нет транзакции. Нет восстанавливающих операций. Надо продумать,
// как сделать так, чтобы сохранялась целостность данных.

type MetaDeletionService interface {
	Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}

type TrackDeletionService interface {
	Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}

type FavoritesDeletionService interface {
	Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}

type PlaylistCoverDeletionService interface {
	Delete(ctx context.Context, claims *entity.Claims, objectID uuid.UUID) error
}

type PlaylistDeleter struct {
	meta      MetaDeletionService
	tracks    TrackDeletionService
	favorites FavoritesDeletionService
	cover     PlaylistCoverDeletionService
}

func New(meta MetaDeletionService, track TrackDeletionService,
	favorites FavoritesDeletionService, cover PlaylistCoverDeletionService) *PlaylistDeleter {
	return &PlaylistDeleter{
		meta:      meta,
		tracks:    track,
		favorites: favorites,
		cover:     cover,
	}
}

func (p *PlaylistDeleter) Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	err := p.favorites.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	err = p.tracks.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	err = p.meta.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	err = p.cover.Delete(ctx, claims, playlistID)
	if err != nil {
		return err
	}

	return nil
}

type DeletionService interface {
	Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error
}

type PlaylistDeleterV2 struct {
	deleters []DeletionService
}

func NewV2(deleters []DeletionService) *PlaylistDeleterV2 {
	return &PlaylistDeleterV2{
		deleters: deleters,
	}
}

func (p *PlaylistDeleterV2) Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	for i := range p.deleters {
		err := p.deleters[i].Delete(ctx, claims, playlistID)
		if err != nil {
			return err
		}
	}

	return nil
}

type Deleter func(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error

type PlaylistDeleterV3 struct {
	deleters []Deleter
}

func NewV3(deleters []Deleter) *PlaylistDeleterV3 {
	return &PlaylistDeleterV3{
		deleters: deleters,
	}
}

func (p *PlaylistDeleterV3) Delete(ctx context.Context, claims *entity.Claims, playlistID uuid.UUID) error {
	for i := range p.deleters {
		err := p.deleters[i](ctx, claims, playlistID)
		if err != nil {
			return err
		}
	}

	return nil
}
